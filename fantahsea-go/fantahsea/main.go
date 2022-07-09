package main

import (
	"fantahsea/config"
	"fantahsea/web/controller"
	"fmt"
	"log"
	"os"
)

func main() {

	profile := config.ParseProfile(os.Args[1:])
	log.Printf("Using profile: %v", profile)

	conf, err := config.ParseJsonConfig(fmt.Sprintf("app-conf-%v.json", profile))
	if err != nil {
		panic(err)
	}

	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}

	err = controller.BootstrapServer(&conf.ServerConf)
	if err != nil {
		panic(err)
	}

}
