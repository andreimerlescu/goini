package main

import (
	"log"
	"strings"

	"github.com/go-ini/ini"
)

// executePrimaryTask determines which primary task to run based on flags.
func executePrimaryTask(pCfg *ProgramConfig, cfg *ini.File, filePath string) {
	if *pCfg.AddSection != "" {
		ignore(Task7(pCfg, cfg, filePath))
		return
	}
	if *pCfg.AddKey {
		ignore(Task8(pCfg, cfg, filePath))
		return
	}
	if *pCfg.ModifyKey {
		ignore(Task9(pCfg, cfg, filePath))
		return
	}
	if *pCfg.SectionName != "" && *pCfg.HasSectionKey != "" {
		ignore(Task1(pCfg, cfg))
		// Task1 calls ExitFunc internally on success/failure
		return
	}
	if *pCfg.SectionName != "" && *pCfg.KeyName != "" && *pCfg.HasSectionKeyValue != "" {
		ignore(Task2(pCfg, cfg))
		return
	}
	if *pCfg.AreSectionsPresent != "" {
		ignore(Task3(pCfg, cfg))
		return
	}
	if *pCfg.PrintSections {
		ignore(Task4(pCfg, cfg))
		return
	}
	if *pCfg.ListKeys {
		ignore(Task5(pCfg, cfg))
		return
	}
	if *pCfg.ListKeyValues {
		ignore(Task6(pCfg, cfg))
		return
	}
	if *pCfg.HasSection != "" {
		ignore(TaskX(pCfg, cfg))
		return
	}
	ExitFunc(0)
}

// Task1 relying on exit code, does "section" have "key"? ---
func Task1(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.SectionName != "" && *pCfg.HasSectionKey != "" {
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			ExitFunc(1) // Section not found
			return nil  // Return to allow defer calls to execute
		}
		if section.HasKey(*pCfg.HasSectionKey) {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}

// Task2 relying on exit code, does "section" "key" have value "value"? ---
func Task2(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.SectionName != "" && *pCfg.KeyName != "" && *pCfg.HasSectionKeyValue != "" {
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			ExitFunc(1) // Section not found
			return nil
		}
		key, err := section.GetKey(*pCfg.KeyName)
		if err != nil {
			ExitFunc(1) // Key not found
			return nil
		}
		if key.String() == *pCfg.HasSectionKeyValue {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}

// Task3 relying on exit code, are "section1" and "section2" present? ---
func Task3(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.AreSectionsPresent != "" {
		sectionsToCheck := strings.Split(*pCfg.AreSectionsPresent, ",")
		allPresent := true
		for _, sec := range sectionsToCheck {
			if !cfg.HasSection(strings.TrimSpace(sec)) {
				allPresent = false
				break
			}
		}
		if allPresent {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}

// Task4 using STDOUT, return list of "sections" in ini file ---
func Task4(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.PrintSections {
		sections := make([]string, 0)
		for _, section := range cfg.Sections() {
			// Filter out the synthetic DEFAULT section if it's not explicitly in the INI.
			// The go-ini library creates a DEFAULT section for global keys or as a fallback.
			// If your INI file contains an explicit '[default]' section,
			// the uppercase 'DEFAULT' often appears as a separate, empty section or a global one.
			if section.Name() == "DEFAULT" && len(section.Keys()) == 0 && !cfg.HasSection("default") {
				continue // Skip empty, synthetic DEFAULT if no explicit [default]
			}
			sections = append(sections, section.Name())
		}
		outputData(pCfg, sections) // Pass pCfg to outputData
		ExitFunc(0)
	}
	return nil
}

// Task5 using STDOUT, return a list of keys in "section" (by name) ---
func Task5(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.ListKeys {
		if *pCfg.SectionName == "" {
			log.Println("Error: --section is required for --list-keys.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			log.Printf("Error: Section '%s' not found.", *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		keys := make([]string, 0)
		for _, key := range section.Keys() {
			keys = append(keys, key.Name())
		}
		outputData(pCfg, keys) // Pass pCfg to outputData
		ExitFunc(0)
	}
	return nil
}

// Task6 using STDOUT, return a list of key/values in "section" (by name) ---
func Task6(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.ListKeyValues {
		if *pCfg.SectionName == "" {
			log.Println("Error: --section is required for --list-key-values.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			log.Printf("Error: Section '%s' not found.", *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		keyValuePairs := make(map[string]string)
		for _, key := range section.Keys() {
			keyValuePairs[key.Name()] = key.String()
		}
		outputData(pCfg, keyValuePairs) // Pass pCfg to outputData
		ExitFunc(0)
	}
	return nil
}

// Task7 using exit code for success status, add new section to ini file ---
func Task7(pCfg *ProgramConfig, cfg *ini.File, filePath string) error {
	if *pCfg.AddSection != "" {
		if cfg.HasSection(*pCfg.AddSection) {
			log.Printf("Section '%s' already exists.", *pCfg.AddSection)
			ExitFunc(1)
			return nil
		}
		_, err := cfg.NewSection(*pCfg.AddSection)
		if err != nil {
			log.Printf("Error adding section '%s': %v", *pCfg.AddSection, err)
			ExitFunc(1)
			return nil
		}
		if err := cfg.SaveTo(filePath); err != nil {
			log.Printf("Error saving INI file: %v", err)
			ExitFunc(1)
			return nil
		}
		ExitFunc(0)
	}
	return nil
}

// Task8 using exit code for success status, in section "section", add "key" with value "value" ---
func Task8(pCfg *ProgramConfig, cfg *ini.File, filePath string) error {
	if *pCfg.AddKey {
		if *pCfg.SectionName == "" || *pCfg.KeyName == "" || *pCfg.KeyValue == "" {
			log.Println("Error: --section, --key, and --value are required for --add-key.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			log.Printf("Error: Section '%s' not found. Cannot add key.", *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		if section.HasKey(*pCfg.KeyName) {
			log.Printf("Key '%s' already exists in section '%s'. Use --modify-key to change its value.", *pCfg.KeyName, *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		section.Key(*pCfg.KeyName).SetValue(*pCfg.KeyValue)
		if err := cfg.SaveTo(filePath); err != nil {
			log.Printf("Error saving INI file: %v", err)
			ExitFunc(1)
			return nil
		}
		ExitFunc(0)
	}
	return nil
}

// Task9 using exit code for success status, in section "section", modify "key" with new value "value" ---
func Task9(pCfg *ProgramConfig, cfg *ini.File, filePath string) error {
	if *pCfg.ModifyKey {
		if *pCfg.SectionName == "" || *pCfg.KeyName == "" || *pCfg.KeyValue == "" {
			log.Println("Error: --section, --key, and --value are required for --modify-key.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*pCfg.SectionName)
		if err != nil {
			log.Printf("Error: Section '%s' not found. Cannot modify key.", *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		if !section.HasKey(*pCfg.KeyName) {
			log.Printf("Key '%s' does not exist in section '%s'. Use --add-key to add it.", *pCfg.KeyName, *pCfg.SectionName)
			ExitFunc(1)
			return nil
		}
		section.Key(*pCfg.KeyName).SetValue(*pCfg.KeyValue)
		if err := cfg.SaveTo(filePath); err != nil {
			log.Printf("Error saving INI file: %v", err)
			ExitFunc(1)
			return nil
		}
		ExitFunc(0)
	}
	return nil
}

// TaskX This is the original logic for --has-section.
func TaskX(pCfg *ProgramConfig, cfg *ini.File) error {
	if *pCfg.HasSection != "" {
		if cfg.HasSection(*pCfg.HasSection) {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}
