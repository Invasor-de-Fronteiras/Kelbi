package adapters

import (
	"erupe-ce/common/byteframe"

	"github.com/jmoiron/sqlx"
)

type QuestLoaderInDb struct {
	db *sqlx.DB
}

func NewQuestLoaderInDb(db *sqlx.DB) QuestLoader {
	return &QuestLoaderInDb{db: db}
}

func (ql *QuestLoaderInDb) QuestBinById(id string) (questBin []byte, err error) {
	err = ql.db.QueryRow("SELECT quest_bin FROM quests WHERE id = $1", id).Scan(&questBin)
	return
}

func (ql *QuestLoaderInDb) QuestListBinById(id string) (questListBin []byte, err error) {
	err = ql.db.QueryRow("SELECT quest_list_bin FROM quests WHERE id = $1", id).Scan(&questListBin)
	return
}

func (ql *QuestLoaderInDb) QuestCount() (count uint16, err error) {
	err = ql.db.QueryRow("SELECT COUNT(*) FROM quests WHERE enabled = true").Scan(&count)
	return
}

func (ql *QuestLoaderInDb) Quests(take uint16, skip uint16) (questList []byte, err error) {
	bf := byteframe.NewByteFrame()

	buffer, returnedCount, err := ql.NextQuests(take, skip)

	if err != nil {
		return nil, err
	}

	totalCount, err := ql.QuestCount()

	if totalCount == 0 {
		return make([]byte, 18), nil
	}

	bf.WriteUint16(returnedCount)
	bf.WriteBytes(buffer.Data())
	bf.WriteUint16(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint32(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint16(totalCount)
	bf.WriteUint16(take)

	questList = bf.Data()

	return
}

type IteratorQuest struct {
	Id            string `db:"id"`
	QuestListSize uint16 `db:"quest_list_bin_size"`
}

func (ql *QuestLoaderInDb) NextQuest(skip uint16) (iter IteratorQuest, err error) {
	iter = IteratorQuest{}
	err = ql.db.Get(iter, "SELECT id,quest_list_bin_size FROM quests LIMIT 1 OFFSET $1", skip)
	return
}

func (ql *QuestLoaderInDb) NextQuests(take uint16, skip uint16) (buffer *byteframe.ByteFrame, count uint16, err error) {
	buffer = byteframe.NewByteFrame()
	count = 0

	var bufferSize uint16 = 0
	var maxBufferSize uint16 = 64_000

	for bufferSize <= maxBufferSize {
		if count >= take {
			break
		}

		iter, err := ql.NextQuest(skip)

		if err != nil {
			return nil, 0, err
		}

		if (iter.QuestListSize + bufferSize) > maxBufferSize {
			break
		}

		skip++

		questListBin, err := ql.QuestListBinById(iter.Id)

		if err != nil {
			return nil, 0, err
		}

		buffer.WriteBytes(questListBin)
	}

	return
}
