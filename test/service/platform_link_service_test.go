package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestPlatformLinkService_Create_NeedsPlatformNameURL(t *testing.T) {
	repo := &mockrepo.MockPlatformLinkRepository{CreateFunc: func(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.PlatformLinkService{PlatformLinkRepository: repo}

	_, err := svc.Create(domain.Attributes{"platform": "c64", "name": "VICE"})

	assert.ErrorContains(t, err, "url")
}

func TestPlatformLinkService_Create_DelegatesAttributes(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockPlatformLinkRepository{CreateFunc: func(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
		got = attributes
		return &domain.PlatformLinkModel{ID: 1}, nil
	}}
	svc := domain.PlatformLinkService{PlatformLinkRepository: repo}
	attributes := domain.Attributes{"platform": "c64", "name": "VICE", "url": "https://vice.example.org"}

	_, err := svc.Create(attributes)

	assert.NoError(t, err)
	assert.Equal(t, attributes, got)
}

func TestPlatformLinkService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.PlatformLinkService{PlatformLinkRepository: &mockrepo.MockPlatformLinkRepository{}}

	_, err := svc.Update(1, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestPlatformLinkService_Destroy_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockPlatformLinkRepository{DestroyFunc: func(platform_link_id int) error { got_id = platform_link_id; return nil }}
	svc := domain.PlatformLinkService{PlatformLinkRepository: repo}

	assert.NoError(t, svc.Destroy(6))
	assert.Equal(t, 6, got_id)
}
