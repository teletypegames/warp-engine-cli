package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockProfileRepository mocks [domain.ProfileRepositoryInterface] in
// memory: Load returns Config (an empty one when nil), Save stores into it.
type MockProfileRepository struct {
	// Config is the configuration Load returns and Save replaces.
	Config *domain.ConfigModel
	// Saves counts the Save calls.
	Saves int
	// LoadFunc, when set, replaces Load.
	LoadFunc func() (*domain.ConfigModel, error)
}

// Load returns LoadFunc's answer, or a copy of Config.
func (m *MockProfileRepository) Load() (*domain.ConfigModel, error) {
	if m.LoadFunc != nil {
		return m.LoadFunc()
	}
	config := &domain.ConfigModel{Profiles: map[string]domain.ProfileModel{}}
	if m.Config != nil {
		config.Current = m.Config.Current
		for name, profile := range m.Config.Profiles {
			config.Profiles[name] = profile
		}
	}
	return config, nil
}

// Save stores config in Config.
func (m *MockProfileRepository) Save(config *domain.ConfigModel) error {
	m.Config = config
	m.Saves++
	return nil
}
