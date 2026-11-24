package service_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

// buildFile writes a file named name with content into a temporary
// directory and returns its path.
func buildFile(t *testing.T, name string, content string) string {
	file_path := filepath.Join(t.TempDir(), name)
	assert.NoError(t, os.WriteFile(file_path, []byte(content), 0o600))
	return file_path
}

func TestBuildService_Upload_SendsEveryFileWithItsDigest(t *testing.T) {
	first := buildFile(t, "pong-1.2.0-mac-x64.zip", "a")
	second := buildFile(t, "pong-1.2.0-mac-arm64.zip", "b")
	sent := map[string]string{}
	repo := &mockrepo.MockBuildRepository{UploadFunc: func(target domain.BuildTargetModel, file_path string, sha256 string) (*domain.BuildUploadModel, error) {
		sent[filepath.Base(file_path)] = sha256
		return &domain.BuildUploadModel{File: filepath.Base(file_path)}, nil
	}}

	uploads, err := domain.BuildService{BuildRepository: repo}.Upload(domain.BuildTargetModel{Name: "pong", Version: "1.2.0"}, []string{first, second})

	assert.NoError(t, err)
	assert.Len(t, uploads, 2)
	assert.Equal(t, "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb", sent["pong-1.2.0-mac-x64.zip"])
	assert.Equal(t, "3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d", sent["pong-1.2.0-mac-arm64.zip"])
}

func TestBuildService_Upload_RefusesAFileOfAnotherRelease(t *testing.T) {
	file_path := buildFile(t, "pong-1.1.0-mac-x64.zip", "a")
	repo := &mockrepo.MockBuildRepository{UploadFunc: func(target domain.BuildTargetModel, file_path string, sha256 string) (*domain.BuildUploadModel, error) {
		t.Fatal("nothing may be uploaded")
		return nil, nil
	}}

	_, err := domain.BuildService{BuildRepository: repo}.Upload(domain.BuildTargetModel{Name: "pong", Version: "1.2.0"}, []string{file_path})

	assert.ErrorContains(t, err, "must start with pong-1.2.0")
}

func TestBuildService_Upload_NeedsNameVersionAndFiles(t *testing.T) {
	service := domain.BuildService{BuildRepository: &mockrepo.MockBuildRepository{}}

	_, err := service.Upload(domain.BuildTargetModel{Version: "1"}, []string{"x"})
	assert.ErrorContains(t, err, "--name")
	_, err = service.Upload(domain.BuildTargetModel{Name: "pong"}, []string{"x"})
	assert.ErrorContains(t, err, "--version")
	_, err = service.Upload(domain.BuildTargetModel{Name: "pong", Version: "1"}, nil)
	assert.ErrorIs(t, err, domain.ErrNoBuildFiles)
}

func TestBuildService_Upload_StopsAtTheFirstFailure(t *testing.T) {
	first := buildFile(t, "pong-1-a.zip", "a")
	second := buildFile(t, "pong-1-b.zip", "b")
	calls := 0
	repo := &mockrepo.MockBuildRepository{UploadFunc: func(target domain.BuildTargetModel, file_path string, sha256 string) (*domain.BuildUploadModel, error) {
		calls++
		return nil, errors.New("401")
	}}

	_, err := domain.BuildService{BuildRepository: repo}.Upload(domain.BuildTargetModel{Name: "pong", Version: "1"}, []string{first, second})

	assert.Error(t, err)
	assert.Equal(t, 1, calls)
}

func TestBuildService_Publish_ChecksPlatformAndApp(t *testing.T) {
	published := false
	repo := &mockrepo.MockBuildRepository{PublishFunc: func(target domain.BuildTargetModel) (*domain.BuildPublishModel, error) {
		published = true
		return &domain.BuildPublishModel{Published: true}, nil
	}}
	service := domain.BuildService{BuildRepository: repo, BuildPlatformService: buildPlatformService(buildPlatforms())}

	_, err := service.Publish(domain.BuildTargetModel{Name: "x", Version: "1", Platform: "nope"})
	assert.Error(t, err)
	_, err = service.Publish(domain.BuildTargetModel{Name: "x", Version: "1", Platform: "golang", App: "tui"})
	assert.Error(t, err)
	assert.False(t, published)
	_, err = service.Publish(domain.BuildTargetModel{Name: "x", Version: "1", Platform: "golang", App: "gui"})
	assert.NoError(t, err)
	assert.True(t, published)
}

func TestBuildService_Publish_NeedsAPlatform(t *testing.T) {
	_, err := domain.BuildService{BuildRepository: &mockrepo.MockBuildRepository{}}.Publish(domain.BuildTargetModel{Name: "x", Version: "1"})

	assert.ErrorContains(t, err, "--platform")
}
