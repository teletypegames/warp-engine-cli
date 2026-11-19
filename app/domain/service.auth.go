package domain

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"
)

// ErrDenied is returned by [AuthService.WaitForApproval] when the person
// refused the sign-in in the browser.
var ErrDenied = errors.New("the sign-in was denied")

// ErrExpired is returned by [AuthService.WaitForApproval] when the codes
// expired before anyone answered.
var ErrExpired = errors.New("the sign-in code expired; start again")

// minimumPollInterval is the shortest wait between two polls of a device
// sign-in, whatever the server asks for.
const minimumPollInterval = 2 * time.Second

// AuthServiceInterface signs in to servers and decides which server and
// token a command talks to.
type AuthServiceInterface interface {
	// Connect picks the server and token for selection and points the
	// client at them; it returns the profile it used (an unsaved one when
	// the server came from --server).
	Connect(selection Selection) (*ProfileModel, error)
	// Discover points the client at server_url, without a token, and reads
	// what the engine says about itself.
	Discover(server_url string) (*ServerModel, error)
	// StartDeviceLogin starts a device sign-in at server_url.
	StartDeviceLogin(server_url string, client_name string, scopes []string) (*DeviceGrantModel, error)
	// WaitForApproval polls grant until it is approved, denied or expired,
	// and returns the token.
	WaitForApproval(grant DeviceGrantModel) (string, error)
	// FinishLogin checks token at server_url and saves it as the current
	// profile name (a name made from the URL when empty).
	FinishLogin(server_url string, token string, name string) (*ProfileModel, error)
	// LoginLocal saves a profile for WarpEngine Local at server_url that
	// reads its token from the token file.
	LoginLocal(server_url string, name string) (*ProfileModel, error)
	// Logout removes the profile name (the current one when empty) and,
	// with revoke, revokes its token on the server first.
	Logout(name string, revoke bool) (*ProfileModel, error)
	// Whoami returns who the server takes the current token for.
	Whoami() (*PrincipalModel, error)
}

// AuthService is the [AuthServiceInterface] over the auth, profile and
// local-token repositories.
type AuthService struct {
	// AuthRepository talks to the server.
	AuthRepository AuthRepositoryInterface
	// ProfileRepository keeps the profiles.
	ProfileRepository ProfileRepositoryInterface
	// LocalTokenRepository reads WarpEngine Local's token.
	LocalTokenRepository LocalTokenRepositoryInterface
	// Sleep waits between two polls; nil uses time.Sleep (tests replace it).
	Sleep func(time.Duration)
}

// Connect resolves selection, first match wins:
//
//  1. Server set: that URL, with Token if set, else the token of a saved
//     profile with the same URL, else (on a loopback address) the WarpEngine
//     Local token file, else no token.
//  2. Otherwise the profile Profile names, or the current profile:
//     its URL and token (read from the token file for a Local profile);
//     Token, when set, replaces the token.
//  3. Neither: [ErrNoProfile].
//
// The client is pointed at the result before Connect returns.
func (s AuthService) Connect(selection Selection) (*ProfileModel, error) {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return nil, err
	}
	var profile ProfileModel
	if selection.Server != "" {
		profile = ProfileModel{Name: ProfileName(selection.Server), URL: strings.TrimRight(selection.Server, "/"), Token: selection.Token}
		if profile.Token == "" {
			for _, saved := range config.Profiles {
				if strings.TrimRight(saved.URL, "/") == profile.URL {
					profile = saved
					break
				}
			}
		}
		if profile.Token == "" && (profile.TokenFile || IsLoopback(profile.URL)) {
			profile.Token, _ = s.LocalTokenRepository.Read()
		}
	} else {
		name := selection.Profile
		if name == "" {
			name = config.Current
		}
		if name == "" {
			return nil, ErrNoProfile
		}
		saved, ok := config.Profiles[name]
		if !ok {
			return nil, errors.New("no profile named \"" + name + "\" (see \"warp-engine profile list\")")
		}
		profile = saved
		if profile.TokenFile {
			if profile.Token, err = s.LocalTokenRepository.Read(); err != nil {
				return nil, err
			}
		}
		if selection.Token != "" {
			profile.Token = selection.Token
		}
	}
	s.AuthRepository.Use(profile.URL, profile.Token)
	return &profile, nil
}

// Discover asks server_url about itself without signing in.
func (s AuthService) Discover(server_url string) (*ServerModel, error) {
	s.AuthRepository.Use(server_url, "")
	return s.AuthRepository.Discover()
}

// StartDeviceLogin points the client at server_url without a token and
// starts a device sign-in for client_name; nil or empty scopes ask for
// [DefaultScopes].
func (s AuthService) StartDeviceLogin(server_url string, client_name string, scopes []string) (*DeviceGrantModel, error) {
	if len(scopes) == 0 {
		scopes = DefaultScopes
	}
	s.AuthRepository.Use(server_url, "")
	return s.AuthRepository.StartDeviceGrant(client_name, scopes)
}

// WaitForApproval polls every grant.Interval seconds (at least two) until
// the sign-in is approved — returning the token —, denied ([ErrDenied]) or
// expired ([ErrExpired]); an error of the poll itself is returned as it is.
func (s AuthService) WaitForApproval(grant DeviceGrantModel) (string, error) {
	sleep := s.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	interval := time.Duration(grant.Interval) * time.Second
	if interval < minimumPollInterval {
		interval = minimumPollInterval
	}
	waited := time.Duration(0)
	limit := time.Duration(grant.ExpiresIn) * time.Second
	for waited < limit {
		sleep(interval)
		waited += interval
		poll, err := s.AuthRepository.PollDeviceGrant(grant.DeviceCode)
		if err != nil {
			return "", err
		}
		switch poll.State {
		case TokenApproved:
			if poll.Token != "" {
				return poll.Token, nil
			}
		case TokenDenied:
			return "", ErrDenied
		case TokenExpired:
			return "", ErrExpired
		}
	}
	return "", ErrExpired
}

// FinishLogin points the client at server_url with token, asks who it is
// (an Admin API error, 401 included, is returned and nothing is saved) and
// saves the current profile with the label and scopes the server reported.
func (s AuthService) FinishLogin(server_url string, token string, name string) (*ProfileModel, error) {
	server_url = strings.TrimRight(server_url, "/")
	s.AuthRepository.Use(server_url, token)
	principal, err := s.AuthRepository.Whoami()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = ProfileName(server_url)
	}
	profile := ProfileModel{Name: name, URL: server_url, Token: token, Scopes: principal.Scopes}
	if principal.Label != nil {
		profile.Label = *principal.Label
	}
	return &profile, s.save(profile)
}

// LoginLocal reads the token file (its errors, [ErrNoLocalToken] among them,
// are returned), checks the token at server_url and saves the current
// profile name ("local" when empty) without the token: it is read again on
// every run, so a reinstall that changes it does not break the profile.
func (s AuthService) LoginLocal(server_url string, name string) (*ProfileModel, error) {
	token, err := s.LocalTokenRepository.Read()
	if err != nil {
		return nil, err
	}
	server_url = strings.TrimRight(server_url, "/")
	s.AuthRepository.Use(server_url, token)
	principal, err := s.AuthRepository.Whoami()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = "local"
	}
	profile := ProfileModel{Name: name, URL: server_url, TokenFile: true, Scopes: principal.Scopes}
	if principal.Label != nil {
		profile.Label = *principal.Label
	}
	return &profile, s.save(profile)
}

// Logout removes a profile. With revoke, and when the profile holds its own
// token (not a Local token file), the token is revoked on the server first;
// a server that cannot be reached or refuses is reported and the profile is
// kept, so the token is not forgotten while still valid.
func (s AuthService) Logout(name string, revoke bool) (*ProfileModel, error) {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = config.Current
	}
	profile, ok := config.Profiles[name]
	if !ok {
		return nil, ErrNoProfile
	}
	if revoke && !profile.TokenFile && profile.Token != "" {
		s.AuthRepository.Use(profile.URL, profile.Token)
		if err = s.AuthRepository.Revoke(); err != nil {
			return nil, err
		}
	}
	delete(config.Profiles, name)
	if config.Current == name {
		config.Current = ""
	}
	return &profile, s.ProfileRepository.Save(config)
}

// Whoami asks the server about the token the client carries.
func (s AuthService) Whoami() (*PrincipalModel, error) {
	return s.AuthRepository.Whoami()
}

// save stores profile as the current profile.
func (s AuthService) save(profile ProfileModel) error {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return err
	}
	config.Profiles[profile.Name] = profile
	config.Current = profile.Name
	return s.ProfileRepository.Save(config)
}

// ProfileName makes a profile name from a server URL: "local" for a
// loopback address, otherwise the host name ("games.example.org"), with the
// port when it has one.
func ProfileName(server_url string) string {
	if IsLoopback(server_url) {
		return "local"
	}
	parsed, err := url.Parse(server_url)
	if err != nil || parsed.Host == "" {
		return server_url
	}
	return parsed.Host
}

// IsLoopback reports whether server_url points at this machine: localhost,
// 127.0.0.0/8 or ::1.
func IsLoopback(server_url string) bool {
	parsed, err := url.Parse(server_url)
	if err != nil {
		return false
	}
	host := parsed.Hostname()
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
