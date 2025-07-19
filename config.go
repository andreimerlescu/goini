package main

import "github.com/andreimerlescu/configurable"

// ProgramConfig holds all the application's configuration flags.
// This struct will be initialized per test to ensure flag isolation.
type ProgramConfig struct {
	Config             configurable.IConfigurable
	IniFile            *string
	PrintSections      *bool
	HasSection         *string
	HasSectionKey      *string
	HasSectionKeyValue *string
	SectionName        *string
	KeyName            *string
	KeyValue           *string
	AreSectionsPresent *string
	ListKeys           *bool
	ListKeyValues      *bool
	AddSection         *string
	AddKey             *bool
	ModifyKey          *bool
	AsCsv              *bool
	AsJson             *bool
	AsYaml             *bool
}

// NewProgramConfig initializes a new set of configurable flags.
func NewProgramConfig() *ProgramConfig {
	cfg := configurable.New()
	return &ProgramConfig{
		Config:             cfg,
		IniFile:            cfg.NewString("ini", "", "Path to ini file to process"),
		PrintSections:      cfg.NewBool("sections", false, "If true, only the sections from the ini file will be displayed"),
		HasSection:         cfg.NewString("has-section", "", "If set, exit code will respond if the section exists or not in the --ini file"),
		HasSectionKey:      cfg.NewString("has-section-key", "", "If set, exit code will respond if the section has the key in the --ini file. Requires --section."),
		HasSectionKeyValue: cfg.NewString("has-section-key-value", "", "If set, exit code will respond if the section key has the specified value in the --ini file. Requires --section and --key."),
		SectionName:        cfg.NewString("section", "", "Specify section name for operations like --has-section-key, --list-keys, --list-key-values, --add-key, --modify-key."),
		KeyName:            cfg.NewString("key", "", "Specify key name for operations like --has-section-key-value, --add-key, --modify-key."),
		KeyValue:           cfg.NewString("value", "", "Specify value for operations like --has-section-key-value, --add-key, --modify-key."),
		AreSectionsPresent: cfg.NewString("are-sections-present", "", "Comma-separated list of section names. Exit code will be 0 if all are present, 1 otherwise."),
		ListKeys:           cfg.NewBool("list-keys", false, "If true, returns a list of keys in the specified --section using STDOUT."),
		ListKeyValues:      cfg.NewBool("list-key-values", false, "If true, returns a list of key/values in the specified --section using STDOUT."),
		AddSection:         cfg.NewString("add-section", "", "If set, adds a new section to the ini file. Exit code 0 on success."),
		AddKey:             cfg.NewBool("add-key", false, "If true, adds a new key-value pair to the specified --section. Requires --key and --value. Exit code 0 on success."),
		ModifyKey:          cfg.NewBool("modify-key", false, "If true, modifies an existing key's value in the specified --section. Requires --key and --value. Exit code 0 on success."),
		AsCsv:              cfg.NewBool("csv", false, "Output as CSV."),
		AsJson:             cfg.NewBool("json", false, "Output as JSON."),
		AsYaml:             cfg.NewBool("yaml", false, "Output as YAML."),
	}
}
