package domain

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// BuildPlatformServiceInterface knows the platforms, asset kinds and app
// variants of the server, read from its build description once per run,
// and checks values against them. When the description cannot be read,
// the checks fall back to the shape of a value only and leave the rest to
// the server.
type BuildPlatformServiceInterface interface {
	// Platforms returns the platform names of the server, sorted, or nil
	// when they cannot be read.
	Platforms() []string
	// AssetKinds returns every asset kind of every platform and app variant
	// of the server, sorted, or nil when they cannot be read.
	AssetKinds() []string
	// Apps returns the app variants of platform, and false when the
	// platforms cannot be read or platform is not one of them.
	Apps(platform string) ([]string, bool)
	// CheckPlatform returns a [*ValueError] when value is not a platform of
	// the server nor one of extra.
	CheckPlatform(value string, extra ...string) error
	// CheckAssetKind returns a [*ValueError] when value is not an asset kind
	// of the server.
	CheckAssetKind(value string) error
	// CheckApp returns a [*ValueError] when app is not a variant of platform.
	CheckApp(platform string, app string) error
}

// PlatformPattern is the shape of a platform name, checked when the
// server's platforms cannot be read.
var PlatformPattern = regexp.MustCompile(`^[a-z0-9_]{2,32}$`)

// AssetKindPattern is the shape of an asset kind or an app variant, checked
// when the server's platforms cannot be read.
var AssetKindPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// ValueError is a value the server would refuse: not one of the platforms,
// asset kinds or app variants it knows, or not even of their shape.
type ValueError struct {
	// Field is the attribute ("platform", "kind", "app").
	Field string
	// Value is what was given.
	Value string
	// Allowed are the values the server knows; empty when only the shape
	// was checked.
	Allowed []string
	// Reason, when set, replaces the list of allowed values in the message.
	Reason string
}

// Error names the field, the value and what it may be.
func (e *ValueError) Error() string {
	switch {
	case e.Reason != "":
		return fmt.Sprintf("%s %q: %s", e.Field, e.Value, e.Reason)
	case len(e.Allowed) > 0:
		return fmt.Sprintf("%s takes one of %s, not %q", e.Field, strings.Join(e.Allowed, ", "), e.Value)
	}
	return fmt.Sprintf("%s %q is no valid %s name", e.Field, e.Value, e.Field)
}

// BuildPlatformCache keeps the platforms read once per run; share one
// between the copies of a [BuildPlatformService].
type BuildPlatformCache struct {
	once      sync.Once
	platforms []BuildPlatformModel
}

// BuildPlatformService is the [BuildPlatformServiceInterface] over a
// [BuildPlatformRepositoryInterface].
type BuildPlatformService struct {
	// BuildPlatformRepository reads the server's platforms.
	BuildPlatformRepository BuildPlatformRepositoryInterface
	// Cache keeps them for the run; nil reads them on every call.
	Cache *BuildPlatformCache
}

// load returns the server's platforms, read once when there is a Cache;
// nil when they cannot be read or there are none.
func (s BuildPlatformService) load() []BuildPlatformModel {
	read := func() []BuildPlatformModel {
		platforms, err := s.BuildPlatformRepository.All()
		if err != nil {
			return nil
		}
		return platforms
	}
	if s.Cache == nil {
		return read()
	}
	s.Cache.once.Do(func() { s.Cache.platforms = read() })
	return s.Cache.platforms
}

// Platforms returns the server's platform names, sorted.
func (s BuildPlatformService) Platforms() []string {
	platforms := s.load()
	if len(platforms) == 0 {
		return nil
	}
	names := make([]string, 0, len(platforms))
	for _, platform := range platforms {
		names = append(names, platform.Name)
	}
	sort.Strings(names)
	return names
}

// AssetKinds returns the union of the kinds of every platform and app
// variant, sorted.
func (s BuildPlatformService) AssetKinds() []string {
	seen := map[string]bool{}
	add := func(kinds []BuildKindModel) {
		for _, kind := range kinds {
			seen[kind.Name] = true
		}
	}
	for _, platform := range s.load() {
		add(platform.Kinds)
		for _, app := range platform.Apps {
			add(app.Kinds)
		}
	}
	if len(seen) == 0 {
		return nil
	}
	kinds := make([]string, 0, len(seen))
	for kind := range seen {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	return kinds
}

// Apps returns the variants of platform in the server's order.
func (s BuildPlatformService) Apps(platform string) ([]string, bool) {
	for _, known := range s.load() {
		if known.Name == platform {
			apps := []string{}
			for _, app := range known.Apps {
				apps = append(apps, app.Name)
			}
			return apps, true
		}
	}
	return nil, false
}

// CheckPlatform checks value against the server's platforms and extra, or
// only its shape ([PlatformPattern]) when they cannot be read.
func (s BuildPlatformService) CheckPlatform(value string, extra ...string) error {
	platforms := s.Platforms()
	if platforms == nil {
		if PlatformPattern.MatchString(value) {
			return nil
		}
		return &ValueError{Field: "platform", Value: value}
	}
	return checkOneOf("platform", value, append(platforms, extra...))
}

// CheckAssetKind checks value against [BuildPlatformService.AssetKinds], or
// only its shape ([AssetKindPattern]) when they cannot be read.
func (s BuildPlatformService) CheckAssetKind(value string) error {
	kinds := s.AssetKinds()
	if kinds == nil {
		if AssetKindPattern.MatchString(value) {
			return nil
		}
		return &ValueError{Field: "kind", Value: value}
	}
	return checkOneOf("kind", value, kinds)
}

// CheckApp checks app against the variants of platform; when the platforms
// cannot be read or platform is unknown, only the shape of app
// ([AssetKindPattern]) is checked.
func (s BuildPlatformService) CheckApp(platform string, app string) error {
	apps, ok := s.Apps(platform)
	if !ok {
		if AssetKindPattern.MatchString(app) {
			return nil
		}
		return &ValueError{Field: "app", Value: app}
	}
	if len(apps) == 0 {
		return &ValueError{Field: "app", Value: app, Reason: "the platform " + platform + " has no app variants"}
	}
	return checkOneOf("app", app, apps)
}

// checkOneOf returns a [*ValueError] when value is not in allowed.
func checkOneOf(field string, value string, allowed []string) error {
	for _, candidate := range allowed {
		if candidate == value {
			return nil
		}
	}
	return &ValueError{Field: field, Value: value, Allowed: allowed}
}

// checkAttribute runs check on attributes[key] when it is a non-empty
// string; any other value is left to the server.
func checkAttribute(attributes Attributes, key string, check func(string) error) error {
	text, ok := attributes[key].(string)
	if !ok || text == "" {
		return nil
	}
	return check(text)
}

// checkPlatformAttribute checks the platform of attributes with platforms,
// allowing extra too; a nil platforms checks nothing.
func checkPlatformAttribute(platforms BuildPlatformServiceInterface, attributes Attributes, extra ...string) error {
	if platforms == nil {
		return nil
	}
	return checkAttribute(attributes, "platform", func(text string) error { return platforms.CheckPlatform(text, extra...) })
}
