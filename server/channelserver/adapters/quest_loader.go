package adapters

type Quest struct {
	Id            string
	QuestListType string
	QuestListBin  string
	QuestBin      string
}

type QuestLoader interface {
	QuestBinById(id string) ([]byte, error)
	Quests(take uint16, skip uint16) ([]byte, error)
}
