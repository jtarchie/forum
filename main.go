package main

import (
	"github.com/alecthomas/kong"
	"github.com/jtarchie/forum/cmd"
)

func main() {
	cli := &cmd.CLI{}
	ctx := kong.Parse(cli,
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
