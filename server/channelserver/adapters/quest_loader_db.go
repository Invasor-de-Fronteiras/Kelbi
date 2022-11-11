package adapters

import (
	"erupe-ce/common/byteframe"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type QuestLoaderInDb struct {
	db *sqlx.DB
}

func NewQuestLoaderInDb(db *sqlx.DB) QuestLoader {
	return &QuestLoaderInDb{db: db}
}

func (_ *QuestLoaderInDb) IdFromFilename(id string) (questId uint64, period string, season uint64, err error) {
	questId, _ = strconv.ParseUint(id[0:5], 10, 32)
	period = string(id[5])
	season, _ = strconv.ParseUint(string(id[6]), 10, 32)
	return
}

func (ql *QuestLoaderInDb) QuestBinById(id string) (questBin []byte, err error) {
	questId, period, season, _ := ql.IdFromFilename(id)
	query := "SELECT quest_bin FROM quests WHERE quest_id = $1 AND period = $2 AND season = $3"
	err = ql.db.QueryRow(query, questId, period, season).Scan(&questBin)
	return
}

func (ql *QuestLoaderInDb) QuestCount() (count uint16, err error) {
	err = ql.db.QueryRow("SELECT COUNT(*) FROM quests WHERE enabled = true").Scan(&count)
	return
}

func (ql *QuestLoaderInDb) Quests(_take uint16, skip uint16) (questList []byte, err error) {
	bf := byteframe.NewByteFrame()

	buffer, returnedCount, err := ql.QuestListBuffer(100, skip)

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
	bf.WriteUint16(skip)

	questList = bf.Data()

	return
}

type QuestListBin struct {
	QuestList []byte `db:"quest_list_bin"`
}

func (ql *QuestLoaderInDb) QuestListBin(take uint16, skip uint16) (quests []QuestListBin, err error) {
	quests = []QuestListBin{}
	query := "SELECT quest_list_bin FROM quests WHERE enabled = true LIMIT $1 OFFSET $2"
	err = ql.db.Select(&quests, query, take, skip)
	return
}

func (ql *QuestLoaderInDb) QuestListBuffer(take uint16, skip uint16) (buffer *byteframe.ByteFrame, count uint16, err error) {
	buffer = byteframe.NewByteFrame()
	count = 0

	var bufferSize = 0
	var maxBufferSize = 64_000

	for bufferSize <= maxBufferSize {

		quests, err := ql.QuestListBin(take, skip+count)

		if err != nil {
			break
		}

		if len(quests) == 0 {
			break
		}

		for _, quest := range quests {
			size := len(quest.QuestList)
			bufferSize += size

			if bufferSize > maxBufferSize {
				break
			}

			count++

			buffer.WriteBytes(quest.QuestList)
		}
	}

	return
}
