package controller_test

import (
	"bytes"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
)

// execute runs the command line args on the tree of router and returns
// what it printed on stdout and stderr, and its error. A router without an
// AuthController gets one over a [mocksvc.MockAuthService] that connects
// to https://games.example.org.
func execute(router cli.Router, args ...string) (string, string, error) {
	if router.AuthController == nil {
		router.AuthController = cli.AuthController{AuthService: &mocksvc.MockAuthService{}, Session: &cli.Session{}}
	}
	root := router.Init()
	var stdout, stderr bytes.Buffer
	root.SetArgs(args)
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

// text returns a pointer to value, for the optional fields of the models.
func text(value string) *string {
	return &value
}
