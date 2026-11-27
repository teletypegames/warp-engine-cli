package domain

import "errors"

// LocalProviderServiceInterface manages the providers of WarpEngine
// Local (the git server, the CI) and the local CI of a software. The
// tasks that start or stop containers run in the background on the
// server and answer a queued [LocalTaskModel]; follow it with
// [LocalTaskServiceInterface.Wait].
type LocalProviderServiceInterface interface {
	// All lists the providers with their state.
	All() (*PageModel[LocalProviderModel], error)
	// Enable starts the provider key and what it needs.
	Enable(key string) (*LocalTaskModel, error)
	// Disable stops the provider key and what needs it; its data stays.
	Disable(key string) (*LocalTaskModel, error)
	// Remove removes the containers of the provider key; with purge its
	// data and generated sign-ins too.
	Remove(key string, purge bool) (*LocalTaskModel, error)
	// Credentials returns the generated sign-in of the provider key.
	Credentials(key string) (*LocalCredentialsModel, error)
	// SetupCI sets the software reference up on the local git server and CI.
	SetupCI(reference string) (*LocalTaskModel, error)
	// TeardownCI stops the local CI of the software reference.
	TeardownCI(reference string) error
}

// LocalProviderService is the [LocalProviderServiceInterface] over a
// [LocalProviderRepositoryInterface]; it resolves software references
// through the [SoftwareRepositoryInterface].
type LocalProviderService struct {
	// LocalProviderRepository talks to the local_providers resource.
	LocalProviderRepository LocalProviderRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// errNoProvider is returned for an empty provider key.
var errNoProvider = errors.New("no provider given")

// All lists the providers.
func (s LocalProviderService) All() (*PageModel[LocalProviderModel], error) {
	return s.LocalProviderRepository.All()
}

// Enable queues enabling key.
func (s LocalProviderService) Enable(key string) (*LocalTaskModel, error) {
	if key == "" {
		return nil, errNoProvider
	}
	return s.LocalProviderRepository.Enable(key)
}

// Disable queues disabling key.
func (s LocalProviderService) Disable(key string) (*LocalTaskModel, error) {
	if key == "" {
		return nil, errNoProvider
	}
	return s.LocalProviderRepository.Disable(key)
}

// Remove queues removing key (purging it with purge).
func (s LocalProviderService) Remove(key string, purge bool) (*LocalTaskModel, error) {
	if key == "" {
		return nil, errNoProvider
	}
	return s.LocalProviderRepository.Destroy(key, purge)
}

// Credentials gets the sign-in of key.
func (s LocalProviderService) Credentials(key string) (*LocalCredentialsModel, error) {
	if key == "" {
		return nil, errNoProvider
	}
	return s.LocalProviderRepository.Credentials(key)
}

// SetupCI resolves reference and queues its local CI setup.
func (s LocalProviderService) SetupCI(reference string) (*LocalTaskModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.LocalProviderRepository.SetupCI(software_id)
}

// TeardownCI resolves reference and stops its local CI.
func (s LocalProviderService) TeardownCI(reference string) error {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return err
	}
	return s.LocalProviderRepository.TeardownCI(software_id)
}
