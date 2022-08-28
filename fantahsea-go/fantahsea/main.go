package main

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/fantahsea/web/controller"

	"github.com/curtisnewbie/fantahsea/config"

	log "github.com/sirupsen/logrus"
)

func main() {

	profile := config.ParseProfile(os.Args[1:])
	log.Printf("Using profile: %v", profile)

	conf, err := config.ParseJsonConfig(fmt.Sprintf("app-conf-%v.json", profile))
	if err != nil {
		panic(err)
	}
	config.SetGlobalConfig(conf)

	if err := config.InitDBFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}

	isProd := profile == "prod"
	err = controller.BootstrapServer(&conf.ServerConf, isProd)
	if err != nil {
		panic(err)
	}

}
