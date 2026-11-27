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
	// RawService sends any request.
	RawService RawServiceInterface
	// ReleaseService manages the releases.
	ReleaseService ReleaseServiceInterface
	// ReleaseAssetService manages the files of the releases.
	ReleaseAssetService ReleaseAssetServiceInterface
	// SoftwareLinkService manages the links of the softwares.
	SoftwareLinkService SoftwareLinkServiceInterface
	// SoftwareSourceService manages the source repositories of the softwares.
	SoftwareSourceService SoftwareSourceServiceInterface
	// SoftwareForkService moderates the forks other engines announced.
	SoftwareForkService SoftwareForkServiceInterface
	// SoftwareImageService manages the images of the softwares.
	SoftwareImageService SoftwareImageServiceInterface
	// ExtraService manages the host's additional resources of the softwares.
	ExtraService ExtraServiceInterface
	// SoftwareStatusService manages the software statuses.
	SoftwareStatusService SoftwareStatusServiceInterface
	// SoftwareCategoryService manages the software categories.
	SoftwareCategoryService SoftwareCategoryServiceInterface
	// PlatformLinkService manages the platform links.
	PlatformLinkService PlatformLinkServiceInterface
	// DownloadService reads the download log and its statistics.
	DownloadService DownloadServiceInterface
	// FeedService reads the feeds.
	FeedService FeedServiceInterface
	// CiRepositoryService manages the CI pipelines.
	CiRepositoryService CiRepositoryServiceInterface
	// FileService manages the artifact store.
	FileService FileServiceInterface
	// RuntimeService manages the runtimes.
	RuntimeService RuntimeServiceInterface
	// LocalProviderService manages the components of WarpEngine Local.
	LocalProviderService LocalProviderServiceInterface
	// LocalTaskService follows the background jobs of WarpEngine Local.
	LocalTaskService LocalTaskServiceInterface
	// SyncPeerService manages the sync peers.
	SyncPeerService SyncPeerServiceInterface
	// SyncPairingService manages the caller's own sync pairings.
	SyncPairingService SyncPairingServiceInterface
	// SyncEventService reads the sync events and states.
	SyncEventService SyncEventServiceInterface
	// ApplicationTokenService manages the application tokens.
	ApplicationTokenService ApplicationTokenServiceInterface
	// BuildPlatformService knows the server's platforms, asset kinds and app
	// variants.
	BuildPlatformService BuildPlatformServiceInterface
	// BuildService releases files through the build endpoints.
	BuildService BuildServiceInterface
	// MakefileService checks a project's Makefile against the contract.
	MakefileService MakefileServiceInterface

	// AuthRepository talks to the sign-in endpoints.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
	// RawRepository sends any request.
	RawRepository RawRepositoryInterface
	// ReleaseRepository talks to the releases resource.
	ReleaseRepository ReleaseRepositoryInterface
	// ReleaseAssetRepository talks to the release assets resource.
	ReleaseAssetRepository ReleaseAssetRepositoryInterface
	// SoftwareLinkRepository talks to the links resource.
	SoftwareLinkRepository SoftwareLinkRepositoryInterface
	// SoftwareSourceRepository talks to the source of a software.
	SoftwareSourceRepository SoftwareSourceRepositoryInterface
	// SoftwareForkRepository talks to the forks resource.
	SoftwareForkRepository SoftwareForkRepositoryInterface
	// SoftwareImageRepository talks to the software images resource.
	SoftwareImageRepository SoftwareImageRepositoryInterface
	// ExtraRepository talks to the extras resource.
	ExtraRepository ExtraRepositoryInterface
	// SoftwareStatusRepository talks to the software statuses resource.
	SoftwareStatusRepository SoftwareStatusRepositoryInterface
	// SoftwareCategoryRepository talks to the software categories resource.
	SoftwareCategoryRepository SoftwareCategoryRepositoryInterface
	// PlatformLinkRepository talks to the platform links resource.
	PlatformLinkRepository PlatformLinkRepositoryInterface
	// DownloadRepository talks to the downloads resource.
	DownloadRepository DownloadRepositoryInterface
	// FeedRepository talks to the feeds resource.
	FeedRepository FeedRepositoryInterface
	// CiRepositoryRepository talks to the ci_repositories resource.
	CiRepositoryRepository CiRepositoryRepositoryInterface
	// FileRepository talks to the files resource.
	FileRepository FileRepositoryInterface
	// RuntimeRepository talks to the runtimes resource.
	RuntimeRepository RuntimeRepositoryInterface
	// LocalProviderRepository talks to the local_providers resource.
	LocalProviderRepository LocalProviderRepositoryInterface
	// LocalTaskRepository talks to the local_tasks resource.
	LocalTaskRepository LocalTaskRepositoryInterface
	// SyncPeerRepository talks to the sync_peers resource.
	SyncPeerRepository SyncPeerRepositoryInterface
	// SyncPairingRepository talks to the sync_pairings resource.
	SyncPairingRepository SyncPairingRepositoryInterface
	// SyncEventRepository talks to the sync_events resource.
	SyncEventRepository SyncEventRepositoryInterface
	// ApplicationTokenRepository talks to the application_tokens resource.
	ApplicationTokenRepository ApplicationTokenRepositoryInterface
	// BuildPlatformRepository reads the server's build description.
	BuildPlatformRepository BuildPlatformRepositoryInterface
	// BuildRepository talks to the build endpoints.
	BuildRepository BuildRepositoryInterface
	// MakefileRepository asks make about a Makefile.
	MakefileRepository MakefileRepositoryInterface
	// ProjectRepository reads a project's CI marker.
	ProjectRepository ProjectRepositoryInterface
}

// NewDomain wires the repositories and services around client (shared by
// every repository that talks to the server), the configuration file at
// config_path and the WarpEngine Local data directory local_directory.
func NewDomain(client *api.Client, config_path string, local_directory string) Domain {
	auth_repository := AuthRepository{Client: client}
	profile_repository := ProfileRepository{Path: config_path}
	local_token_repository := LocalTokenRepository{Directory: local_directory}
	software_repository := SoftwareRepository{Client: client}
	raw_repository := RawRepository{Client: client}
	release_repository := ReleaseRepository{Client: client}
	release_asset_repository := ReleaseAssetRepository{Client: client}
	software_link_repository := SoftwareLinkRepository{Client: client}
	software_source_repository := SoftwareSourceRepository{Client: client}
	software_fork_repository := SoftwareForkRepository{Client: client}
	software_image_repository := SoftwareImageRepository{Client: client}
	extra_repository := ExtraRepository{Client: client}
	software_status_repository := SoftwareStatusRepository{Client: client}
	software_category_repository := SoftwareCategoryRepository{Client: client}
	platform_link_repository := PlatformLinkRepository{Client: client}
	download_repository := DownloadRepository{Client: client}
	feed_repository := FeedRepository{Client: client}
	ci_repository_repository := CiRepositoryRepository{Client: client}
	file_repository := FileRepository{Client: client}
	runtime_repository := RuntimeRepository{Client: client}
	local_provider_repository := LocalProviderRepository{Client: client}
	local_task_repository := LocalTaskRepository{Client: client}
	sync_peer_repository := SyncPeerRepository{Client: client}
	sync_pairing_repository := SyncPairingRepository{Client: client}
	sync_event_repository := SyncEventRepository{Client: client}
	application_token_repository := ApplicationTokenRepository{Client: client}
	build_platform_repository := BuildPlatformRepository{Client: client}
	build_repository := BuildRepository{Client: client}
	makefile_repository := MakefileRepository{}
	project_repository := ProjectRepository{}
	build_platform_service := BuildPlatformService{BuildPlatformRepository: build_platform_repository, Cache: &BuildPlatformCache{}}

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:          ProfileService{ProfileRepository: profile_repository},
		SoftwareService:         SoftwareService{SoftwareRepository: software_repository, BuildPlatformService: build_platform_service},
		RawService:              RawService{RawRepository: raw_repository},
		ReleaseService:          ReleaseService{ReleaseRepository: release_repository, SoftwareRepository: software_repository},
		ReleaseAssetService:     ReleaseAssetService{ReleaseAssetRepository: release_asset_repository, BuildPlatformService: build_platform_service},
		SoftwareLinkService:     SoftwareLinkService{SoftwareLinkRepository: software_link_repository, SoftwareRepository: software_repository},
		SoftwareSourceService:   SoftwareSourceService{SoftwareSourceRepository: software_source_repository, SoftwareRepository: software_repository},
		SoftwareForkService:     SoftwareForkService{SoftwareForkRepository: software_fork_repository, SoftwareRepository: software_repository},
		SoftwareImageService:    SoftwareImageService{SoftwareImageRepository: software_image_repository, SoftwareRepository: software_repository},
		ExtraService:            ExtraService{ExtraRepository: extra_repository, SoftwareRepository: software_repository},
		SoftwareStatusService:   SoftwareStatusService{SoftwareStatusRepository: software_status_repository},
		SoftwareCategoryService: SoftwareCategoryService{SoftwareCategoryRepository: software_category_repository},
		PlatformLinkService:     PlatformLinkService{PlatformLinkRepository: platform_link_repository, BuildPlatformService: build_platform_service},
		DownloadService:         DownloadService{DownloadRepository: download_repository, SoftwareRepository: software_repository},
		FeedService:             FeedService{FeedRepository: feed_repository},
		CiRepositoryService:     CiRepositoryService{CiRepositoryRepository: ci_repository_repository, SoftwareRepository: software_repository, BuildPlatformService: build_platform_service},
		FileService:             FileService{FileRepository: file_repository},
		RuntimeService:          RuntimeService{RuntimeRepository: runtime_repository},
		LocalProviderService:    LocalProviderService{LocalProviderRepository: local_provider_repository, SoftwareRepository: software_repository},
		LocalTaskService:        LocalTaskService{LocalTaskRepository: local_task_repository},
		SyncPeerService:         SyncPeerService{SyncPeerRepository: sync_peer_repository},
		SyncPairingService:      SyncPairingService{SyncPairingRepository: sync_pairing_repository},
		SyncEventService:        SyncEventService{SyncEventRepository: sync_event_repository, SoftwareRepository: software_repository},
		ApplicationTokenService: ApplicationTokenService{ApplicationTokenRepository: application_token_repository},
		BuildPlatformService:    build_platform_service,
		BuildService:            BuildService{BuildRepository: build_repository, BuildPlatformService: build_platform_service},
		MakefileService: MakefileService{
			MakefileRepository:      makefile_repository,
			ProjectRepository:       project_repository,
			BuildPlatformRepository: build_platform_repository,
		},

		AuthRepository:             auth_repository,
		ProfileRepository:          profile_repository,
		LocalTokenRepository:       local_token_repository,
		SoftwareRepository:         software_repository,
		RawRepository:              raw_repository,
		ReleaseRepository:          release_repository,
		ReleaseAssetRepository:     release_asset_repository,
		SoftwareLinkRepository:     software_link_repository,
		SoftwareSourceRepository:   software_source_repository,
		SoftwareForkRepository:     software_fork_repository,
		SoftwareImageRepository:    software_image_repository,
		ExtraRepository:            extra_repository,
		SoftwareStatusRepository:   software_status_repository,
		SoftwareCategoryRepository: software_category_repository,
		PlatformLinkRepository:     platform_link_repository,
		DownloadRepository:         download_repository,
		FeedRepository:             feed_repository,
		CiRepositoryRepository:     ci_repository_repository,
		FileRepository:             file_repository,
		RuntimeRepository:          runtime_repository,
		LocalProviderRepository:    local_provider_repository,
		LocalTaskRepository:        local_task_repository,
		SyncPeerRepository:         sync_peer_repository,
		SyncPairingRepository:      sync_pairing_repository,
		SyncEventRepository:        sync_event_repository,
		ApplicationTokenRepository: application_token_repository,
		BuildPlatformRepository:    build_platform_repository,
		BuildRepository:            build_repository,
		MakefileRepository:         makefile_repository,
		ProjectRepository:          project_repository,
	}
}
