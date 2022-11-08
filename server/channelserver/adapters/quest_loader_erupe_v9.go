package adapters

import (
	"erupe-ce/config"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type QuestLoaderErupeV9 struct {
	erupeConfig *config.Config
}

func (ql *QuestLoaderErupeV9) QuestBinById(id string) (quest []byte, err error) {
	filepath := filepath.Join(ql.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", id))
	quest, err = ioutil.ReadFile(filepath)
	return
}

func (ql *QuestLoaderErupeV9) Quests(_take uint16, skip uint16) (quests []byte, err error) {
	filename := fmt.Sprintf("questlists/list_%d.bin", skip)
	filepath := filepath.Join(ql.erupeConfig.BinPath, filename)
	quests, err = ioutil.ReadFile(filepath)
	return
}
