package adapters

import "erupe-ce/config"

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

func QuestLoaderAdapter(erupeConfig config.Config) QuestLoader {
	return &QuestLoaderErupeV9{
		erupeConfig: &erupeConfig,
	}
}
