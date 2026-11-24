package domain

import (
	"net/url"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// BuildRepositoryInterface is the build endpoints of the server, the ones a
// CI pipeline releases through. They take an application token with the
// upload (upload) or update (publish) scope instead of an Admin API token.
type BuildRepositoryInterface interface {
	// Upload sends the file at file_path into the folder of target's
	// release; sha256, when not empty, is the digest the server checks.
	Upload(target BuildTargetModel, file_path string, sha256 string) (*BuildUploadModel, error)
	// Publish registers the uploaded files of target as a release.
	Publish(target BuildTargetModel) (*BuildPublishModel, error)
}

// BuildRepository is the [BuildRepositoryInterface] over POST /build/upload
// and POST /build/publish (outside the Admin API), signed with the client's
// token in [api.UpdateHeader].
type BuildRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// Upload posts POST /build/upload?name=&version=[&sha256=] with the file.
func (r BuildRepository) Upload(target BuildTargetModel, file_path string, sha256 string) (*BuildUploadModel, error) {
	query := url.Values{"name": {target.Name}, "version": {target.Version}}
	if sha256 != "" {
		query.Set("sha256", sha256)
	}
	var answer BuildUploadModel
	if err := r.Client.BuildUpload("/build/upload", query, file_path, &answer); err != nil {
		return nil, err
	}
	return &answer, nil
}

// Publish posts POST /build/publish?name=&platform=&version=[&commit=][&app=].
func (r BuildRepository) Publish(target BuildTargetModel) (*BuildPublishModel, error) {
	query := url.Values{"name": {target.Name}, "platform": {target.Platform}, "version": {target.Version}}
	if target.Commit != "" {
		query.Set("commit", target.Commit)
	}
	if target.App != "" {
		query.Set("app", target.App)
	}
	var answer BuildPublishModel
	if err := r.Client.BuildPost("/build/publish", query, &answer); err != nil {
		return nil, err
	}
	return &answer, nil
}
