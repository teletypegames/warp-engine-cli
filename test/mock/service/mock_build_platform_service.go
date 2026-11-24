package mocksvc

// MockBuildPlatformService mocks [domain.BuildPlatformServiceInterface].
type MockBuildPlatformService struct {
	// PlatformsFunc replaces Platforms.
	PlatformsFunc func() []string
	// AssetKindsFunc replaces AssetKinds.
	AssetKindsFunc func() []string
	// AppsFunc replaces Apps.
	AppsFunc func(platform string) ([]string, bool)
	// CheckPlatformFunc replaces CheckPlatform.
	CheckPlatformFunc func(value string, extra ...string) error
	// CheckAssetKindFunc replaces CheckAssetKind.
	CheckAssetKindFunc func(value string) error
	// CheckAppFunc replaces CheckApp.
	CheckAppFunc func(platform string, app string) error
}

// Platforms calls PlatformsFunc, or returns nil.
func (m *MockBuildPlatformService) Platforms() []string {
	if m.PlatformsFunc != nil {
		return m.PlatformsFunc()
	}
	return nil
}

// AssetKinds calls AssetKindsFunc, or returns nil.
func (m *MockBuildPlatformService) AssetKinds() []string {
	if m.AssetKindsFunc != nil {
		return m.AssetKindsFunc()
	}
	return nil
}

// Apps calls AppsFunc, or returns nil and false.
func (m *MockBuildPlatformService) Apps(platform string) ([]string, bool) {
	if m.AppsFunc != nil {
		return m.AppsFunc(platform)
	}
	return nil, false
}

// CheckPlatform calls CheckPlatformFunc, or returns nil.
func (m *MockBuildPlatformService) CheckPlatform(value string, extra ...string) error {
	if m.CheckPlatformFunc != nil {
		return m.CheckPlatformFunc(value, extra...)
	}
	return nil
}

// CheckAssetKind calls CheckAssetKindFunc, or returns nil.
func (m *MockBuildPlatformService) CheckAssetKind(value string) error {
	if m.CheckAssetKindFunc != nil {
		return m.CheckAssetKindFunc(value)
	}
	return nil
}

// CheckApp calls CheckAppFunc, or returns nil.
func (m *MockBuildPlatformService) CheckApp(platform string, app string) error {
	if m.CheckAppFunc != nil {
		return m.CheckAppFunc(platform, app)
	}
	return nil
}
