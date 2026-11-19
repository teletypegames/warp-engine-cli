package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareImageService_All_ResolvesTheSoftwareName(t *testing.T) {
	var got_id int
	images := &mockrepo.MockSoftwareImageRepository{
		AllFunc: func(software_id int) ([]domain.SoftwareImageModel, error) {
			got_id = software_id
			return nil, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) { return &domain.SoftwareModel{ID: 6}, nil },
	}
	svc := domain.SoftwareImageService{SoftwareImageRepository: images, SoftwareRepository: softwares}

	_, err := svc.All("pong")

	assert.NoError(t, err)
	assert.Equal(t, 6, got_id)
}

func TestSoftwareImageService_Create_NeedsAnImageID(t *testing.T) {
	svc := domain.SoftwareImageService{SoftwareImageRepository: &mockrepo.MockSoftwareImageRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Create("3", domain.Attributes{"isDefault": true})

	assert.ErrorContains(t, err, "imageId")
}

func TestSoftwareImageService_Upload_PassesFileAndAttributes(t *testing.T) {
	var got_id int
	var got_path string
	var got domain.Attributes
	images := &mockrepo.MockSoftwareImageRepository{
		UploadFunc: func(software_id int, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
			got_id, got_path, got = software_id, file_path, attributes
			return &domain.SoftwareImageModel{ID: 31}, nil
		},
	}
	svc := domain.SoftwareImageService{SoftwareImageRepository: images, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	image, err := svc.Upload("3", "cover.png", domain.Attributes{"isDefault": true})

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, "cover.png", got_path)
	assert.Equal(t, domain.Attributes{"isDefault": true}, got)
	assert.Equal(t, 31, image.ID)
}

func TestSoftwareImageService_Upload_NeedsAFile(t *testing.T) {
	svc := domain.SoftwareImageService{SoftwareImageRepository: &mockrepo.MockSoftwareImageRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Upload("3", "", nil)

	assert.ErrorContains(t, err, "file")
}

func TestSoftwareImageService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SoftwareImageService{SoftwareImageRepository: &mockrepo.MockSoftwareImageRepository{}}

	_, err := svc.Update("31", domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSoftwareImageService_Destroy_DelegatesTheParsedID(t *testing.T) {
	var got_id int
	svc := domain.SoftwareImageService{SoftwareImageRepository: &mockrepo.MockSoftwareImageRepository{
		DestroyFunc: func(image_id int) error { got_id = image_id; return nil },
	}}

	err := svc.Destroy("31")

	assert.NoError(t, err)
	assert.Equal(t, 31, got_id)
}
