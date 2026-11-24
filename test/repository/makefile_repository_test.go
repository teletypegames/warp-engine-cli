package repository_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakefileRepository_Probe_TellsPresentMissingAndFailing(t *testing.T) {
	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not installed")
	}
	directory := t.TempDir()
	makefile := "build:\n\techo build\n\nrun: missing-file\n\techo run\n"
	assert.NoError(t, os.WriteFile(filepath.Join(directory, "Makefile"), []byte(makefile), 0o600))
	repo := domain.MakefileRepository{}

	assert.True(t, repo.Exists(directory))
	status, _, err := repo.Probe(directory, "build")
	assert.NoError(t, err)
	assert.Equal(t, domain.MakefileTargetPresent, status)
	status, _, err = repo.Probe(directory, "clean")
	assert.NoError(t, err)
	assert.Equal(t, domain.MakefileTargetMissing, status)
	status, detail, err := repo.Probe(directory, "run")
	assert.NoError(t, err)
	assert.Equal(t, domain.MakefileTargetFailing, status)
	assert.Contains(t, detail, "missing-file")
}

func TestMakefileRepository_Exists_FalseWithoutMakefile(t *testing.T) {
	assert.False(t, domain.MakefileRepository{}.Exists(t.TempDir()))
}

func TestProjectRepository_Marker_ReadsTheWoodpeckerMarker(t *testing.T) {
	directory := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(directory, ".woodpecker.yaml"), []byte("platform: golang\napp: gui # the gio variant\n"), 0o600))

	marker, err := domain.ProjectRepository{}.Marker(directory)

	assert.NoError(t, err)
	assert.Equal(t, &domain.ProjectMarkerModel{File: ".woodpecker.yaml", Platform: "golang", App: "gui"}, marker)
}

func TestProjectRepository_Marker_EmptyWithoutMarker(t *testing.T) {
	directory := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(directory, "metadata.json"), []byte(`{"name":"pong"}`), 0o600))

	marker, err := domain.ProjectRepository{}.Marker(directory)

	assert.NoError(t, err)
	assert.Equal(t, "", marker.Platform)
}
