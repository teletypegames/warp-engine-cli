package domain

import (
	"errors"
	"fmt"
	"strings"
)

// MakefileServiceInterface checks a project's Makefile against the shared
// contract of the WarpEngine templates.
type MakefileServiceInterface interface {
	// Check reports the targets the project in directory must have and
	// whether it has them. platform and app override what the project's
	// CI marker names.
	Check(directory string, platform string, app string) (*MakefileReportModel, error)
}

// MakefileService is the [MakefileServiceInterface] over a
// [MakefileRepositoryInterface], a [ProjectRepositoryInterface] and the
// server's build description.
type MakefileService struct {
	// MakefileRepository asks make.
	MakefileRepository MakefileRepositoryInterface
	// ProjectRepository reads the CI marker.
	ProjectRepository ProjectRepositoryInterface
	// BuildPlatformRepository reads the server's platforms and their kinds.
	BuildPlatformRepository BuildPlatformRepositoryInterface
}

// ErrNoMakefile is returned by [MakefileService.Check] for a directory
// without a Makefile.
var ErrNoMakefile = errors.New("no Makefile in the directory")

// Check requires [MakefileContract] of every project; from the server's
// platforms, [MakefileWebTarget] when the release of the platform (its app
// variant, the default one without app) has an html kind, and
// [MakefileManualTarget] when it has a kind the CI does not build. When
// the platforms cannot be read, only the contract is checked and the
// report says so. Every target is probed with make -n.
func (s MakefileService) Check(directory string, platform string, app string) (*MakefileReportModel, error) {
	if directory == "" {
		directory = "."
	}
	if !s.MakefileRepository.Exists(directory) {
		return nil, fmt.Errorf("%s: %w", directory, ErrNoMakefile)
	}
	if platform == "" {
		marker, err := s.ProjectRepository.Marker(directory)
		if err != nil {
			return nil, err
		}
		platform = marker.Platform
		if app == "" {
			app = marker.App
		}
	}
	if platform == "" {
		return nil, errors.New("the project names no platform (.woodpecker.yaml \"platform: …\"): pass --platform")
	}

	report := &MakefileReportModel{Directory: directory, Platform: platform, App: app}
	required := make([]MakefileTargetModel, 0, len(MakefileContract)+2)
	for _, name := range MakefileContract {
		required = append(required, MakefileTargetModel{Name: name, Reason: "every project"})
	}
	platforms, err := s.BuildPlatformRepository.All()
	if err != nil || len(platforms) == 0 {
		report.Note = "the server's platforms could not be read, so web and mac-release were not checked"
		if err != nil {
			report.Note += ": " + err.Error()
		}
	} else {
		kinds, variant, err := releaseKinds(platforms, platform, app)
		if err != nil {
			return nil, err
		}
		report.PlatformKnown = true
		report.App = variant
		manual := []string{}
		web := false
		for _, kind := range kinds {
			if kind.Name == "html" {
				web = true
			}
			if !kind.Built {
				manual = append(manual, kind.Name)
			}
		}
		if web {
			required = append(required, MakefileTargetModel{Name: MakefileWebTarget, Reason: "the platform has a web build (html)"})
		}
		if len(manual) > 0 {
			required = append(required, MakefileTargetModel{Name: MakefileManualTarget, Reason: "the CI does not build " + strings.Join(manual, ", ")})
		}
	}

	for _, target := range required {
		status, detail, err := s.MakefileRepository.Probe(directory, target.Name)
		if err != nil {
			return nil, err
		}
		target.Status, target.Detail = status, detail
		report.Targets = append(report.Targets, target)
	}
	return report, nil
}

// releaseKinds returns the kinds of a release of platform: those of its app
// variant (the default one when app is empty) when it has variants, its own
// otherwise, and the variant they were taken from.
func releaseKinds(platforms []BuildPlatformModel, platform string, app string) ([]BuildKindModel, string, error) {
	names := make([]string, 0, len(platforms))
	for _, known := range platforms {
		names = append(names, known.Name)
		if known.Name != platform {
			continue
		}
		if len(known.Apps) == 0 {
			if app != "" {
				return nil, "", &ValueError{Field: "app", Value: app, Reason: "the platform " + platform + " has no app variants"}
			}
			return known.Kinds, "", nil
		}
		variants := make([]string, 0, len(known.Apps))
		for _, variant := range known.Apps {
			variants = append(variants, variant.Name)
			if variant.Name == app || (app == "" && variant.Default) {
				return variant.Kinds, variant.Name, nil
			}
		}
		if app == "" {
			return known.Kinds, "", nil
		}
		return nil, "", &ValueError{Field: "app", Value: app, Allowed: variants}
	}
	return nil, "", &ValueError{Field: "platform", Value: platform, Allowed: names}
}
