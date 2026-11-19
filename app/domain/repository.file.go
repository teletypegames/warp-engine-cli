package domain

import (
	"net/http"
	"net/url"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// FileRepositoryInterface is the Admin API's files resource: the artifact
// store. Paths are relative to the store's root.
type FileRepositoryInterface interface {
	// All lists the directory dir ("" for the root).
	All(dir string) (*PageModel[FileEntryModel], error)
	// Upload sends the local file file_path into the directory dir.
	Upload(dir string, file_path string) (*FilePathModel, error)
	// Rename gives the file or directory path the name new_name, in place.
	Rename(path string, new_name string) (*FilePathModel, error)
	// Destroy deletes the file or empty directory path.
	Destroy(path string) error
	// CreateDirectory makes the directory name in the directory dir.
	CreateDirectory(dir string, name string) (*FilePathModel, error)
}

// FileRepository is the [FileRepositoryInterface] over the Admin API.
type FileRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /files?dir=…. The list has no paging.
func (r FileRepository) All(dir string) (*PageModel[FileEntryModel], error) {
	return list[FileEntryModel](r.Client, "/files", Filter{"dir": dir}, ListQuery{})
}

// Upload posts POST /files as multipart/form-data: the parts dir and file.
func (r FileRepository) Upload(dir string, file_path string) (*FilePathModel, error) {
	var envelope api.Envelope[FilePathModel]
	if err := r.Client.Upload("/files", map[string]string{"dir": dir}, "file", file_path, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Rename sends PATCH /files with path and newName.
func (r FileRepository) Rename(path string, new_name string) (*FilePathModel, error) {
	return api.Send[FilePathModel](r.Client, http.MethodPatch, "/files", Attributes{"path": path, "newName": new_name})
}

// Destroy sends DELETE /files?path=….
func (r FileRepository) Destroy(path string) error {
	return api.Delete(r.Client, "/files", url.Values{"path": {path}})
}

// CreateDirectory posts POST /files/directories with dir and name.
func (r FileRepository) CreateDirectory(dir string, name string) (*FilePathModel, error) {
	return api.Send[FilePathModel](r.Client, http.MethodPost, "/files/directories", Attributes{"dir": dir, "name": name})
}
