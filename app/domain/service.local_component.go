package domain

import "errors"

// LocalComponentServiceInterface manages the components of WarpEngine
// Local (the git server, the CI) and the local CI of a software. The
// operations that start or stop containers run in the background on the
// server and answer a queued [LocalOperationModel]; follow it with
// [LocalOperationServiceInterface.Wait].
type LocalComponentServiceInterface interface {
	// All lists the components with their state.
	All() (*PageModel[LocalComponentModel], error)
	// Enable starts the component key and what it needs.
	Enable(key string) (*LocalOperationModel, error)
	// Disable stops the component key and what needs it; its data stays.
	Disable(key string) (*LocalOperationModel, error)
	// Remove removes the containers of the component key; with purge its
	// data and generated sign-ins too.
	Remove(key string, purge bool) (*LocalOperationModel, error)
	// Credentials returns the generated sign-in of the component key.
	Credentials(key string) (*LocalCredentialsModel, error)
	// SetupCI sets the software reference up on the local git server and CI.
	SetupCI(reference string) (*LocalOperationModel, error)
	// TeardownCI stops the local CI of the software reference.
	TeardownCI(reference string) error
}

// LocalComponentService is the [LocalComponentServiceInterface] over a
// [LocalComponentRepositoryInterface]; it resolves software references
// through the [SoftwareRepositoryInterface].
type LocalComponentService struct {
	// LocalComponentRepository talks to the local_components resource.
	LocalComponentRepository LocalComponentRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// errNoComponent is returned for an empty component key.
var errNoComponent = errors.New("no component given")

// All lists the components.
func (s LocalComponentService) All() (*PageModel[LocalComponentModel], error) {
	return s.LocalComponentRepository.All()
}

// Enable queues enabling key.
func (s LocalComponentService) Enable(key string) (*LocalOperationModel, error) {
	if key == "" {
		return nil, errNoComponent
	}
	return s.LocalComponentRepository.Enable(key)
}

// Disable queues disabling key.
func (s LocalComponentService) Disable(key string) (*LocalOperationModel, error) {
	if key == "" {
		return nil, errNoComponent
	}
	return s.LocalComponentRepository.Disable(key)
}

// Remove queues removing key (purging it with purge).
func (s LocalComponentService) Remove(key string, purge bool) (*LocalOperationModel, error) {
	if key == "" {
		return nil, errNoComponent
	}
	return s.LocalComponentRepository.Destroy(key, purge)
}

// Credentials gets the sign-in of key.
func (s LocalComponentService) Credentials(key string) (*LocalCredentialsModel, error) {
	if key == "" {
		return nil, errNoComponent
	}
	return s.LocalComponentRepository.Credentials(key)
}

// SetupCI resolves reference and queues its local CI setup.
func (s LocalComponentService) SetupCI(reference string) (*LocalOperationModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.LocalComponentRepository.SetupCI(software_id)
}

// TeardownCI resolves reference and stops its local CI.
func (s LocalComponentService) TeardownCI(reference string) error {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return err
	}
	return s.LocalComponentRepository.TeardownCI(software_id)
}
