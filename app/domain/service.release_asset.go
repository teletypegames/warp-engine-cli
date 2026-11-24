package domain

import "errors"

// ReleaseAssetServiceInterface manages the files of the releases. Releases
// and assets are named by their IDs.
type ReleaseAssetServiceInterface interface {
	// All lists the assets of the release release.
	All(release string) ([]ReleaseAssetModel, error)
	// Show returns the asset id.
	Show(id string) (*ReleaseAssetModel, error)
	// Create adds an asset to the release release.
	Create(release string, attributes Attributes) (*ReleaseAssetModel, error)
	// Update changes the asset id.
	Update(id string, attributes Attributes) (*ReleaseAssetModel, error)
	// Destroy deletes the asset id.
	Destroy(id string) error
}

// ReleaseAssetService is the [ReleaseAssetServiceInterface] over a
// [ReleaseAssetRepositoryInterface].
type ReleaseAssetService struct {
	// ReleaseAssetRepository talks to the release assets resource.
	ReleaseAssetRepository ReleaseAssetRepositoryInterface
	// BuildPlatformService checks the asset kind
	// against the server's; nil checks nothing.
	BuildPlatformService BuildPlatformServiceInterface
}

// All parses release and lists its assets.
func (s ReleaseAssetService) All(release string) ([]ReleaseAssetModel, error) {
	release_id, err := ParseID("release", release)
	if err != nil {
		return nil, err
	}
	return s.ReleaseAssetRepository.All(release_id)
}

// Show parses id and gets the asset.
func (s ReleaseAssetService) Show(id string) (*ReleaseAssetModel, error) {
	asset_id, err := ParseID("asset", id)
	if err != nil {
		return nil, err
	}
	return s.ReleaseAssetRepository.FindByID(asset_id)
}

// Create checks that attributes hold a kind and a path — the server would
// refuse without them — and that the kind is one of the server's, parses
// release and adds the asset.
func (s ReleaseAssetService) Create(release string, attributes Attributes) (*ReleaseAssetModel, error) {
	for _, key := range []string{"kind", "path"} {
		if attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("an asset needs a " + key)
		}
	}
	if err := s.checkKind(attributes); err != nil {
		return nil, err
	}
	release_id, err := ParseID("release", release)
	if err != nil {
		return nil, err
	}
	return s.ReleaseAssetRepository.Create(release_id, attributes)
}

// Update checks a given kind against the server's, parses id and changes
// the asset; empty attributes are [ErrNothingToChange].
func (s ReleaseAssetService) Update(id string, attributes Attributes) (*ReleaseAssetModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	if err := s.checkKind(attributes); err != nil {
		return nil, err
	}
	asset_id, err := ParseID("asset", id)
	if err != nil {
		return nil, err
	}
	return s.ReleaseAssetRepository.Update(asset_id, attributes)
}

// Destroy parses id and deletes the asset.
func (s ReleaseAssetService) Destroy(id string) error {
	asset_id, err := ParseID("asset", id)
	if err != nil {
		return err
	}
	return s.ReleaseAssetRepository.Destroy(asset_id)
}

// checkKind checks the kind of attributes against the server's asset kinds;
// without a BuildPlatformService it checks nothing.
func (s ReleaseAssetService) checkKind(attributes Attributes) error {
	if s.BuildPlatformService == nil {
		return nil
	}
	return checkAttribute(attributes, "kind", s.BuildPlatformService.CheckAssetKind)
}
