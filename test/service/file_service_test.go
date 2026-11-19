package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestFileService_All_TrimsSlashes(t *testing.T) {
	var got string
	repo := &mockrepo.MockFileRepository{AllFunc: func(dir string) (*domain.PageModel[domain.FileEntryModel], error) {
		got = dir
		return &domain.PageModel[domain.FileEntryModel]{}, nil
	}}

	_, err := domain.FileService{FileRepository: repo}.All("/games/pong/")

	assert.NoError(t, err)
	assert.Equal(t, "games/pong", got)
}

func TestFileService_Upload_PassesDirAndFile(t *testing.T) {
	repo := &mockrepo.MockFileRepository{UploadFunc: func(dir string, file_path string) (*domain.FilePathModel, error) {
		assert.Equal(t, "docs", dir)
		assert.Equal(t, "./manual.pdf", file_path)
		return &domain.FilePathModel{Path: "docs/manual.pdf"}, nil
	}}

	result, err := domain.FileService{FileRepository: repo}.Upload("./manual.pdf", "docs/")

	assert.NoError(t, err)
	assert.Equal(t, "docs/manual.pdf", result.Path)
}

func TestFileService_Rename_RefusesAPath(t *testing.T) {
	repo := &mockrepo.MockFileRepository{RenameFunc: func(path string, new_name string) (*domain.FilePathModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.FileService{FileRepository: repo}

	_, err := svc.Rename("docs/a.pdf", "other/b.pdf")
	assert.ErrorContains(t, err, "without \"/\"")

	_, err = svc.Rename("docs/a.pdf", "")
	assert.Error(t, err)

	_, err = svc.Rename("/", "b.pdf")
	assert.Error(t, err)
}

func TestFileService_Destroy_RootIsAnError(t *testing.T) {
	called := false
	repo := &mockrepo.MockFileRepository{DestroyFunc: func(path string) error { called = true; return nil }}

	err := domain.FileService{FileRepository: repo}.Destroy("/")

	assert.Error(t, err)
	assert.False(t, called)
}

func TestFileService_CreateDirectory_SplitsParentAndName(t *testing.T) {
	var got [][2]string
	repo := &mockrepo.MockFileRepository{CreateDirectoryFunc: func(dir string, name string) (*domain.FilePathModel, error) {
		got = append(got, [2]string{dir, name})
		return &domain.FilePathModel{}, nil
	}}
	svc := domain.FileService{FileRepository: repo}

	_, err := svc.CreateDirectory("games/pong/docs/")
	assert.NoError(t, err)
	_, err = svc.CreateDirectory("top")
	assert.NoError(t, err)
	_, err = svc.CreateDirectory("/")
	assert.Error(t, err)

	assert.Equal(t, [][2]string{{"games/pong", "docs"}, {"", "top"}}, got)
}
