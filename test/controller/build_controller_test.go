package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestBuildController_Upload_ConnectsWithWarpServerAndWarpToken(t *testing.T) {
	t.Setenv("WARP_SERVER", "https://games.example.org")
	t.Setenv("WARP_TOKEN", "we_upload")
	t.Setenv("WARP_ENGINE_URL", "")
	t.Setenv("WARP_ENGINE_TOKEN", "")
	var selection domain.Selection
	auth := &mocksvc.MockAuthService{ConnectFunc: func(given domain.Selection) (*domain.ProfileModel, error) {
		selection = given
		return &domain.ProfileModel{}, nil
	}}
	var target domain.BuildTargetModel
	var files []string
	service := &mocksvc.MockBuildService{UploadFunc: func(given domain.BuildTargetModel, file_paths []string) ([]domain.BuildUploadModel, error) {
		target, files = given, file_paths
		return []domain.BuildUploadModel{{File: "pong-1.2.0-mac-x64.zip", Size: 2048, SHA256: "abc"}}, nil
	}}

	stdout, _, err := execute(cli.Router{BuildController: cli.BuildController{BuildService: service, AuthService: auth}},
		"release", "upload", "--name", "pong", "--version", "1.2.0", "pong-1.2.0-mac-x64.zip")

	assert.NoError(t, err)
	assert.Equal(t, domain.Selection{Server: "https://games.example.org", Token: "we_upload"}, selection)
	assert.Equal(t, domain.BuildTargetModel{Name: "pong", Version: "1.2.0"}, target)
	assert.Equal(t, []string{"pong-1.2.0-mac-x64.zip"}, files)
	assert.Contains(t, stdout, "pong-1.2.0-mac-x64.zip")
}

func TestBuildController_Upload_FlagsWinOverTheEnvironment(t *testing.T) {
	t.Setenv("WARP_SERVER", "https://env.example.org")
	t.Setenv("WARP_TOKEN", "we_env")
	var selection domain.Selection
	auth := &mocksvc.MockAuthService{ConnectFunc: func(given domain.Selection) (*domain.ProfileModel, error) {
		selection = given
		return &domain.ProfileModel{}, nil
	}}

	_, _, err := execute(cli.Router{BuildController: cli.BuildController{BuildService: &mocksvc.MockBuildService{}, AuthService: auth}},
		"release", "upload", "--server", "https://flag.example.org", "--token", "we_flag", "--name", "p", "--version", "1", "p-1.zip")

	assert.NoError(t, err)
	assert.Equal(t, "https://flag.example.org", selection.Server)
	assert.Equal(t, "we_flag", selection.Token)
}

func TestBuildController_Upload_NeedsAFile(t *testing.T) {
	_, _, err := execute(cli.Router{BuildController: cli.BuildController{BuildService: &mocksvc.MockBuildService{}, AuthService: &mocksvc.MockAuthService{}}},
		"release", "upload", "--name", "p", "--version", "1")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestBuildController_Publish_SendsTheTarget(t *testing.T) {
	var target domain.BuildTargetModel
	service := &mocksvc.MockBuildService{PublishFunc: func(given domain.BuildTargetModel) (*domain.BuildPublishModel, error) {
		target = given
		return &domain.BuildPublishModel{Published: true, Name: "gio", Platform: "golang", Version: "0.3.0"}, nil
	}}

	stdout, _, err := execute(cli.Router{BuildController: cli.BuildController{BuildService: service, AuthService: &mocksvc.MockAuthService{}}},
		"release", "publish", "--name", "gio", "--version", "0.3.0", "--platform", "golang", "--app", "gui", "--commit", "3fb3ed8")

	assert.NoError(t, err)
	assert.Equal(t, domain.BuildTargetModel{Name: "gio", Version: "0.3.0", Platform: "golang", App: "gui", Commit: "3fb3ed8"}, target)
	assert.Contains(t, stdout, "golang")
}

func TestMakefileController_Check_FailsOnABrokenContract(t *testing.T) {
	var directory string
	service := &mocksvc.MockMakefileService{CheckFunc: func(given string, platform string, app string) (*domain.MakefileReportModel, error) {
		directory = given
		return &domain.MakefileReportModel{Directory: given, Platform: "godot", Targets: []domain.MakefileTargetModel{
			{Name: "build", Status: domain.MakefileTargetPresent},
			{Name: "check", Status: domain.MakefileTargetMissing},
		}}, nil
	}}

	stdout, stderr, err := execute(cli.Router{MakefileController: cli.MakefileController{MakefileService: service, AuthService: &mocksvc.MockAuthService{}}},
		"check", "makefile", "../pong")

	assert.Error(t, err)
	assert.Equal(t, cli.ExitError, cli.ExitCode(err))
	assert.Equal(t, "../pong", directory)
	assert.Contains(t, stdout, "missing")
	assert.Contains(t, stderr, "platform godot")
}

func TestMakefileController_Check_PassesWithEveryTarget(t *testing.T) {
	service := &mocksvc.MockMakefileService{CheckFunc: func(given string, platform string, app string) (*domain.MakefileReportModel, error) {
		assert.Equal(t, "golang", platform)
		assert.Equal(t, "gui", app)
		return &domain.MakefileReportModel{Directory: given, Platform: platform, Targets: []domain.MakefileTargetModel{{Name: "build", Status: domain.MakefileTargetPresent}}}, nil
	}}

	_, _, err := execute(cli.Router{MakefileController: cli.MakefileController{MakefileService: service, AuthService: &mocksvc.MockAuthService{}}},
		"check", "makefile", "--platform", "golang", "--app", "gui")

	assert.NoError(t, err)
}
