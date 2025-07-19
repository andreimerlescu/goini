package main

import (
	"log"
	"path/filepath"

	"github.com/go-ini/ini"
)

// Run is the primary program entry point.
func Run(pCfg *ProgramConfig) {
	// Parse command-line arguments. This should happen once at the start.
	// We pass os.Args as the source, configurable library parses it
	_ = pCfg.Config.Parse("")
	filePath := filepath.Join(*pCfg.IniFile)

	// Input validation for INI file path
	if *pCfg.IniFile == "" {
		log.Println("Error: --ini file path is required.")
		ExitFunc(1)
	}

	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("Error loading INI file: %v", err) // log.Fatalf calls ExitFunc(1) implicitly
	}

	// Pass the specific ProgramConfig to the task executor
	executePrimaryTask(pCfg, cfg, filePath)
}
