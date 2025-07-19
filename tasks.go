package main

import (
	"log"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
	"github.com/go-ini/ini"
)

// executePrimaryTask determines which primary task to run based on flags.
func executePrimaryTask(figs figtree.Plant, cfg *ini.File, filePath string) {
	if *figs.String(argAddSection) != "" {
		ignore(Task7(figs, cfg, filePath))
		return
	}
	if *figs.Bool(argAddKey) {
		ignore(Task8(figs, cfg, filePath))
		return
	}
	if *figs.Bool(argModifyKey) {
		ignore(Task9(figs, cfg, filePath))
		return
	}
	if *figs.String(argSectionName) != "" && *figs.String(argHasSectionKey) != "" {
		ignore(Task1(figs, cfg))
		return
	}
	if *figs.String(argSectionName) != "" && *figs.String(argKeyName) != "" && *figs.String(argHasSectionKeyValue) != "" {
		ignore(Task2(figs, cfg))
		return
	}
	if *figs.String(argAreSectionsPresent) != "" {
		ignore(Task3(figs, cfg))
		return
	}
	if *figs.Bool(argPrintSections) {
		ignore(Task4(figs, cfg))
		return
	}
	if *figs.Bool(argListKeys) {
		ignore(Task5(figs, cfg))
		return
	}
	if *figs.Bool(argListKeyValues) {
		ignore(Task6(figs, cfg))
		return
	}
	if *figs.String(argHasSection) != "" {
		ignore(TaskX(figs, cfg))
		return
	}
	ExitFunc(0)
}

// Task1 relying on exit code, does "section" have "key"? ---
func Task1(figs figtree.Plant, cfg *ini.File) error {
	sectionName, keyName := *figs.String(argSectionName), *figs.String(argHasSectionKey)
	if sectionName != "" && keyName != "" {
		section, err := cfg.GetSection(sectionName)
		if err != nil {
			ExitFunc(1) // section not found
			return nil
		}
		if section.HasKey(keyName) {
			ExitFunc(0)
		} else {
			ExitFunc(1) // key not found
		}
	}
	return nil
}

// Task2 relying on exit code, does "section" "key" have value "value"? ---
func Task2(figs figtree.Plant, cfg *ini.File) error {
	if *figs.String(argSectionName) != "" && *figs.String(argKeyName) != "" && *figs.String(argHasSectionKeyValue) != "" {
		section, err := cfg.GetSection(*figs.String(argSectionName))
		if err != nil {
			ExitFunc(1) // Section not found
			return nil
		}
		key, err := section.GetKey(*figs.String(argKeyName))
		if err != nil {
			ExitFunc(1) // Key not found
			return nil
		}
		if strings.EqualFold(key.String(), *figs.String(argHasSectionKeyValue)) {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}

// Task3 relying on exit code, are "section1" and "section2" present? ---
func Task3(figs figtree.Plant, cfg *ini.File) error {
	if *figs.String(argAreSectionsPresent) != "" {
		sectionsToCheck := strings.Split(*figs.String(argAreSectionsPresent), ",")
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
func Task4(figs figtree.Plant, cfg *ini.File) error {
	if *figs.Bool(argPrintSections) {
		sections := make([]string, 0)
		for _, section := range cfg.Sections() {
			if section.Name() == "DEFAULT" && len(section.Keys()) == 0 && !cfg.HasSection("default") {
				continue // Skip empty, synthetic DEFAULT if no explicit [default]
			}
			sections = append(sections, section.Name())
		}
		outputData(figs, sections)
		ExitFunc(0)
	}
	return nil
}

// Task5 using STDOUT, return a list of keys in "section" (by name) ---
func Task5(figs figtree.Plant, cfg *ini.File) error {
	if *figs.Bool(argListKeys) {
		if *figs.String(argSectionName) == "" {
			log.Println("Error: --section is required for --list-keys.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*figs.String(argSectionName))
		if err != nil {
			log.Printf("Error: Section '%s' not found.", *figs.String(argSectionName))
			ExitFunc(1)
			return nil
		}
		keys := make([]string, 0)
		for _, key := range section.Keys() {
			keys = append(keys, key.Name())
		}
		outputData(figs, keys)
		ExitFunc(0)
	}
	return nil
}

// Task6 using STDOUT, return a list of key/values in "section" (by name) ---
func Task6(figs figtree.Plant, cfg *ini.File) error {
	if *figs.Bool(argListKeyValues) {
		if *figs.String(argSectionName) == "" {
			log.Println("Error: --section is required for --list-key-values.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*figs.String(argSectionName))
		if err != nil {
			log.Printf("Error: Section '%s' not found.", *figs.String(argSectionName))
			ExitFunc(1)
			return nil
		}
		keyValuePairs := make(map[string]string)
		for _, key := range section.Keys() {
			keyValuePairs[key.Name()] = key.String()
		}
		outputData(figs, keyValuePairs)
		ExitFunc(0)
	}
	return nil
}

// Task7 using exit code for success status, add new section to ini file ---
func Task7(figs figtree.Plant, cfg *ini.File, filePath string) error {
	if *figs.String(argAddSection) != "" {
		if cfg.HasSection(*figs.String(argAddSection)) {
			log.Printf("Section '%s' already exists.", *figs.String(argAddSection))
			ExitFunc(1)
			return nil
		}
		_, err := cfg.NewSection(*figs.String(argAddSection))
		if err != nil {
			log.Printf("Error adding section '%s': %v", *figs.String(argAddSection), err)
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
func Task8(figs figtree.Plant, cfg *ini.File, filePath string) error {
	if *figs.Bool(argAddKey) {
		if *figs.String(argSectionName) == "" || *figs.String(argKeyName) == "" || *figs.String(argKeyValue) == "" {
			log.Println("Error: --section, --key, and --value are required for --add-key.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*figs.String(argSectionName))
		if err != nil {
			log.Printf("Error: Section '%s' not found. Cannot add key.", *figs.String(argSectionName))
			ExitFunc(1)
			return nil
		}
		if section.HasKey(*figs.String(argKeyName)) {
			log.Printf("Key '%s' already exists in section '%s'. Use --modify-key to change its value.", *figs.String(argKeyName), *figs.String(argHasSectionKey))
			ExitFunc(1)
			return nil
		}
		section.Key(*figs.String(argKeyName)).SetValue(*figs.String(argKeyValue))
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
func Task9(figs figtree.Plant, cfg *ini.File, filePath string) error {
	if *figs.Bool(argModifyKey) {
		if *figs.String(argSectionName) == "" || *figs.String(argKeyName) == "" || *figs.String(argKeyValue) == "" {
			log.Println("Error: --section, --key, and --value are required for --modify-key.")
			ExitFunc(1)
			return nil
		}
		section, err := cfg.GetSection(*figs.String(argSectionName))
		if err != nil {
			log.Printf("Error: Section '%s' not found. Cannot modify key.", *figs.String(argSectionName))
			ExitFunc(1)
			return nil
		}
		if !section.HasKey(*figs.String(argKeyName)) {
			log.Printf("Key '%s' does not exist in section '%s'. Use --add-key to add it.", *figs.String(argKeyName), *figs.String(argHasSectionKey))
			ExitFunc(1)
			return nil
		}
		section.Key(*figs.String(argKeyName)).SetValue(*figs.String(argKeyValue))
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
func TaskX(figs figtree.Plant, cfg *ini.File) error {
	if *figs.String(argHasSection) != "" {
		if cfg.HasSection(*figs.String(argHasSection)) {
			ExitFunc(0)
		} else {
			ExitFunc(1)
		}
	}
	return nil
}
