package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var ExitFunc = os.Exit

// outputData handles printing various data types in the requested format (CSV, JSON, YAML, or plain)
func outputData(pCfg *ProgramConfig, data interface{}) {
	if *pCfg.AsJson {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err) // log.Fatalf calls ExitFunc(1) implicitly
		}
		ignoreFprint(fmt.Fprintln(os.Stdout, string(jsonData)))
	} else if *pCfg.AsYaml {
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("Error marshaling to YAML: %v", err) // log.Fatalf calls ExitFunc(1) implicitly
		}
		ignoreFprint(fmt.Fprintln(os.Stdout, string(yamlData)))
	} else if *pCfg.AsCsv {
		switch v := data.(type) {
		case []string:
			ignoreFprint(fmt.Fprintln(os.Stdout, strings.Join(v, ",")))
		case map[string]string:
			var parts []string
			for k, val := range v {
				parts = append(parts, fmt.Sprintf("%s=%s", k, val))
			}
			ignoreFprint(fmt.Fprintln(os.Stdout, strings.Join(parts, ",")))
		default:
			ignoreFprint(fmt.Fprintln(os.Stdout, fmt.Sprintf("%v", data)))
		}
	} else {
		switch v := data.(type) {
		case []string:
			ignoreFprint(fmt.Fprintln(os.Stdout, strings.Join(v, "\n")))
		case map[string]string:
			for k, val := range v {
				ignoreFprint(fmt.Fprintf(os.Stdout, "%s = %s\n", k, val))
			}
		default:
			ignoreFprint(fmt.Fprintln(os.Stdout, fmt.Sprintf("%v", data)))
		}
	}
	ignoreFprint(fmt.Fprintln(os.Stdout, ""))
}
