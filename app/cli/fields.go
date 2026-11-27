package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// Kinds of [FieldFlag]: how the flag's text becomes the API's value.
const (
	// KindString sends the text as it is.
	KindString = "string"
	// KindNullable sends the text, or null for an empty text: --author ""
	// clears the author.
	KindNullable = "nullable"
	// KindInt sends a whole number.
	KindInt = "int"
	// KindNullableInt sends a whole number, or null for an empty text.
	KindNullableInt = "nullable-int"
	// KindBool sends true or false (--highlighted, --highlighted=false).
	KindBool = "bool"
	// KindStrings sends a list: the flag repeats or takes commas
	// (--scopes admin:read,admin:write).
	KindStrings = "strings"
)

// FieldFlag ties a command-line flag to a field of the API: the serializers
// list their fields with it, and [DeclareFields] and [FieldsToAttributes]
// do the rest.
type FieldFlag struct {
	// Flag is the flag name, kebab-case ("directory-name").
	Flag string
	// Key is the API field it sets, camelCase ("directoryName").
	Key string
	// Kind is one of the Kind constants; empty is [KindString].
	Kind string
	// Usage is the help text of the flag.
	Usage string
	// Enum, when set, lists the values the flag accepts; the help shows them.
	Enum []string
}

// DataFlag is the flag every create and update command has besides its
// field flags: a JSON object of attributes, inline or "@file.json", merged
// under the field flags. It reaches fields that have no flag of their own
// (nested objects such as a sync peer's transport options).
const DataFlag = "data"

// DeclareFields adds a flag for every field of fields to flags, and the
// --data flag.
func DeclareFields(flags *pflag.FlagSet, fields []FieldFlag) {
	for _, field := range fields {
		usage := field.Usage
		if len(field.Enum) > 0 {
			usage += " (" + strings.Join(field.Enum, ", ") + ")"
		}
		switch field.Kind {
		case KindBool:
			flags.Bool(field.Flag, false, usage)
		case KindStrings:
			flags.StringSlice(field.Flag, nil, usage)
		case KindNullable, KindNullableInt:
			flags.String(field.Flag, "", usage+"; \"\" clears it")
		default:
			flags.String(field.Flag, "", usage)
		}
	}
	flags.String(DataFlag, "", "attributes as a JSON object, or @file.json; the flags win over it")
}

// FieldsToAttributes builds the attributes of the flags given on the
// command line: the --data object first, then every field flag that was set
// (an unset flag sends nothing, so an update leaves that field alone). A
// value outside a field's Enum, a malformed number or bad JSON is an error
// naming the flag.
func FieldsToAttributes(flags *pflag.FlagSet, fields []FieldFlag) (domain.Attributes, error) {
	attributes := domain.Attributes{}
	if data, _ := flags.GetString(DataFlag); data != "" {
		parsed, err := readData(data)
		if err != nil {
			return nil, err
		}
		for key, value := range parsed {
			attributes[key] = value
		}
	}
	for _, field := range fields {
		if !flags.Changed(field.Flag) {
			continue
		}
		value, err := fieldValue(flags, field)
		if err != nil {
			return nil, err
		}
		attributes[field.Key] = value
	}
	return attributes, nil
}

// fieldValue reads the value of one set flag in the API's type.
func fieldValue(flags *pflag.FlagSet, field FieldFlag) (any, error) {
	switch field.Kind {
	case KindBool:
		return flags.GetBool(field.Flag)
	case KindStrings:
		values, err := flags.GetStringSlice(field.Flag)
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			if err = checkEnum(field, value); err != nil {
				return nil, err
			}
		}
		return values, nil
	}
	text, err := flags.GetString(field.Flag)
	if err != nil {
		return nil, err
	}
	if text == "" && (field.Kind == KindNullable || field.Kind == KindNullableInt) {
		return nil, nil
	}
	if field.Kind == KindInt || field.Kind == KindNullableInt {
		var number int
		if _, err = fmt.Sscanf(text, "%d", &number); err != nil || fmt.Sprint(number) != text {
			return nil, &UsageError{Err: fmt.Errorf("--%s takes a whole number, not %q", field.Flag, text)}
		}
		return number, nil
	}
	if err = checkEnum(field, text); err != nil {
		return nil, err
	}
	return text, nil
}

// checkEnum returns a [*UsageError] when field has an Enum without value.
func checkEnum(field FieldFlag, value string) error {
	if len(field.Enum) == 0 {
		return nil
	}
	for _, allowed := range field.Enum {
		if allowed == value {
			return nil
		}
	}
	return &UsageError{Err: fmt.Errorf("--%s takes one of %s, not %q", field.Flag, strings.Join(field.Enum, ", "), value)}
}

// readData parses the --data value: a JSON object, or "@path" to read one
// from a file ("@-" reads standard input).
func readData(data string) (map[string]any, error) {
	raw := []byte(data)
	if strings.HasPrefix(data, "@") {
		var err error
		if data == "@-" {
			raw, err = readAll(os.Stdin)
		} else {
			raw, err = os.ReadFile(strings.TrimPrefix(data, "@"))
		}
		if err != nil {
			return nil, fmt.Errorf("--%s: %w", DataFlag, err)
		}
	}
	parsed := map[string]any{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, &UsageError{Err: fmt.Errorf("--%s is no JSON object: %w", DataFlag, err)}
	}
	return parsed, nil
}
