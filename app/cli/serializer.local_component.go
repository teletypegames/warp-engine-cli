package cli

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// LocalComponentSerializer turns the components of WarpEngine Local and
// their sign-ins into tables and records.
type LocalComponentSerializer struct{}

// ModelToTable renders components as a table: key, label, role, state,
// requires, message.
func (s LocalComponentSerializer) ModelToTable(components []domain.LocalComponentModel) Table {
	table := Table{Headers: []string{"Key", "Label", "Role", "State", "Requires", "Message"}}
	for _, component := range components {
		table.Rows = append(table.Rows, []string{
			component.Key, component.Label, component.Role, component.State,
			List(component.Requires), Truncate(Text(component.Message), 50),
		})
	}
	return table
}

// CredentialsToRecord renders the sign-in of a component: URL, user,
// password.
func (s LocalComponentSerializer) CredentialsToRecord(credentials domain.LocalCredentialsModel) Record {
	return Record{
		{Name: "URL", Value: credentials.URL},
		{Name: "User", Value: credentials.User},
		{Name: "Password", Value: credentials.Password},
	}
}
