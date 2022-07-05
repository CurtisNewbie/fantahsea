package main

import (
	"fantahsea/config"
	"fantahsea/web/controller"
)


func main() {

	conf, err := config.ParseJsonConfig("app-conf.json");
	if err != nil {
		panic(err)
	}

	if err := config.InitDbFromConfig(&conf.DBConf); err != nil {
		panic(err)
	}

	err = controller.BootstrapServer(&conf.ServerConf)
	if err != nil {
		panic(err)
	}

}