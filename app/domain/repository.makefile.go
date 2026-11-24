package domain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// MakefileRepositoryInterface asks make about the targets of a project's
// Makefile without building anything.
type MakefileRepositoryInterface interface {
	// Exists reports whether directory has a Makefile make reads.
	Exists(directory string) bool
	// Probe runs make -n target in directory and returns
	// [MakefileTargetPresent], [MakefileTargetMissing] or
	// [MakefileTargetFailing] with make's message; an error only when make
	// cannot be started.
	Probe(directory string, target string) (string, string, error)
}

// MakefileRepository is the [MakefileRepositoryInterface] over the make
// program.
type MakefileRepository struct {
	// Make is the program to run; empty runs "make".
	Make string
}

// MakefileNames are the file names GNU make reads, in its order.
var MakefileNames = []string{"GNUmakefile", "makefile", "Makefile"}

// Exists looks for one of [MakefileNames] in directory.
func (r MakefileRepository) Exists(directory string) bool {
	for _, name := range MakefileNames {
		if info, err := os.Stat(filepath.Join(directory, name)); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

// Probe runs "make -n target" in directory. make -n prints the recipes
// instead of running them (lines that call $(MAKE) run, also with -n); a
// "No rule to make target" naming target itself means the target is
// missing, any other failure that it fails.
func (r MakefileRepository) Probe(directory string, target string) (string, string, error) {
	program := r.Make
	if program == "" {
		program = "make"
	}
	command := exec.Command(program, "-n", target)
	command.Dir = directory
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	err := command.Run()
	if err == nil {
		return MakefileTargetPresent, "", nil
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		return "", "", err
	}
	text := strings.TrimSpace(output.String())
	missing := regexp.MustCompile("No rule to make target [`'\"]" + regexp.QuoteMeta(target) + "['\"]")
	if missing.MatchString(text) {
		return MakefileTargetMissing, "", nil
	}
	return MakefileTargetFailing, lastLine(text), nil
}

// lastLine returns the last non-empty line of text.
func lastLine(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

// ProjectRepositoryInterface reads what a project's files say about it.
type ProjectRepositoryInterface interface {
	// Marker returns the platform and app variant the project's CI marker
	// names; an empty Platform when none does.
	Marker(directory string) (*ProjectMarkerModel, error)
}

// ProjectRepository is the [ProjectRepositoryInterface] over the files of
// the project.
type ProjectRepository struct{}

// MarkerFiles are the files [ProjectRepository.Marker] reads, in order:
// the Woodpecker marker, a warp.yaml, then metadata.json's "platform".
var MarkerFiles = []string{".woodpecker.yaml", ".woodpecker.yml", "warp.yaml", "metadata.json"}

// markerLine matches a "platform:" or "app:" line of a YAML marker.
var markerLine = regexp.MustCompile(`^(platform|app):\s*["']?([A-Za-z0-9_-]+)["']?\s*(#.*)?$`)

// Marker reads the first of [MarkerFiles] that names a platform.
func (r ProjectRepository) Marker(directory string) (*ProjectMarkerModel, error) {
	for _, name := range MarkerFiles {
		content, err := os.ReadFile(filepath.Join(directory, name))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		marker := ProjectMarkerModel{File: name}
		if strings.HasSuffix(name, ".json") {
			var metadata struct {
				Platform string `json:"platform"`
				App      string `json:"app"`
			}
			if json.Unmarshal(content, &metadata) == nil {
				marker.Platform, marker.App = metadata.Platform, metadata.App
			}
		} else {
			scanner := bufio.NewScanner(bytes.NewReader(content))
			for scanner.Scan() {
				match := markerLine.FindStringSubmatch(strings.TrimRight(scanner.Text(), "\r"))
				if match == nil {
					continue
				}
				if match[1] == "platform" {
					marker.Platform = match[2]
				} else {
					marker.App = match[2]
				}
			}
		}
		if marker.Platform != "" {
			return &marker, nil
		}
	}
	return &ProjectMarkerModel{}, nil
}
