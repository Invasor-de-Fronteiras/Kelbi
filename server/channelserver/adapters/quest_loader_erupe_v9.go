package adapters

import (
	"erupe-ce/config"
	"fmt"
	"os"
	"path/filepath"
)

type QuestLoaderErupeV9 struct {
	erupeConfig *config.Config
}

func NewQuestLoaderErupeV9(erupeConfig *config.Config) QuestLoader {
	var questLoader = &QuestLoaderErupeV9{
		erupeConfig: erupeConfig,
	}

	return questLoader
}

func (ql *QuestLoaderErupeV9) QuestBinById(id string) (quest []byte, err error) {
	filepath := filepath.Join(ql.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", id))
	quest, err = os.ReadFile(filepath)
	return
}

func (ql *QuestLoaderErupeV9) Quests(_take uint16, skip uint16, _dev bool) (quests []byte, err error) {
	filename := fmt.Sprintf("questlists/list_%d.bin", skip)
	filepath := filepath.Join(ql.erupeConfig.BinPath, filename)
	quests, err = os.ReadFile(filepath)
	return
}
