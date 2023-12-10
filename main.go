package main

import (
	"os"

	"github.com/curtisnewbie/fantahsea/fantahsea"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func main() {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(func(rail miso.Rail) error { return fantahsea.RegisterRoutes(rail) })
	miso.PreServerBootstrap(func(rail miso.Rail) error { return fantahsea.PrepareEventBus(rail) })
	miso.BootstrapServer(os.Args)
}
