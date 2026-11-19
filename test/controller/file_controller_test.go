package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestFileController_Index_DirIsOptional(t *testing.T) {
	var got []string
	size := int64(2048)
	svc := &mocksvc.MockFileService{AllFunc: func(dir string) (*domain.PageModel[domain.FileEntryModel], error) {
		got = append(got, dir)
		return &domain.PageModel[domain.FileEntryModel]{Items: []domain.FileEntryModel{{Name: "a.prg", Path: "pong/a.prg", Type: "file", Size: &size}}}, nil
	}}
	router := cli.Router{FileController: cli.FileController{FileService: svc}}

	stdout, _, err := execute(router, "file", "list")
	assert.NoError(t, err)
	_, _, err = execute(router, "files", "ls", "pong")
	assert.NoError(t, err)

	assert.Equal(t, []string{"", "pong"}, got)
	assert.Contains(t, stdout, "2.0 KiB")
}

func TestFileController_Upload_PassesFileAndDir(t *testing.T) {
	svc := &mocksvc.MockFileService{UploadFunc: func(file_path string, dir string) (*domain.FilePathModel, error) {
		assert.Equal(t, "./manual.pdf", file_path)
		assert.Equal(t, "docs", dir)
		return &domain.FilePathModel{Path: "docs/manual.pdf"}, nil
	}}

	stdout, _, err := execute(cli.Router{FileController: cli.FileController{FileService: svc}}, "file", "upload", "./manual.pdf", "--dir", "docs")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "docs/manual.pdf")
}

func TestFileController_Rename_PassesPathAndName(t *testing.T) {
	svc := &mocksvc.MockFileService{RenameFunc: func(file_path string, new_name string) (*domain.FilePathModel, error) {
		assert.Equal(t, "docs/a.pdf", file_path)
		assert.Equal(t, "b.pdf", new_name)
		return &domain.FilePathModel{Path: "docs/b.pdf"}, nil
	}}

	stdout, _, err := execute(cli.Router{FileController: cli.FileController{FileService: svc}}, "file", "rename", "docs/a.pdf", "b.pdf", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, "docs/b.pdf\n", stdout)
}

func TestFileController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	var got string
	svc := &mocksvc.MockFileService{DestroyFunc: func(file_path string) error { got = file_path; return nil }}
	router := cli.Router{FileController: cli.FileController{FileService: svc}}

	_, _, err := execute(router, "file", "delete", "docs/a.pdf")
	assert.Error(t, err, "no --yes and no answer: nothing is done")
	assert.Empty(t, got)

	stdout, _, err := execute(router, "file", "rm", "docs/a.pdf", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, "docs/a.pdf", got)
	assert.Contains(t, stdout, "deleted")
}

func TestFileController_CreateDirectory_PassesPath(t *testing.T) {
	var got string
	svc := &mocksvc.MockFileService{CreateDirectoryFunc: func(dir_path string) (*domain.FilePathModel, error) {
		got = dir_path
		return &domain.FilePathModel{Path: dir_path}, nil
	}}

	_, _, err := execute(cli.Router{FileController: cli.FileController{FileService: svc}}, "file", "mkdir", "pong/docs")

	assert.NoError(t, err)
	assert.Equal(t, "pong/docs", got)
}
