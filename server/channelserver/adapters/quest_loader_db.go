package adapters

import (
	"erupe-ce/common/byteframe"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type QuestLoaderDb struct {
	db *sqlx.DB
}

func NewQuestLoaderDb(db *sqlx.DB) QuestLoader {
	return &QuestLoaderDb{db: db}
}

func (*QuestLoaderDb) IdFromFilename(id string) (questId uint64, period string, season string, err error) {
	questId, _ = strconv.ParseUint(id[0:5], 10, 32)

	if string(id[5]) == "d" {
		period = "DAY"
	} else {
		period = "NIGHT"
	}

	seasonId, _ := strconv.ParseUint(string(id[6]), 10, 32)

	switch seasonId {
	case 0:
		season = "WARM"
	case 1:
		season = "COLD"
	default:
		season = "BREED"
	}

	return
}

func (ql *QuestLoaderDb) QuestBinById(id string) (questBin []byte, err error) {
	questId, period, season, _ := ql.IdFromFilename(id)
	query := "SELECT quest_bin FROM quests WHERE quest_id = $1 AND period = $2 AND season = $3"
	err = ql.db.QueryRow(query, questId, period, season).Scan(&questBin)
	return
}

func (ql *QuestLoaderDb) QuestCount(dev bool) (count uint16, err error) {
	onlyDevFilter := ""
	if !dev {
		onlyDevFilter = " AND only_dev=false "
	}
	err = ql.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM questlist WHERE enable = true %s", onlyDevFilter)).Scan(&count)
	return
}

func (ql *QuestLoaderDb) Quests(_take uint16, skip uint16, dev bool) (questList []byte, err error) {
	bf := byteframe.NewByteFrame()

	buffer, returnedCount, err := ql.QuestListBuffer(100, skip, dev)

	if err != nil {
		return nil, err
	}

	totalCount, err := ql.QuestCount(dev)

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
	QuestList []byte `db:"questlist_bin"`
}

func (ql *QuestLoaderDb) QuestListBin(take uint16, skip uint16, dev bool) (quests []QuestListBin, err error) {
	quests = []QuestListBin{}
	onlyDevFilter := ""
	if !dev {
		onlyDevFilter = " AND only_dev=false "
	}

	query := fmt.Sprintf("SELECT questlist_bin FROM questlist WHERE enable = true %s ORDER BY position LIMIT $1 OFFSET $2", onlyDevFilter)
	err = ql.db.Select(&quests, query, take, skip)
	return
}

func (ql *QuestLoaderDb) QuestListBuffer(take uint16, skip uint16, dev bool) (buffer *byteframe.ByteFrame, count uint16, err error) {
	buffer = byteframe.NewByteFrame()
	count = 0

	var bufferSize = 0
	var maxBufferSize = 64_000

	for bufferSize <= maxBufferSize {

		quests, err := ql.QuestListBin(take, skip+count, dev)

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
