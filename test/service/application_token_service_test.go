package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestApplicationTokenService_Create_NeedsNameAndScopes(t *testing.T) {
	repo := &mockrepo.MockApplicationTokenRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}
	svc := domain.ApplicationTokenService{ApplicationTokenRepository: repo}

	_, err := svc.Create(domain.Attributes{"scopes": []string{"admin:read"}})
	assert.ErrorContains(t, err, "name")

	_, err = svc.Create(domain.Attributes{"name": "ci", "scopes": []string{}})
	assert.ErrorContains(t, err, "scope")
}

func TestApplicationTokenService_Create_TakesScopesFromData(t *testing.T) {
	called := false
	repo := &mockrepo.MockApplicationTokenRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
			called = true
			return &domain.ApplicationTokenModel{ID: 5}, nil
		},
	}
	svc := domain.ApplicationTokenService{ApplicationTokenRepository: repo}

	_, err := svc.Create(domain.Attributes{"name": "ci", "scopes": []any{"admin:read"}})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestApplicationTokenService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.ApplicationTokenService{ApplicationTokenRepository: &mockrepo.MockApplicationTokenRepository{}}

	_, err := svc.Update(5, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestApplicationTokenService_Rotate_DelegatesTheID(t *testing.T) {
	var got_id int
	plain := "new-secret"
	repo := &mockrepo.MockApplicationTokenRepository{
		RotateFunc: func(application_token_id int) (*domain.ApplicationTokenModel, error) {
			got_id = application_token_id
			return &domain.ApplicationTokenModel{ID: 6, PlainToken: &plain}, nil
		},
	}
	svc := domain.ApplicationTokenService{ApplicationTokenRepository: repo}

	token, err := svc.Rotate(5)

	assert.NoError(t, err)
	assert.Equal(t, 5, got_id)
	assert.Equal(t, "new-secret", *token.PlainToken)
}
