package adapters

type QuestCategory = string

const (
	Diva QuestCategory = "DIVA"
)

type Quest struct {
	Id            string // auto incriment
	Name          string
	Objective     string
	Category      QuestCategory
	Enabled       bool
	Metadata      string //
	QuestBinSize  int32
	QuestListSize int32
	QuestListBin  string
	QuestBin      string
}

type QuestLoader interface {
	QuestBinById(id string) ([]byte, error)
	Quests(take uint16, skip uint16) ([]byte, error)
}
