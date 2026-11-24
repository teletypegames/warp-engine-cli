package service_test

import (
	"errors"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// buildPlatforms is a server with a platform without variants and golang
// with cli and gui.
func buildPlatforms() []domain.BuildPlatformModel {
	return []domain.BuildPlatformModel{
		{Name: "tic80", Kinds: []domain.BuildKindModel{{Name: "cartridge"}, {Name: "html"}}},
		{Name: "golang", Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64"}}, Apps: []domain.BuildAppModel{
			{Name: "cli", Default: true, Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64"}}},
			{Name: "gui", Kinds: []domain.BuildKindModel{{Name: "linux_x64"}}},
		}},
	}
}

// buildPlatformService is a [domain.BuildPlatformService] over platforms,
// or over a failing repository when platforms is nil.
func buildPlatformService(platforms []domain.BuildPlatformModel) domain.BuildPlatformService {
	repo := &mockrepo.MockBuildPlatformRepository{AllFunc: func() ([]domain.BuildPlatformModel, error) {
		if platforms == nil {
			return nil, errors.New("connection refused")
		}
		return platforms, nil
	}}
	return domain.BuildPlatformService{BuildPlatformRepository: repo, Cache: &domain.BuildPlatformCache{}}
}

func TestBuildPlatformService_Platforms_ReadsOncePerCache(t *testing.T) {
	calls := 0
	repo := &mockrepo.MockBuildPlatformRepository{AllFunc: func() ([]domain.BuildPlatformModel, error) {
		calls++
		return buildPlatforms(), nil
	}}
	svc := domain.BuildPlatformService{BuildPlatformRepository: repo, Cache: &domain.BuildPlatformCache{}}

	assert.Equal(t, []string{"golang", "tic80"}, svc.Platforms())
	assert.NoError(t, svc.CheckPlatform("golang"))
	assert.Equal(t, 1, calls)
}

func TestBuildPlatformService_AssetKinds_UnionOfPlatformsAndApps(t *testing.T) {
	svc := buildPlatformService(buildPlatforms())

	assert.Equal(t, []string{"cartridge", "cli_linux_x64", "html", "linux_x64"}, svc.AssetKinds())
}

func TestBuildPlatformService_CheckPlatform_Known(t *testing.T) {
	svc := buildPlatformService(buildPlatforms())

	assert.NoError(t, svc.CheckPlatform("tic80"))
	assert.NoError(t, svc.CheckPlatform("unknown", "unknown"))
	var value *domain.ValueError
	err := svc.CheckPlatform("atari")
	assert.ErrorAs(t, err, &value)
	assert.ErrorContains(t, err, "platform takes one of golang, tic80, not \"atari\"")
}

func TestBuildPlatformService_CheckPlatform_FallsBackToTheShape(t *testing.T) {
	svc := buildPlatformService(nil)

	assert.Nil(t, svc.Platforms())
	assert.NoError(t, svc.CheckPlatform("retroarch"))
	assert.Error(t, svc.CheckPlatform("Not A Platform"))
	assert.Error(t, svc.CheckPlatform("x"))
}

func TestBuildPlatformService_CheckAssetKind(t *testing.T) {
	assert.NoError(t, buildPlatformService(buildPlatforms()).CheckAssetKind("linux_x64"))
	assert.Error(t, buildPlatformService(buildPlatforms()).CheckAssetKind("linux_armhf"))
	assert.NoError(t, buildPlatformService(nil).CheckAssetKind("linux_armhf"))
	assert.Error(t, buildPlatformService(nil).CheckAssetKind("linux-armhf"))
}

func TestBuildPlatformService_CheckApp(t *testing.T) {
	svc := buildPlatformService(buildPlatforms())

	assert.NoError(t, svc.CheckApp("golang", "gui"))
	assert.ErrorContains(t, svc.CheckApp("golang", "tui"), "app takes one of cli, gui")
	assert.ErrorContains(t, svc.CheckApp("tic80", "cli"), "has no app variants")
	assert.NoError(t, buildPlatformService(nil).CheckApp("golang", "tui"))
}

func TestSoftwareService_Create_ChecksThePlatform(t *testing.T) {
	repo := &mockrepo.MockSoftwareRepository{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.SoftwareService{SoftwareRepository: repo, BuildPlatformService: buildPlatformService(buildPlatforms())}

	_, err := svc.Create(domain.Attributes{"name": "x", "title": "X", "platform": "atari"})

	var value *domain.ValueError
	assert.ErrorAs(t, err, &value)
}

func TestSoftwareService_Scaffold_ChecksTheApp(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockSoftwareRepository{ScaffoldFunc: func(attributes domain.Attributes) (*domain.ScaffoldModel, error) {
		got = attributes
		return &domain.ScaffoldModel{}, nil
	}}
	svc := domain.SoftwareService{SoftwareRepository: repo, BuildPlatformService: buildPlatformService(buildPlatforms())}

	_, err := svc.Scaffold(domain.Attributes{"name": "x", "title": "X", "platform": "golang", "app": "tui"})
	assert.ErrorContains(t, err, "app takes one of cli, gui")

	_, err = svc.Scaffold(domain.Attributes{"name": "x", "title": "X", "platform": "golang", "app": "cli"})
	assert.NoError(t, err)
	assert.Equal(t, "cli", got["app"])
}

func TestPipelineService_Update_TakesUnknownBesidesThePlatforms(t *testing.T) {
	platforms := &mocksvc.MockBuildPlatformService{CheckPlatformFunc: func(value string, extra ...string) error {
		assert.Equal(t, []string{domain.PipelineUnknownPlatform}, extra)
		return nil
	}}
	svc := domain.PipelineService{PipelineRepository: &mockrepo.MockPipelineRepository{}, BuildPlatformService: platforms}

	_, err := svc.Update(4, domain.Attributes{"platform": "unknown"})

	assert.NoError(t, err)
}

func TestPlatformLinkService_Update_ChecksThePlatform(t *testing.T) {
	svc := domain.PlatformLinkService{PlatformLinkRepository: &mockrepo.MockPlatformLinkRepository{}, BuildPlatformService: buildPlatformService(buildPlatforms())}

	_, err := svc.Update(3, domain.Attributes{"platform": "atari"})

	assert.Error(t, err)
}

func TestReleaseAssetService_Create_ChecksTheKind(t *testing.T) {
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{}, BuildPlatformService: buildPlatformService(buildPlatforms())}

	_, err := svc.Create("42", domain.Attributes{"kind": "amiga", "path": "/x"})
	assert.Error(t, err)

	_, err = svc.Create("42", domain.Attributes{"kind": "cli_linux_x64", "path": "/x"})
	assert.NoError(t, err)
}
