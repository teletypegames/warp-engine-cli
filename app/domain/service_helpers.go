package domain

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrNothingToChange is returned by an update given no attributes.
var ErrNothingToChange = errors.New("nothing to change: give at least one field")

// ResolveSoftwareID turns a software reference into its ID: a number is the
// ID itself (not checked here; the next request fails if it is wrong),
// anything else is looked up as an exact name through repository.
func ResolveSoftwareID(repository SoftwareRepositoryInterface, reference string) (int, error) {
	if reference == "" {
		return 0, errors.New("no software given")
	}
	if id, err := strconv.Atoi(reference); err == nil {
		return id, nil
	}
	software, err := repository.FindByName(reference)
	if err != nil {
		return 0, err
	}
	return software.ID, nil
}

// ParseID parses the numeric identifier id of a kind of record ("release");
// anything but a positive number is an error naming kind.
func ParseID(kind string, id string) (int, error) {
	value, err := strconv.Atoi(id)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%q is no %s ID: give a positive number", id, kind)
	}
	return value, nil
}
