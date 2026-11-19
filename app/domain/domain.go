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
	// SoftwareStatusService manages the software statuses.
	SoftwareStatusService SoftwareStatusServiceInterface
	// SoftwareCategoryService manages the software categories.
	SoftwareCategoryService SoftwareCategoryServiceInterface
	// PlatformLinkService manages the platform links.
	PlatformLinkService PlatformLinkServiceInterface
	// DownloadService reads the download log and its statistics.
	DownloadService DownloadServiceInterface

	// AuthRepository talks to the sign-in endpoints.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
	// SoftwareStatusRepository talks to the software statuses resource.
	SoftwareStatusRepository SoftwareStatusRepositoryInterface
	// SoftwareCategoryRepository talks to the software categories resource.
	SoftwareCategoryRepository SoftwareCategoryRepositoryInterface
	// PlatformLinkRepository talks to the platform links resource.
	PlatformLinkRepository PlatformLinkRepositoryInterface
	// DownloadRepository talks to the downloads resource.
	DownloadRepository DownloadRepositoryInterface
}

// NewDomain wires the repositories and services around client (shared by
// every repository that talks to the server), the configuration file at
// config_path and the WarpEngine Local data directory local_directory.
func NewDomain(client *api.Client, config_path string, local_directory string) Domain {
	auth_repository := AuthRepository{Client: client}
	profile_repository := ProfileRepository{Path: config_path}
	local_token_repository := LocalTokenRepository{Directory: local_directory}
	software_repository := SoftwareRepository{Client: client}
	software_status_repository := SoftwareStatusRepository{Client: client}
	software_category_repository := SoftwareCategoryRepository{Client: client}
	platform_link_repository := PlatformLinkRepository{Client: client}
	download_repository := DownloadRepository{Client: client}

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:          ProfileService{ProfileRepository: profile_repository},
		SoftwareService:         SoftwareService{SoftwareRepository: software_repository},
		SoftwareStatusService:   SoftwareStatusService{SoftwareStatusRepository: software_status_repository},
		SoftwareCategoryService: SoftwareCategoryService{SoftwareCategoryRepository: software_category_repository},
		PlatformLinkService:     PlatformLinkService{PlatformLinkRepository: platform_link_repository},
		DownloadService:         DownloadService{DownloadRepository: download_repository, SoftwareRepository: software_repository},

		AuthRepository:             auth_repository,
		ProfileRepository:          profile_repository,
		LocalTokenRepository:       local_token_repository,
		SoftwareRepository:         software_repository,
		SoftwareStatusRepository:   software_status_repository,
		SoftwareCategoryRepository: software_category_repository,
		PlatformLinkRepository:     platform_link_repository,
		DownloadRepository:         download_repository,
	}
}
