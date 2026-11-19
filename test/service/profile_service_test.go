package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestProfileService_All_SortedByName(t *testing.T) {
	repo := &mockrepo.MockProfileRepository{Config: config("b", domain.ProfileModel{Name: "b"}, domain.ProfileModel{Name: "a"})}

	profiles, current, err := domain.ProfileService{ProfileRepository: repo}.All()

	assert.NoError(t, err)
	assert.Equal(t, "b", current)
	assert.Equal(t, "a", (*profiles)[0].Name)
}

func TestProfileService_Show_EmptyNameIsTheCurrent(t *testing.T) {
	repo := &mockrepo.MockProfileRepository{Config: config("b", domain.ProfileModel{Name: "b", URL: "u"})}

	profile, err := domain.ProfileService{ProfileRepository: repo}.Show("")

	assert.NoError(t, err)
	assert.Equal(t, "u", profile.URL)
}

func TestProfileService_Use_UnknownIsAnError(t *testing.T) {
	repo := &mockrepo.MockProfileRepository{}

	err := domain.ProfileService{ProfileRepository: repo}.Use("nope")

	assert.Error(t, err)
	assert.Equal(t, 0, repo.Saves)
}

func TestProfileService_Destroy_ClearsTheCurrent(t *testing.T) {
	repo := &mockrepo.MockProfileRepository{Config: config("a", domain.ProfileModel{Name: "a"})}

	err := domain.ProfileService{ProfileRepository: repo}.Destroy("a")

	assert.NoError(t, err)
	assert.Empty(t, repo.Config.Current)
	assert.Empty(t, repo.Config.Profiles)
}

func TestProfileService_Save_FirstProfileBecomesCurrent(t *testing.T) {
	repo := &mockrepo.MockProfileRepository{}

	err := domain.ProfileService{ProfileRepository: repo}.Save(domain.ProfileModel{Name: "a"}, false)

	assert.NoError(t, err)
	assert.Equal(t, "a", repo.Config.Current)
}
