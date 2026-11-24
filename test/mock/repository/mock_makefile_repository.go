package mockrepo

// MockMakefileRepository mocks [domain.MakefileRepositoryInterface].
type MockMakefileRepository struct {
	// ExistsFunc replaces Exists.
	ExistsFunc func(directory string) bool
	// ProbeFunc replaces Probe.
	ProbeFunc func(directory string, target string) (string, string, error)
}

// Exists calls ExistsFunc, or returns false.
func (m *MockMakefileRepository) Exists(directory string) bool {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(directory)
	}
	return false
}

// Probe calls ProbeFunc, or returns empty values.
func (m *MockMakefileRepository) Probe(directory string, target string) (string, string, error) {
	if m.ProbeFunc != nil {
		return m.ProbeFunc(directory, target)
	}
	return "", "", nil
}
