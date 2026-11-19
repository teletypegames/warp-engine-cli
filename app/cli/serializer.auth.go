package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// PrincipalSerializer renders the answer of whoami.
type PrincipalSerializer struct{}

// ModelToRecord renders principal, and the server and profile it was asked
// through when profile is not nil.
func (s PrincipalSerializer) ModelToRecord(principal domain.PrincipalModel, profile *domain.ProfileModel) Record {
	record := Record{{Name: "Label", Value: Text(principal.Label)}}
	if profile != nil {
		record = append(record, Field{Name: "Server", Value: profile.URL}, Field{Name: "Profile", Value: profile.Name})
	}
	owner := Text(principal.OwnerType)
	if principal.OwnerID != nil {
		owner += " #" + strconv.Itoa(*principal.OwnerID)
	}
	return append(record,
		Field{Name: "Owner", Value: owner},
		Field{Name: "Token", Value: Text(principal.TokenPrefix)},
		Field{Name: "Scopes", Value: List(principal.Scopes)},
		Field{Name: "Unrestricted", Value: YesNo(principal.Unrestricted)},
		Field{Name: "Expires", Value: Text(principal.ExpiresAt)},
		Field{Name: "Resources", Value: List(principal.Resources)},
		Field{Name: "Read-only", Value: YesNo(principal.ReadOnly)},
		Field{Name: "MCP", Value: YesNo(principal.MCP)},
	)
}

// ServerSerializer renders the answer of "warp-engine server".
type ServerSerializer struct{}

// ModelToRecord renders server, asked at server_url. An engine older than
// 5.2 says nothing about itself: its fields print as "-".
func (s ServerSerializer) ModelToRecord(server domain.ServerModel, server_url string) Record {
	return Record{
		{Name: "URL", Value: server_url},
		{Name: "Name", Value: Text(&server.Name)},
		{Name: "Engine", Value: Text(&server.Engine)},
		{Name: "Version", Value: Text(&server.Version)},
		{Name: "Catalog API", Value: strconv.Itoa(server.API)},
		{Name: "Profile", Value: Text(&server.Profile)},
		{Name: "Sign-in", Value: List(server.SignIn)},
		{Name: "Admin API", Value: Text(&server.AdminAPI)},
		{Name: "Admin UI", Value: YesNo(server.Admin)},
	}
}
