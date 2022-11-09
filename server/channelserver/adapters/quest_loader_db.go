package adapters

import (
	"erupe-ce/common/byteframe"
	"io"

	"github.com/jmoiron/sqlx"
)

type QuestLoaderInDb struct {
	db *sqlx.DB
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
	bf.WriteUint16(0)

	returnedCount, err := ql.NextQuests(bf, take, skip)

	if err != nil {
		return nil, err
	}

	totalCount, err := ql.QuestCount()

	if totalCount == 0 {
		return make([]byte, 18), nil
	}

	bf.WriteUint16(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint32(0) // Unk
	bf.WriteUint16(0) // Unk
	bf.WriteUint16(totalCount)
	bf.WriteUint16(take)
	bf.Seek(0, io.SeekStart)
	bf.WriteUint16(returnedCount)

	return
}

type IteratorQuest struct {
	Id            string
	QuestListSize uint16
}

func (ql *QuestLoaderInDb) NextQuest(skip uint16) (iter IteratorQuest, err error) {

	return
}

func (ql *QuestLoaderInDb) NextQuests(buffer *byteframe.ByteFrame, take uint16, skip uint16) (count uint16, err error) {
	count = 0

	var bufferSize uint16 = 0
	var maxBufferSize uint16 = 64_000

	for bufferSize <= maxBufferSize {
		if count >= take {
			break
		}

		iter, err := ql.NextQuest(skip)

		if err != nil {
			return 0, err
		}

		if (iter.QuestListSize + bufferSize) > maxBufferSize {
			break
		}

		skip++

		questListBin, err := ql.QuestListBinById(iter.Id)

		if err != nil {
			return 0, err
		}

		buffer.WriteBytes(questListBin)
	}

	return
}
