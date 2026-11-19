package cli

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// ProfileSerializer renders the saved profiles.
type ProfileSerializer struct{}

// ModelToTable renders profiles; the current one is marked with "*".
func (s ProfileSerializer) ModelToTable(profiles []domain.ProfileModel, current string) Table {
	table := Table{Headers: []string{"Name", "Current", "URL", "Signed in as", "Scopes"}}
	for _, profile := range profiles {
		marker := ""
		if profile.Name == current {
			marker = "*"
		}
		table.Rows = append(table.Rows, []string{profile.Name, marker, profile.URL, labelOf(&profile), List(profile.Scopes)})
	}
	return table
}

// ModelToRecord renders one profile; the token itself is never printed,
// only whether there is one.
func (s ProfileSerializer) ModelToRecord(profile domain.ProfileModel, current string) Record {
	token := "saved"
	switch {
	case profile.TokenFile:
		token = "the WarpEngine Local token file"
	case profile.Token == "":
		token = "none"
	}
	return Record{
		{Name: "Name", Value: profile.Name},
		{Name: "Current", Value: YesNo(profile.Name == current)},
		{Name: "URL", Value: profile.URL},
		{Name: "Signed in as", Value: labelOf(&profile)},
		{Name: "Token", Value: token},
		{Name: "Scopes", Value: List(profile.Scopes)},
	}
}

// ModelToResponse returns profile without its token, for the JSON output.
func (s ProfileSerializer) ModelToResponse(profile domain.ProfileModel) domain.ProfileModel {
	profile.Token = ""
	return profile
}

// ModelToResponseMany is [ProfileSerializer.ModelToResponse] for a list.
func (s ProfileSerializer) ModelToResponseMany(profiles []domain.ProfileModel) []domain.ProfileModel {
	result := make([]domain.ProfileModel, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, s.ModelToResponse(profile))
	}
	return result
}
