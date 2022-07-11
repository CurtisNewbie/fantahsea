package main

import (
	"fantahsea/config"
	"fantahsea/web/controller"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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


	isProd := profile == "prod" 
	err = controller.BootstrapServer(&conf.ServerConf, isProd)
	if err != nil {
		panic(err)
	}

}
