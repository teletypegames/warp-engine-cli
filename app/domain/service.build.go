package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// BuildServiceInterface releases files built outside the CI (the macOS
// packages, which only a Mac signs) the way the CI does: upload them into
// the release's folder, then publish the release again.
type BuildServiceInterface interface {
	// Upload sends every file of file_paths into target's release, in
	// order, and returns what the server stored; it stops at the first
	// failure.
	Upload(target BuildTargetModel, file_paths []string) ([]BuildUploadModel, error)
	// Publish registers target's release with every file in its folder.
	Publish(target BuildTargetModel) (*BuildPublishModel, error)
}

// BuildService is the [BuildServiceInterface] over a
// [BuildRepositoryInterface].
type BuildService struct {
	// BuildRepository talks to the build endpoints.
	BuildRepository BuildRepositoryInterface
	// BuildPlatformService checks the platform and the app variant; nil
	// leaves them to the server.
	BuildPlatformService BuildPlatformServiceInterface
}

// ErrNoBuildFiles is returned by [BuildService.Upload] without files.
var ErrNoBuildFiles = errors.New("no file given")

// Upload checks that target has a name and a version and that every file
// exists and is named <name>-<version>… as the server requires, then
// uploads them one by one with their SHA-256 digest, which the server
// verifies.
func (s BuildService) Upload(target BuildTargetModel, file_paths []string) ([]BuildUploadModel, error) {
	if err := requireTarget(target, false); err != nil {
		return nil, err
	}
	if len(file_paths) == 0 {
		return nil, ErrNoBuildFiles
	}
	prefix := target.Name + "-" + target.Version
	for _, file_path := range file_paths {
		if !strings.HasPrefix(filepath.Base(file_path), prefix) {
			return nil, fmt.Errorf("%s: the file name must start with %s", file_path, prefix)
		}
		info, err := os.Stat(file_path)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			return nil, fmt.Errorf("%s is a directory", file_path)
		}
	}
	uploads := make([]BuildUploadModel, 0, len(file_paths))
	for _, file_path := range file_paths {
		digest, err := fileDigest(file_path)
		if err != nil {
			return uploads, err
		}
		upload, err := s.BuildRepository.Upload(target, file_path, digest)
		if err != nil {
			return uploads, fmt.Errorf("upload %s: %w", file_path, err)
		}
		uploads = append(uploads, *upload)
	}
	return uploads, nil
}

// Publish checks that target has a name, a version and a platform (and,
// when the server's platforms can be read, that the platform and the app
// variant exist), then publishes it.
func (s BuildService) Publish(target BuildTargetModel) (*BuildPublishModel, error) {
	if err := requireTarget(target, true); err != nil {
		return nil, err
	}
	if s.BuildPlatformService != nil {
		if err := s.BuildPlatformService.CheckPlatform(target.Platform); err != nil {
			return nil, err
		}
		if target.App != "" {
			if err := s.BuildPlatformService.CheckApp(target.Platform, target.App); err != nil {
				return nil, err
			}
		}
	}
	return s.BuildRepository.Publish(target)
}

// requireTarget returns an error when target lacks its name, its version
// or, with platform, its platform.
func requireTarget(target BuildTargetModel, platform bool) error {
	switch {
	case target.Name == "":
		return errors.New("no software name given (--name)")
	case target.Version == "":
		return errors.New("no version given (--version)")
	case platform && target.Platform == "":
		return errors.New("no platform given (--platform)")
	}
	return nil
}

// fileDigest returns the hex SHA-256 digest of the file at file_path.
func fileDigest(file_path string) (string, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
