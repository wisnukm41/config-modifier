package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

func main() {
	// Parse command line arguments
	key := flag.String("k", "", "Key to update")
	value := flag.String("v", "", "New value for the key")
	flag.Parse()

	if *key == "" || *value == "" {
		fmt.Println("Please provide both key and value to update")
		return
	}

	// Get the current working directory
	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Find all YAML files in the directory
	yamlFilePaths, err := findYAMLFiles(currentDirectory)
	if err != nil {
		fmt.Println("Error finding YAML files:", err)
		return
	}

	// Update each YAML file
	for _, filepath := range yamlFilePaths {
		if err := updateYAMLFile(filepath, *key, *value); err != nil {
			fmt.Println("Failed to update file:", filepath, err)
		} else {
			fmt.Printf("Successfully updated key `%s` to value `%s` in %s\n", *key, *value, filepath)
		}
	}
}

// Recursively find all YAML files in the given directory
func findYAMLFiles(dir string) ([]string, error) {
	var yamlFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})
	return yamlFiles, err
}

// Update the key-value pair in the given YAML file
func updateYAMLFile(filepath string, key, value string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	var yamlNode yaml.Node
	if err := yaml.Unmarshal(data, &yamlNode); err != nil {
		return err
	}

	updateKeyValue(&yamlNode, key, value)

	updatedData, err := yaml.Marshal(&yamlNode)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, updatedData, 0644)
}

// Recursively update the key-value pair in the YAML node
func updateKeyValue(node *yaml.Node, targetKey, newValue string) {
	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			updateKeyValue(child, targetKey, newValue)
		}
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if k.Value == targetKey {
				// Update the value while preserving its type
				if v.Tag == "!!int" {
					v.Value = newValue
					v.Tag = "!!int"
				} else if v.Tag == "!!bool" {
					v.Value = newValue
					v.Tag = "!!bool"
				} else if v.Tag == "!!float" {
					v.Value = newValue
					v.Tag = "!!float"
				} else {
					v.Value = newValue
					v.Tag = "!!str"
				}
				return
			}
			updateKeyValue(v, targetKey, newValue)
		}
	}
	if node.Kind == yaml.SequenceNode {
		for _, child := range node.Content {
			updateKeyValue(child, targetKey, newValue)
		}
	}
}

