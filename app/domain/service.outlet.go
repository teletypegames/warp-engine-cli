package domain

import "errors"

// OutletServiceInterface reads the engine's outlets.
type OutletServiceInterface interface {
	// All lists the configured outlets.
	All() ([]OutletModel, error)
	// Check returns what the outlet key offers and leaves out.
	Check(key string) (*OutletCheckModel, error)
}

// OutletService is the [OutletServiceInterface] over an
// [OutletRepositoryInterface].
type OutletService struct {
	// OutletRepository talks to the outlets resource.
	OutletRepository OutletRepositoryInterface
}

// All lists the outlets.
func (s OutletService) All() ([]OutletModel, error) {
	return s.OutletRepository.All()
}

// Check refuses an empty key and gets the outlet's check.
func (s OutletService) Check(key string) (*OutletCheckModel, error) {
	if key == "" {
		return nil, errors.New("the outlet key is empty")
	}
	return s.OutletRepository.Check(key)
}
