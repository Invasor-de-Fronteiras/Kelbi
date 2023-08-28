package adapters

import (
	"erupe-ce/config"
	"fmt"
	"os"
	"path/filepath"
)

type QuestLoaderQuestlist struct {
	erupeConfig *config.Config
}

func NewQuestLoaderQuestlist(erupeConfig *config.Config) QuestLoader {
	var questLoader = &QuestLoaderQuestlist{
		erupeConfig: erupeConfig,
	}

	return questLoader
}

func (ql *QuestLoaderQuestlist) QuestBinById(id string) (quest []byte, err error) {
	filepath := filepath.Join(ql.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", id))
	quest, err = os.ReadFile(filepath)
	return
}

func (ql *QuestLoaderQuestlist) Quests(_take uint16, skip uint16, _dev bool) (quests []byte, err error) {
	filename := fmt.Sprintf("questlists/list_%d.bin", skip)
	filepath := filepath.Join(ql.erupeConfig.BinPath, filename)
	quests, err = os.ReadFile(filepath)
	return
}
