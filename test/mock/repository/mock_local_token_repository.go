package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalTokenRepository mocks [domain.LocalTokenRepositoryInterface]:
// Read returns Token, or [domain.ErrNoLocalToken] when it is empty.
type MockLocalTokenRepository struct {
	// Token is the token file's content.
	Token string
}

// Read returns Token or [domain.ErrNoLocalToken].
func (m *MockLocalTokenRepository) Read() (string, error) {
	if m.Token == "" {
		return "", domain.ErrNoLocalToken
	}
	return m.Token, nil
}
