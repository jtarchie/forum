package cmd

type CLI struct {
	Server ServerCmd `cmd:"" help:"start the server"`
}
