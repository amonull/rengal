package main

import (
	"github.com/samber/lo"

	"github.com/amonull/rengal/cmd"
	"github.com/amonull/rengal/config"
	"github.com/amonull/rengal/log"
)

func main() {
	lo.Must0(config.Setup())
	lo.Must0(log.Setup())
	cmd.Execute()
}
