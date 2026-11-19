package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeController_Index_PassesNameFilter(t *testing.T) {
	svc := &mocksvc.MockRuntimeService{AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.RuntimeVersionModel], error) {
		assert.Equal(t, domain.Filter{"name": "tic80"}, filter)
		return &domain.PageModel[domain.RuntimeVersionModel]{Items: []domain.RuntimeVersionModel{{ID: 1, Name: "tic80", Version: "1.2", SizeBytes: 3 << 20}}}, nil
	}}

	stdout, _, err := execute(cli.Router{RuntimeController: cli.RuntimeController{RuntimeService: svc}}, "runtime", "list", "--name", "tic80")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "3.0 MiB")
}

func TestRuntimeController_Install_PassesVersionAndDryRun(t *testing.T) {
	var got_name, got_version string
	var got_dry_run bool
	svc := &mocksvc.MockRuntimeService{InstallFunc: func(name string, version string, dry_run bool) (*domain.RuntimeInstallModel, error) {
		got_name, got_version, got_dry_run = name, version, dry_run
		return &domain.RuntimeInstallModel{Name: name, Version: "1.2", DryRun: dry_run, Fetched: []string{"tic80.js"}}, nil
	}}
	router := cli.Router{RuntimeController: cli.RuntimeController{RuntimeService: svc}}

	stdout, _, err := execute(router, "runtime", "install", "tic80", "--dry-run")
	assert.NoError(t, err)
	assert.Equal(t, "tic80", got_name)
	assert.Equal(t, "", got_version)
	assert.True(t, got_dry_run)
	assert.Contains(t, stdout, "tic80.js")

	_, _, err = execute(router, "runtimes", "install", "tic80", "1.2")
	assert.NoError(t, err)
	assert.Equal(t, "1.2", got_version)
	assert.False(t, got_dry_run)
}

func TestRuntimeController_Install_NeedsName(t *testing.T) {
	_, _, err := execute(cli.Router{RuntimeController: cli.RuntimeController{RuntimeService: &mocksvc.MockRuntimeService{}}}, "runtime", "install")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}
