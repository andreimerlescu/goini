package main

import (
	"github.com/andreimerlescu/figtree/v2"
)

const (
	argIniFile            string = "ini"
	argPrintSections      string = "sections"
	argHasSection         string = "has-section"
	argHasSectionKey      string = "has-section-key"
	argHasSectionKeyValue string = "has-section-key-value"
	argSectionName        string = "section"
	argKeyName            string = "key"
	argKeyValue           string = "value"
	argAreSectionsPresent string = "are-sections-present"
	argListKeys           string = "list-keys"
	argListKeyValues      string = "list-key-values"
	argAddSection         string = "add-section"
	argAddKey             string = "add-key"
	argModifyKey          string = "modify-key"
	argAsCsv              string = "csv"
	argAsJson             string = "json"
	argAsYaml             string = "yaml"
)

// NewConfiguration initializes a new set of configurable flags
func NewConfiguration() figtree.Plant {
	figs := figtree.With(figtree.Options{
		IgnoreEnvironment: true,
		Germinate:         true,
	})

	figs = figs.NewString(argIniFile, "", "Path to ini file to process")

	// sections
	figs = figs.NewBool(argPrintSections, false, "If true, only the sections from the ini file will be displayed")
	figs = figs.NewString(argHasSection, "", "If set, exit code will respond if the section exists or not in the --ini file")
	figs = figs.NewString(argHasSectionKey, "", "If set, exit code will respond if the section has the key in the --ini file. Requires --section.")
	figs = figs.NewString(argHasSectionKeyValue, "", "If set, exit code will respond if the section key has the specified value in the --ini file. Requires --section and --key.")
	figs = figs.NewString(argSectionName, "", "Specify section name for operations like --has-section-key, --list-keys, --list-key-values, --add-key, --modify-key.")
	figs = figs.NewString(argKeyName, "", "Specify key name for operations like --has-section-key-value, --add-key, --modify-key.")
	figs = figs.NewString(argKeyValue, "", "Specify value for operations like --has-section-key-value, --add-key, --modify-key.")
	figs = figs.NewString(argAreSectionsPresent, "", "Comma-separated list of section names. Exit code will be 0 if all are present, 1 otherwise.")
	figs = figs.NewString(argAddSection, "", "If set, adds a new section to the ini file. Exit code 0 on success.")

	// lists
	figs = figs.NewBool(argListKeys, false, "If true, returns a list of keys in the specified --section using STDOUT.")
	figs = figs.NewBool(argListKeyValues, false, "If true, returns a list of key/values in the specified --section using STDOUT.")

	// keys
	figs = figs.NewBool(argAddKey, false, "If true, adds a new key-value pair to the specified --section. Requires --key and --value. Exit code 0 on success.")
	figs = figs.NewBool(argModifyKey, false, "If true, modifies an existing key's value in the specified --section. Requires --key and --value. Exit code 0 on success.")

	// output
	figs = figs.NewBool(argAsCsv, false, "Output as CSV.")
	figs = figs.NewBool(argAsJson, false, "Output as JSON.")
	figs = figs.NewBool(argAsYaml, false, "Output as YAML.")
	return figs
}
