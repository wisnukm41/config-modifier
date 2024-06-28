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

	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}

	updateKeyValue(yamlData, key, value)

	updatedData, err := yaml.Marshal(yamlData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, updatedData, 0644)
}

// Recursively update the key-value pair in the YAML data
func updateKeyValue(data interface{}, targetKey, newValue string) {
	switch data := data.(type) {
	case map[string]interface{}:
		for k, v := range data {
			if k == targetKey {
				data[k] = newValue
			} else {
				updateKeyValue(v, targetKey, newValue)
			}
		}
	case []interface{}:
		for _, item := range data {
			updateKeyValue(item, targetKey, newValue)
		}
	}
}
