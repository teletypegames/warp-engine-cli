package main

import (
	"os"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

func main() {
	client := api.New("", "", "warp-engine-cli/"+cli.Version)
	d := domain.NewDomain(client, domain.DefaultConfigPath(), domain.DefaultLocalDirectory())
	os.Exit(cli.Start(d, os.Args[1:], os.Stdout, os.Stderr))
}
