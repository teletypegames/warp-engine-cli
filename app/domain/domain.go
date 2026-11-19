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

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:  ProfileService{ProfileRepository: profile_repository},
		SoftwareService: SoftwareService{SoftwareRepository: software_repository},
		RawService:      RawService{RawRepository: raw_repository},

		AuthRepository:       auth_repository,
		ProfileRepository:    profile_repository,
		LocalTokenRepository: local_token_repository,
		SoftwareRepository:   software_repository,
		RawRepository:        raw_repository,
	}
}
