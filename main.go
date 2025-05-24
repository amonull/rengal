package main

import (
	"github.com/amonull/rengal/cmd"
	"github.com/amonull/rengal/config"
	"github.com/amonull/rengal/log"
	"github.com/samber/lo"
)

func main() {
	lo.Must0(config.Setup())
	lo.Must0(log.Setup())
	cmd.Execute()
}
