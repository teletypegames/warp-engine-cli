package cli

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// LocalProviderSerializer turns the providers of WarpEngine Local and
// their sign-ins into tables and records.
type LocalProviderSerializer struct{}

// ModelToTable renders providers as a table: key, label, role, state,
// requires, message.
func (s LocalProviderSerializer) ModelToTable(providers []domain.LocalProviderModel) Table {
	table := Table{Headers: []string{"Key", "Label", "Role", "State", "Requires", "Message"}}
	for _, provider := range providers {
		table.Rows = append(table.Rows, []string{
			provider.Key, provider.Label, provider.Role, provider.State,
			List(provider.Requires), Truncate(Text(provider.Message), 50),
		})
	}
	return table
}

// CredentialsToRecord renders the sign-in of a provider: URL, user,
// password.
func (s LocalProviderSerializer) CredentialsToRecord(credentials domain.LocalCredentialsModel) Record {
	return Record{
		{Name: "URL", Value: credentials.URL},
		{Name: "User", Value: credentials.User},
		{Name: "Password", Value: credentials.Password},
	}
}
