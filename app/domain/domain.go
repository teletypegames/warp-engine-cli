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
	// ExternalLinkService manages the external links.
	ExternalLinkService ExternalLinkServiceInterface
	// ForkLinkService moderates the forks other engines announced.
	ForkLinkService ForkLinkServiceInterface
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
	// PipelineService manages the CI pipelines.
	PipelineService PipelineServiceInterface
	// FileService manages the artifact store.
	FileService FileServiceInterface
	// RuntimeService manages the runtimes.
	RuntimeService RuntimeServiceInterface
	// LocalComponentService manages the components of WarpEngine Local.
	LocalComponentService LocalComponentServiceInterface
	// LocalOperationService follows the background jobs of WarpEngine Local.
	LocalOperationService LocalOperationServiceInterface
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
	// ExternalLinkRepository talks to the external links resource.
	ExternalLinkRepository ExternalLinkRepositoryInterface
	// ForkLinkRepository talks to the fork_links resource.
	ForkLinkRepository ForkLinkRepositoryInterface
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
	// PipelineRepository talks to the pipelines resource.
	PipelineRepository PipelineRepositoryInterface
	// FileRepository talks to the files resource.
	FileRepository FileRepositoryInterface
	// RuntimeRepository talks to the runtimes resource.
	RuntimeRepository RuntimeRepositoryInterface
	// LocalComponentRepository talks to the local_components resource.
	LocalComponentRepository LocalComponentRepositoryInterface
	// LocalOperationRepository talks to the local_operations resource.
	LocalOperationRepository LocalOperationRepositoryInterface
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
	external_link_repository := ExternalLinkRepository{Client: client}
	fork_link_repository := ForkLinkRepository{Client: client}
	software_image_repository := SoftwareImageRepository{Client: client}
	extra_repository := ExtraRepository{Client: client}
	software_status_repository := SoftwareStatusRepository{Client: client}
	software_category_repository := SoftwareCategoryRepository{Client: client}
	platform_link_repository := PlatformLinkRepository{Client: client}
	download_repository := DownloadRepository{Client: client}
	pipeline_repository := PipelineRepository{Client: client}
	file_repository := FileRepository{Client: client}
	runtime_repository := RuntimeRepository{Client: client}
	local_component_repository := LocalComponentRepository{Client: client}
	local_operation_repository := LocalOperationRepository{Client: client}
	sync_peer_repository := SyncPeerRepository{Client: client}
	sync_pairing_repository := SyncPairingRepository{Client: client}
	sync_event_repository := SyncEventRepository{Client: client}
	application_token_repository := ApplicationTokenRepository{Client: client}
	build_platform_repository := BuildPlatformRepository{Client: client}
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
		ExternalLinkService:     ExternalLinkService{ExternalLinkRepository: external_link_repository, SoftwareRepository: software_repository},
		ForkLinkService:         ForkLinkService{ForkLinkRepository: fork_link_repository, SoftwareRepository: software_repository},
		SoftwareImageService:    SoftwareImageService{SoftwareImageRepository: software_image_repository, SoftwareRepository: software_repository},
		ExtraService:            ExtraService{ExtraRepository: extra_repository, SoftwareRepository: software_repository},
		SoftwareStatusService:   SoftwareStatusService{SoftwareStatusRepository: software_status_repository},
		SoftwareCategoryService: SoftwareCategoryService{SoftwareCategoryRepository: software_category_repository},
		PlatformLinkService:     PlatformLinkService{PlatformLinkRepository: platform_link_repository, BuildPlatformService: build_platform_service},
		DownloadService:         DownloadService{DownloadRepository: download_repository, SoftwareRepository: software_repository},
		PipelineService:         PipelineService{PipelineRepository: pipeline_repository, SoftwareRepository: software_repository, BuildPlatformService: build_platform_service},
		FileService:             FileService{FileRepository: file_repository},
		RuntimeService:          RuntimeService{RuntimeRepository: runtime_repository},
		LocalComponentService:   LocalComponentService{LocalComponentRepository: local_component_repository, SoftwareRepository: software_repository},
		LocalOperationService:   LocalOperationService{LocalOperationRepository: local_operation_repository},
		SyncPeerService:         SyncPeerService{SyncPeerRepository: sync_peer_repository},
		SyncPairingService:      SyncPairingService{SyncPairingRepository: sync_pairing_repository},
		SyncEventService:        SyncEventService{SyncEventRepository: sync_event_repository, SoftwareRepository: software_repository},
		ApplicationTokenService: ApplicationTokenService{ApplicationTokenRepository: application_token_repository},
		BuildPlatformService:    build_platform_service,

		AuthRepository:             auth_repository,
		ProfileRepository:          profile_repository,
		LocalTokenRepository:       local_token_repository,
		SoftwareRepository:         software_repository,
		RawRepository:              raw_repository,
		ReleaseRepository:          release_repository,
		ReleaseAssetRepository:     release_asset_repository,
		ExternalLinkRepository:     external_link_repository,
		ForkLinkRepository:         fork_link_repository,
		SoftwareImageRepository:    software_image_repository,
		ExtraRepository:            extra_repository,
		SoftwareStatusRepository:   software_status_repository,
		SoftwareCategoryRepository: software_category_repository,
		PlatformLinkRepository:     platform_link_repository,
		DownloadRepository:         download_repository,
		PipelineRepository:         pipeline_repository,
		FileRepository:             file_repository,
		RuntimeRepository:          runtime_repository,
		LocalComponentRepository:   local_component_repository,
		LocalOperationRepository:   local_operation_repository,
		SyncPeerRepository:         sync_peer_repository,
		SyncPairingRepository:      sync_pairing_repository,
		SyncEventRepository:        sync_event_repository,
		ApplicationTokenRepository: application_token_repository,
		BuildPlatformRepository:    build_platform_repository,
	}
}
