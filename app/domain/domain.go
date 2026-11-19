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

	// AuthRepository talks to the sign-in endpoints.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
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
}

// NewDomain wires the repositories and services around client (shared by
// every repository that talks to the server), the configuration file at
// config_path and the WarpEngine Local data directory local_directory.
func NewDomain(client *api.Client, config_path string, local_directory string) Domain {
	auth_repository := AuthRepository{Client: client}
	profile_repository := ProfileRepository{Path: config_path}
	local_token_repository := LocalTokenRepository{Directory: local_directory}
	software_repository := SoftwareRepository{Client: client}
	pipeline_repository := PipelineRepository{Client: client}
	file_repository := FileRepository{Client: client}
	runtime_repository := RuntimeRepository{Client: client}
	local_component_repository := LocalComponentRepository{Client: client}
	local_operation_repository := LocalOperationRepository{Client: client}

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:  ProfileService{ProfileRepository: profile_repository},
		SoftwareService: SoftwareService{SoftwareRepository: software_repository},
		PipelineService: PipelineService{PipelineRepository: pipeline_repository, SoftwareRepository: software_repository},
		FileService:     FileService{FileRepository: file_repository},
		RuntimeService:  RuntimeService{RuntimeRepository: runtime_repository},
		LocalComponentService: LocalComponentService{
			LocalComponentRepository: local_component_repository,
			SoftwareRepository:       software_repository,
		},
		LocalOperationService: LocalOperationService{LocalOperationRepository: local_operation_repository},

		AuthRepository:           auth_repository,
		ProfileRepository:        profile_repository,
		LocalTokenRepository:     local_token_repository,
		SoftwareRepository:       software_repository,
		PipelineRepository:       pipeline_repository,
		FileRepository:           file_repository,
		RuntimeRepository:        runtime_repository,
		LocalComponentRepository: local_component_repository,
		LocalOperationRepository: local_operation_repository,
	}
}
