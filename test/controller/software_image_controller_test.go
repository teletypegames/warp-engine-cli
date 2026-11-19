package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareImageController_Index_PrintsTheImages(t *testing.T) {
	svc := &mocksvc.MockSoftwareImageService{
		AllFunc: func(reference string) ([]domain.SoftwareImageModel, error) {
			return []domain.SoftwareImageModel{{ID: 31, ImageID: 118, IsDefault: true, Origin: "local", URL: text("https://img.example.org/118.png")}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareImageController: cli.SoftwareImageController{SoftwareImageService: svc}}, "image", "list", "pong")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "118.png")
	assert.Contains(t, stdout, "yes")
}

func TestSoftwareImageController_Upload_PassesFileDefaultAndPosition(t *testing.T) {
	var got_reference, got_path string
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareImageService{
		UploadFunc: func(reference string, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
			got_reference, got_path, got = reference, file_path, attributes
			return &domain.SoftwareImageModel{ID: 32}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareImageController: cli.SoftwareImageController{SoftwareImageService: svc}},
		"image", "upload", "pong", "cover.png", "--default", "--position", "2")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Equal(t, "cover.png", got_path)
	assert.Equal(t, domain.Attributes{"isDefault": true, "position": 2}, got)
}

func TestSoftwareImageController_Create_SendsTheImageID(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareImageService{
		CreateFunc: func(reference string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
			got = attributes
			return &domain.SoftwareImageModel{ID: 33}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareImageController: cli.SoftwareImageController{SoftwareImageService: svc}},
		"image", "create", "pong", "--image-id", "118")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"imageId": 118}, got)
}

func TestSoftwareImageController_Update_RejectsANonNumericPosition(t *testing.T) {
	_, _, err := execute(cli.Router{SoftwareImageController: cli.SoftwareImageController{SoftwareImageService: &mocksvc.MockSoftwareImageService{}}},
		"image", "update", "31", "--position", "first")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSoftwareImageController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareImageService{DestroyFunc: func(id string) error { called = true; return nil }}
	router := cli.Router{SoftwareImageController: cli.SoftwareImageController{SoftwareImageService: svc}}

	_, _, err := execute(router, "image", "delete", "31")
	assert.Error(t, err)
	assert.False(t, called)

	stdout, _, err := execute(router, "image", "delete", "31", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
	assert.Contains(t, stdout, "unbound")
}
