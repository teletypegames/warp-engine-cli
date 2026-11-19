package domain

import "git.teletypegames.org/engines/warp-engine-cli/app/api"

// Domain holds every service and repository of the program, wired together.
// It is built once by [NewDomain]; the command line takes its controllers'
// services from it.
type Domain struct {
	// AuthService signs in and connects.
	AuthService AuthServiceInterface
	// ProfileService keeps the profiles.
	ProfileService ProfileServiceInterface
	// SoftwareService manages the softwares.
	SoftwareService SoftwareServiceInterface
	// ReleaseService manages the releases.
	ReleaseService ReleaseServiceInterface
	// ReleaseAssetService manages the files of the releases.
	ReleaseAssetService ReleaseAssetServiceInterface
	// ExternalLinkService manages the external links.
	ExternalLinkService ExternalLinkServiceInterface
	// SoftwareImageService manages the images of the softwares.
	SoftwareImageService SoftwareImageServiceInterface
	// ExtraService manages the host's additional resources of the softwares.
	ExtraService ExtraServiceInterface

	// AuthRepository talks to the sign-in endpoints.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
	// ReleaseRepository talks to the releases resource.
	ReleaseRepository ReleaseRepositoryInterface
	// ReleaseAssetRepository talks to the release assets resource.
	ReleaseAssetRepository ReleaseAssetRepositoryInterface
	// ExternalLinkRepository talks to the external links resource.
	ExternalLinkRepository ExternalLinkRepositoryInterface
	// SoftwareImageRepository talks to the software images resource.
	SoftwareImageRepository SoftwareImageRepositoryInterface
	// ExtraRepository talks to the extras resource.
	ExtraRepository ExtraRepositoryInterface
}

// NewDomain wires the repositories and services around client (shared by
// every repository that talks to the server), the configuration file at
// config_path and the WarpEngine Local data directory local_directory.
func NewDomain(client *api.Client, config_path string, local_directory string) Domain {
	auth_repository := AuthRepository{Client: client}
	profile_repository := ProfileRepository{Path: config_path}
	local_token_repository := LocalTokenRepository{Directory: local_directory}
	software_repository := SoftwareRepository{Client: client}
	release_repository := ReleaseRepository{Client: client}
	release_asset_repository := ReleaseAssetRepository{Client: client}
	external_link_repository := ExternalLinkRepository{Client: client}
	software_image_repository := SoftwareImageRepository{Client: client}
	extra_repository := ExtraRepository{Client: client}

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:       ProfileService{ProfileRepository: profile_repository},
		SoftwareService:      SoftwareService{SoftwareRepository: software_repository},
		ReleaseService:       ReleaseService{ReleaseRepository: release_repository, SoftwareRepository: software_repository},
		ReleaseAssetService:  ReleaseAssetService{ReleaseAssetRepository: release_asset_repository},
		ExternalLinkService:  ExternalLinkService{ExternalLinkRepository: external_link_repository, SoftwareRepository: software_repository},
		SoftwareImageService: SoftwareImageService{SoftwareImageRepository: software_image_repository, SoftwareRepository: software_repository},
		ExtraService:         ExtraService{ExtraRepository: extra_repository, SoftwareRepository: software_repository},

		AuthRepository:          auth_repository,
		ProfileRepository:       profile_repository,
		LocalTokenRepository:    local_token_repository,
		SoftwareRepository:      software_repository,
		ReleaseRepository:       release_repository,
		ReleaseAssetRepository:  release_asset_repository,
		ExternalLinkRepository:  external_link_repository,
		SoftwareImageRepository: software_image_repository,
		ExtraRepository:         extra_repository,
	}
}
