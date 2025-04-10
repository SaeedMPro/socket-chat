package config

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

func GetClientConfig(clientType string, clientConfig model.Config) (model.Client, model.Client) {
	switch clientType {
	case "client-one":
		return clientConfig.ClientOne, clientConfig.ClientTwo
	case "client-two":
		return clientConfig.ClientTwo, clientConfig.ClientOne
	default:
		panic("Invalid argument input!!")
	}
}
