package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
	"gopkg.in/yaml.v3"
)

var ExitFunc = os.Exit

// outputData handles printing various data types in the requested format (CSV, JSON, YAML, or plain)
func outputData(figs figtree.Plant, data interface{}) {
	text := outputAsData(figs, data)
	fmt.Println(text)
}

func outputAsData(figs figtree.Plant, data interface{}) string {
	var sb strings.Builder
	j, y, c := *figs.Bool(argAsJson), *figs.Bool(argAsYaml), *figs.Bool(argAsCsv)
	if j { // just
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}
		sb.Write(jsonData)
		sb.WriteString("\n")
		return sb.String()
	}
	if y { // yeshua
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("Error marshaling to YAML: %v", err)
		}
		sb.Write(yamlData)
		sb.WriteString("\n")
		return sb.String()
	}
	if c { // calling
		switch v := data.(type) {
		case []string:
			sb.WriteString(strings.Join(v, ",") + "\n")
		case map[string]string:
			var parts []string
			for k, val := range v {
				parts = append(parts, fmt.Sprintf("%s=%s", k, val))
			}
			sb.WriteString(strings.Join(parts, ","))
		default:
			sb.WriteString(fmt.Sprintf("%v", data))
		}
		sb.WriteString("\n")
		return sb.String()
	}

	// no formatting option
	switch v := data.(type) { // you liked V for Vendetta, right?
	case []string:
		sb.WriteString(strings.Join(v, ","))
	case map[string]string:
		for k, val := range v {
			sb.WriteString(fmt.Sprintf("%s=%s", k, val) + "\n")
		}
	default:
		sb.WriteString(fmt.Sprintf("%v", data))
	}
	sb.WriteString("\n")
	return sb.String()
}
