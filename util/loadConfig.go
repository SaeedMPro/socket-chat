package util

import (
	"encoding/json"
	"os"

	"github.com/SaeedMPro/socket-chat/model"
)

func LoadConfig(fileName string, config *model.Config) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}
