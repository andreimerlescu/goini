package main

import (
	"log"
	"path/filepath"

	"github.com/andreimerlescu/figtree/v2"
	"github.com/go-ini/ini"
)

// Run is the primary program entry point.
func Run(figs figtree.Plant) {
	err := figs.Load()
	if err != nil {
		log.Fatal(err)
	}
	filePath := filepath.Join(*figs.String(argIniFile))

	if *figs.String(argIniFile) == "" {
		log.Println("Error: --ini file path is required.")
		ExitFunc(1)
	}

	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("Error loading INI file: %v", err)
	}

	executePrimaryTask(figs, cfg, filePath)
}
