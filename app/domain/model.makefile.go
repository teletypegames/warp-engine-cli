package domain

// The states of a [MakefileTargetModel].
const (
	// MakefileTargetPresent: make -n ran the target.
	MakefileTargetPresent = "ok"
	// MakefileTargetMissing: the Makefile has no rule for the target.
	MakefileTargetMissing = "missing"
	// MakefileTargetFailing: the target exists, but make -n failed on it
	// (a missing prerequisite, a broken recipe).
	MakefileTargetFailing = "fails"
)

// MakefileContract are the targets every project's Makefile has, whatever
// its platform: build, run, check (the CI's check step runs it) and clean.
var MakefileContract = []string{"build", "run", "check", "clean"}

// MakefileWebTarget is the target of a platform whose releases have a web
// build (the html kind); the Toolbox's Play runs it.
const MakefileWebTarget = "web"

// MakefileManualTarget is the target of a platform with kinds the CI does
// not build (built: false): it builds them on the developer's machine and
// releases them with "warp-engine release upload" and "release publish".
const MakefileManualTarget = "mac-release"

// MakefileTargetModel is one target of a [MakefileReportModel].
type MakefileTargetModel struct {
	// Name is the target ("build").
	Name string `json:"name"`
	// Status is [MakefileTargetPresent], [MakefileTargetMissing] or
	// [MakefileTargetFailing].
	Status string `json:"status"`
	// Reason is why the target is required.
	Reason string `json:"reason"`
	// Detail is make's message when the target fails.
	Detail string `json:"detail,omitempty"`
}

// MakefileReportModel is what "check makefile" found in a project: its
// platform and the state of every target the contract requires of it.
type MakefileReportModel struct {
	// Directory is the project's directory.
	Directory string `json:"directory"`
	// Platform is the project's platform ("ebitengine").
	Platform string `json:"platform"`
	// App is the platform variant the kinds were taken from; empty for a
	// platform without variants.
	App string `json:"app,omitempty"`
	// PlatformKnown is false when the server's platforms could not be read:
	// only [MakefileContract] was checked then.
	PlatformKnown bool `json:"platformKnown"`
	// Note says what was not checked and why; empty when everything was.
	Note string `json:"note,omitempty"`
	// Targets are the required targets in order.
	Targets []MakefileTargetModel `json:"targets"`
}

// Broken returns the names of the targets that are missing or fail.
func (m MakefileReportModel) Broken() []string {
	broken := []string{}
	for _, target := range m.Targets {
		if target.Status != MakefileTargetPresent {
			broken = append(broken, target.Name)
		}
	}
	return broken
}

// ProjectMarkerModel is the platform a project's CI marker
// (".woodpecker.yaml": "platform: ebitengine", "app: gui") names.
type ProjectMarkerModel struct {
	// File is the marker file the platform was read from.
	File string
	// Platform is the platform; empty when no marker names one.
	Platform string
	// App is the platform variant; empty for the default one.
	App string
}
