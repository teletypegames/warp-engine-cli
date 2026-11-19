package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ExtraSerializer turns the flags of the extra commands into the API's
// terms and an extra into a record.
type ExtraSerializer struct{}

// ExtraValueFlag is the flag of "extra update" that holds the new value: a
// JSON object, or an array of objects for a resource that holds many;
// inline, or "@file.json" ("@-" reads standard input).
const ExtraValueFlag = "value"

// DeclareFlags adds --value and --data to the flags of "extra update".
func (s ExtraSerializer) DeclareFlags(flags *pflag.FlagSet) {
	flags.String(ExtraValueFlag, "", "the new value: a JSON object (an array of objects when the resource holds many), or @file.json")
	DeclareFields(flags, nil)
}

// RequestToAttributes reads the body of "extra update": the --data object,
// with --value, when given, as its "value". A --value that is no JSON is a
// [*UsageError].
func (s ExtraSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	attributes, err := FieldsToAttributes(flags, nil)
	if err != nil {
		return nil, err
	}
	text, _ := flags.GetString(ExtraValueFlag)
	if !flags.Changed(ExtraValueFlag) {
		return attributes, nil
	}
	raw := []byte(text)
	if strings.HasPrefix(text, "@") {
		if text == "@-" {
			raw, err = readAll(os.Stdin)
		} else {
			raw, err = os.ReadFile(strings.TrimPrefix(text, "@"))
		}
		if err != nil {
			return nil, fmt.Errorf("--%s: %w", ExtraValueFlag, err)
		}
	}
	var value any
	if err = json.Unmarshal(raw, &value); err != nil {
		return nil, &UsageError{Err: fmt.Errorf("--%s is no JSON: %w", ExtraValueFlag, err)}
	}
	attributes["value"] = value
	return attributes, nil
}

// ModelToRecord renders an extra: the association first, whether it holds
// many, and the value as indented JSON.
func (s ExtraSerializer) ModelToRecord(extra domain.ExtraModel) Record {
	value := "-"
	if extra.Value != nil {
		if data, err := json.MarshalIndent(extra.Value, "", "  "); err == nil {
			value = string(data)
		}
	}
	return Record{
		{Name: "Association", Value: extra.Association},
		{Name: "Multiple", Value: YesNo(extra.Multiple)},
		{Name: "Value", Value: value},
	}
}
