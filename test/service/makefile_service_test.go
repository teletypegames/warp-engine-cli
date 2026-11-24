package service_test

import (
	"errors"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

// makefilePlatforms is a server with ebitengine (web, mac by hand) and
// golang, whose gui variant has both and whose cli variant neither.
func makefilePlatforms() []domain.BuildPlatformModel {
	return []domain.BuildPlatformModel{
		{Name: "ebitengine", Kinds: []domain.BuildKindModel{{Name: "html", Built: true}, {Name: "linux_x64", Built: true}, {Name: "mac_arm64"}}},
		{Name: "c64", Kinds: []domain.BuildKindModel{{Name: "cartridge", Built: true}}},
		{Name: "golang", Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64", Built: true}}, Apps: []domain.BuildAppModel{
			{Name: "cli", Default: true, Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64", Built: true}}},
			{Name: "gui", Kinds: []domain.BuildKindModel{{Name: "html", Built: true}, {Name: "mac_x64"}}},
		}},
	}
}

// makefileService is a [domain.MakefileService] over a Makefile that has
// the targets present, the marker platform/app and the server's platforms
// (a failing server when platforms is nil).
func makefileService(present map[string]bool, platform string, app string, platforms []domain.BuildPlatformModel) domain.MakefileService {
	return domain.MakefileService{
		MakefileRepository: &mockrepo.MockMakefileRepository{
			ExistsFunc: func(directory string) bool { return true },
			ProbeFunc: func(directory string, target string) (string, string, error) {
				if present[target] {
					return domain.MakefileTargetPresent, "", nil
				}
				return domain.MakefileTargetMissing, "", nil
			},
		},
		ProjectRepository: &mockrepo.MockProjectRepository{MarkerFunc: func(directory string) (*domain.ProjectMarkerModel, error) {
			return &domain.ProjectMarkerModel{Platform: platform, App: app}, nil
		}},
		BuildPlatformRepository: &mockrepo.MockBuildPlatformRepository{AllFunc: func() ([]domain.BuildPlatformModel, error) {
			if platforms == nil {
				return nil, errors.New("no server")
			}
			return platforms, nil
		}},
	}
}

// targetNames returns the names of the report's targets.
func targetNames(report *domain.MakefileReportModel) []string {
	names := []string{}
	for _, target := range report.Targets {
		names = append(names, target.Name)
	}
	return names
}

func TestMakefileService_Check_RequiresWebAndMacReleaseFromTheKinds(t *testing.T) {
	present := map[string]bool{"build": true, "run": true, "check": true, "clean": true, "web": true}

	report, err := makefileService(present, "ebitengine", "", makefilePlatforms()).Check(".", "", "")

	assert.NoError(t, err)
	assert.Equal(t, []string{"build", "run", "check", "clean", "web", "mac-release"}, targetNames(report))
	assert.Equal(t, []string{"mac-release"}, report.Broken())
	assert.True(t, report.PlatformKnown)
}

func TestMakefileService_Check_TakesTheKindsOfTheAppVariant(t *testing.T) {
	service := makefileService(map[string]bool{}, "golang", "", makefilePlatforms())

	cli, err := service.Check(".", "", "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"build", "run", "check", "clean"}, targetNames(cli))
	assert.Equal(t, "cli", cli.App)

	gui, err := service.Check(".", "golang", "gui")
	assert.NoError(t, err)
	assert.Equal(t, []string{"build", "run", "check", "clean", "web", "mac-release"}, targetNames(gui))
}

func TestMakefileService_Check_OnlyTheContractWithoutAServer(t *testing.T) {
	report, err := makefileService(map[string]bool{"build": true, "run": true, "check": true, "clean": true}, "ebitengine", "", nil).Check(".", "", "")

	assert.NoError(t, err)
	assert.Equal(t, []string{"build", "run", "check", "clean"}, targetNames(report))
	assert.Empty(t, report.Broken())
	assert.False(t, report.PlatformKnown)
	assert.Contains(t, report.Note, "not checked")
}

func TestMakefileService_Check_RefusesAnUnknownPlatformOrApp(t *testing.T) {
	_, err := makefileService(nil, "nope", "", makefilePlatforms()).Check(".", "", "")
	assert.Error(t, err)
	_, err = makefileService(nil, "c64", "gui", makefilePlatforms()).Check(".", "", "")
	assert.Error(t, err)
}

func TestMakefileService_Check_NeedsAMakefileAndAPlatform(t *testing.T) {
	service := makefileService(nil, "", "", makefilePlatforms())

	_, err := service.Check(".", "", "")
	assert.ErrorContains(t, err, "--platform")

	service.MakefileRepository = &mockrepo.MockMakefileRepository{}
	_, err = service.Check(".", "c64", "")
	assert.ErrorIs(t, err, domain.ErrNoMakefile)
}
