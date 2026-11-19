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
	// SyncPeerService manages the sync peers.
	SyncPeerService SyncPeerServiceInterface
	// SyncPairingService manages the caller's own sync pairings.
	SyncPairingService SyncPairingServiceInterface
	// SyncEventService reads the sync events and states.
	SyncEventService SyncEventServiceInterface
	// ApplicationTokenService manages the application tokens.
	ApplicationTokenService ApplicationTokenServiceInterface

	// AuthRepository talks to the sign-in endpoints.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
	// SyncPeerRepository talks to the sync_peers resource.
	SyncPeerRepository SyncPeerRepositoryInterface
	// SyncPairingRepository talks to the sync_pairings resource.
	SyncPairingRepository SyncPairingRepositoryInterface
	// SyncEventRepository talks to the sync_events resource.
	SyncEventRepository SyncEventRepositoryInterface
	// ApplicationTokenRepository talks to the application_tokens resource.
	ApplicationTokenRepository ApplicationTokenRepositoryInterface
}

// NewDomain wires the repositories and services around client (shared by
// every repository that talks to the server), the configuration file at
// config_path and the WarpEngine Local data directory local_directory.
func NewDomain(client *api.Client, config_path string, local_directory string) Domain {
	auth_repository := AuthRepository{Client: client}
	profile_repository := ProfileRepository{Path: config_path}
	local_token_repository := LocalTokenRepository{Directory: local_directory}
	software_repository := SoftwareRepository{Client: client}
	sync_peer_repository := SyncPeerRepository{Client: client}
	sync_pairing_repository := SyncPairingRepository{Client: client}
	sync_event_repository := SyncEventRepository{Client: client}
	application_token_repository := ApplicationTokenRepository{Client: client}

	return Domain{
		AuthService: AuthService{
			AuthRepository:       auth_repository,
			ProfileRepository:    profile_repository,
			LocalTokenRepository: local_token_repository,
		},
		ProfileService:     ProfileService{ProfileRepository: profile_repository},
		SoftwareService:    SoftwareService{SoftwareRepository: software_repository},
		SyncPeerService:    SyncPeerService{SyncPeerRepository: sync_peer_repository},
		SyncPairingService: SyncPairingService{SyncPairingRepository: sync_pairing_repository},
		SyncEventService: SyncEventService{
			SyncEventRepository: sync_event_repository,
			SoftwareRepository:  software_repository,
		},
		ApplicationTokenService: ApplicationTokenService{ApplicationTokenRepository: application_token_repository},

		AuthRepository:             auth_repository,
		ProfileRepository:          profile_repository,
		LocalTokenRepository:       local_token_repository,
		SoftwareRepository:         software_repository,
		SyncPeerRepository:         sync_peer_repository,
		SyncPairingRepository:      sync_pairing_repository,
		SyncEventRepository:        sync_event_repository,
		ApplicationTokenRepository: application_token_repository,
	}
}
