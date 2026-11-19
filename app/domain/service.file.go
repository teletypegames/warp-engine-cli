package domain

import (
	"errors"
	"path"
	"strings"
)

// FileServiceInterface manages the artifact store: the directory tree of
// release files, runtimes and anything else the engine serves. Paths are
// relative to the store's root and use "/".
type FileServiceInterface interface {
	// All lists the directory dir ("" for the root).
	All(dir string) (*PageModel[FileEntryModel], error)
	// Upload sends the local file file_path into the directory dir.
	Upload(file_path string, dir string) (*FilePathModel, error)
	// Rename gives the file or directory file_path the name new_name.
	Rename(file_path string, new_name string) (*FilePathModel, error)
	// Destroy deletes the file or empty directory file_path.
	Destroy(file_path string) error
	// CreateDirectory makes the directory dir_path; its parent must exist.
	CreateDirectory(dir_path string) (*FilePathModel, error)
}

// FileService is the [FileServiceInterface] over a [FileRepositoryInterface].
type FileService struct {
	// FileRepository talks to the files resource.
	FileRepository FileRepositoryInterface
}

// All lists the directory; leading and trailing slashes of dir are dropped.
func (s FileService) All(dir string) (*PageModel[FileEntryModel], error) {
	return s.FileRepository.All(cleanStorePath(dir))
}

// Upload uploads file_path into dir (the root when empty).
func (s FileService) Upload(file_path string, dir string) (*FilePathModel, error) {
	if file_path == "" {
		return nil, errors.New("no file given")
	}
	return s.FileRepository.Upload(cleanStorePath(dir), file_path)
}

// Rename renames file_path in place: new_name is a name, not a path, so a
// "/" in it is an error (the store cannot move a file to another directory).
func (s FileService) Rename(file_path string, new_name string) (*FilePathModel, error) {
	file_path = cleanStorePath(file_path)
	if file_path == "" {
		return nil, errors.New("no path given")
	}
	if new_name == "" || strings.Contains(new_name, "/") {
		return nil, errors.New("the new name must be a name without \"/\": a file is renamed in its directory")
	}
	return s.FileRepository.Rename(file_path, new_name)
}

// Destroy deletes file_path; the store's root cannot be deleted.
func (s FileService) Destroy(file_path string) error {
	file_path = cleanStorePath(file_path)
	if file_path == "" {
		return errors.New("no path given")
	}
	return s.FileRepository.Destroy(file_path)
}

// CreateDirectory splits dir_path into its parent and its name and makes it.
func (s FileService) CreateDirectory(dir_path string) (*FilePathModel, error) {
	dir_path = cleanStorePath(dir_path)
	if dir_path == "" {
		return nil, errors.New("no directory given")
	}
	parent := path.Dir(dir_path)
	if parent == "." {
		parent = ""
	}
	return s.FileRepository.CreateDirectory(parent, path.Base(dir_path))
}

// cleanStorePath trims the slashes around a store path ("/games/" →
// "games"); the root is "".
func cleanStorePath(store_path string) string {
	return strings.Trim(store_path, "/")
}
