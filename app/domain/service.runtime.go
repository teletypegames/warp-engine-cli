package domain

import "errors"

// RuntimeServiceInterface manages the installed runtimes: the players a
// platform's games run in on the web.
type RuntimeServiceInterface interface {
	// All lists the installed runtime versions matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[RuntimeVersionModel], error)
	// Show returns the runtime version runtime_id.
	Show(runtime_id int) (*RuntimeVersionModel, error)
	// Install installs version (empty: the chosen one) of the downloaded
	// runtime name; with dry_run it only reports what it would fetch.
	Install(name string, version string, dry_run bool) (*RuntimeInstallModel, error)
}

// RuntimeService is the [RuntimeServiceInterface] over a
// [RuntimeRepositoryInterface].
type RuntimeService struct {
	// RuntimeRepository talks to the runtimes resource.
	RuntimeRepository RuntimeRepositoryInterface
}

// All lists the runtime versions.
func (s RuntimeService) All(filter Filter, query ListQuery) (*PageModel[RuntimeVersionModel], error) {
	return s.RuntimeRepository.All(filter, query)
}

// Show gets the runtime version.
func (s RuntimeService) Show(runtime_id int) (*RuntimeVersionModel, error) {
	return s.RuntimeRepository.FindByID(runtime_id)
}

// Install checks that name is given and installs it; an empty version is
// not sent, so the server takes the chosen one.
func (s RuntimeService) Install(name string, version string, dry_run bool) (*RuntimeInstallModel, error) {
	if name == "" {
		return nil, errors.New("no runtime given")
	}
	attributes := Attributes{"name": name, "dryRun": dry_run}
	if version != "" {
		attributes["version"] = version
	}
	return s.RuntimeRepository.Install(attributes)
}
