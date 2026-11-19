package domain

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNoProfile is returned when a profile is needed but none is saved and
// none is current.
var ErrNoProfile = errors.New("no profile: sign in with \"warp-engine login <url>\"")

// ProfileServiceInterface manages the saved profiles.
type ProfileServiceInterface interface {
	// All returns the profiles sorted by name, and the name of the current one.
	All() (*[]ProfileModel, string, error)
	// Show returns the profile name, or the current one when name is empty.
	Show(name string) (*ProfileModel, error)
	// Save stores profile under its name, replacing one with the same name;
	// with current set it becomes the current profile.
	Save(profile ProfileModel, current bool) error
	// Use makes the profile name the current one.
	Use(name string) error
	// Destroy removes the profile name; when it was current, no profile is.
	Destroy(name string) error
}

// ProfileService is the [ProfileServiceInterface] over a [ProfileRepositoryInterface].
type ProfileService struct {
	// ProfileRepository keeps the configuration file.
	ProfileRepository ProfileRepositoryInterface
}

// All loads the configuration and lists its profiles by name.
func (s ProfileService) All() (*[]ProfileModel, string, error) {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return nil, "", err
	}
	profiles := make([]ProfileModel, 0, len(config.Profiles))
	for _, profile := range config.Profiles {
		profiles = append(profiles, profile)
	}
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return &profiles, config.Current, nil
}

// Show returns one profile; [ErrNoProfile] when name is empty and there is
// no current profile, an error naming it when there is no such profile.
func (s ProfileService) Show(name string) (*ProfileModel, error) {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = config.Current
	}
	if name == "" {
		return nil, ErrNoProfile
	}
	profile, ok := config.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("no profile named %q (see \"warp-engine profile list\")", name)
	}
	return &profile, nil
}

// Save stores profile; an empty name is an error.
func (s ProfileService) Save(profile ProfileModel, current bool) error {
	if profile.Name == "" {
		return errors.New("a profile needs a name")
	}
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return err
	}
	config.Profiles[profile.Name] = profile
	if current || config.Current == "" {
		config.Current = profile.Name
	}
	return s.ProfileRepository.Save(config)
}

// Use sets the current profile; an unknown name is an error.
func (s ProfileService) Use(name string) error {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return err
	}
	if _, ok := config.Profiles[name]; !ok {
		return fmt.Errorf("no profile named %q (see \"warp-engine profile list\")", name)
	}
	config.Current = name
	return s.ProfileRepository.Save(config)
}

// Destroy deletes a profile; an unknown name is an error.
func (s ProfileService) Destroy(name string) error {
	config, err := s.ProfileRepository.Load()
	if err != nil {
		return err
	}
	if _, ok := config.Profiles[name]; !ok {
		return fmt.Errorf("no profile named %q (see \"warp-engine profile list\")", name)
	}
	delete(config.Profiles, name)
	if config.Current == name {
		config.Current = ""
	}
	return s.ProfileRepository.Save(config)
}
