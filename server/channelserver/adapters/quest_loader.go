package adapters

type QuestCategory = string

const (
	Diva QuestCategory = "DIVA"
)

type Quest struct {
	Id           uint32        `db:"id"`
	QuestId      uint32        `db:"quest_id"`
	Period       string        `db:"period"`
	Season       string        `db:"season"`
	Name         string        `db:"name"`
	Objective    string        `db:"objective"`
	Category     QuestCategory `db:"category"`
	Enabled      bool          `db:"enabled"`
	Metadata     string        `db:"metadata"`
	QuestBin     string        `db:"quest_bin"`
	QuestListBin string        `db:"quest_list_bin"`
}

type QuestLoader interface {
	QuestBinById(id string) ([]byte, error)
	Quests(take uint16, skip uint16) ([]byte, error)
}
