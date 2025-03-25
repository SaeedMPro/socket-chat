package main

import (
	"fmt"
	"os"

	"github.com/SaeedMPro/socket-chat/model"
	"github.com/SaeedMPro/socket-chat/util"
)

var clientConfig model.Config

func init() {
	if err := util.LoadConfig("./config.json", &clientConfig); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	util.ClearScreen()
}

func main() {

}
