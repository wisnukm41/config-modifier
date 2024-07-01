package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v3"
)

func main() {
	// Parse command line arguments
	key := flag.String("k", "", "Key to update (e.g., database.host)")
	environment := flag.String("e", "stag", "Environment default(stag)")
	flag.Parse()

	if *key == "" {
		fmt.Println("Please provide key to modify")
		return
	}

	var sourceFile string
	if *environment == "prod" {
		sourceFile = "./config-example/dbs/db_prod.yaml"
	} else {
		sourceFile = "./config-example/dbs/db_stg.yaml"
	}

	targetFile := "./config-example/apps/web.yaml"

	// Update the key-value pair in the target YAML file based on the source file
	if err := updateYAMLFile(sourceFile, targetFile, *key); err != nil {
		fmt.Println("Failed to update target file:", targetFile, err)
	} else {
		fmt.Printf("Successfully updated key `%s` in %s based on %s\n", *key, targetFile, sourceFile)
	}
}

// Update the key-value pair in the target YAML file based on the source file
func updateYAMLFile(sourceFile, targetFile, key string) error {
	sourceData, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	var sourceNode yaml.Node
	if err := yaml.Unmarshal(sourceData, &sourceNode); err != nil {
		return err
	}

	targetData, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return err
	}

	var targetNode yaml.Node
	if err := yaml.Unmarshal(targetData, &targetNode); err != nil {
		return err
	}

	// Get the value from the source file
	sourceValue := getNestedKeyValue(&sourceNode, key)
	if sourceValue == nil {
		return fmt.Errorf("key `%s` not found in source file", key)
	}

	// Update the key-value pair in the target YAML node
	updateNestedKeyValue(&targetNode, key, sourceValue)

	updatedData, err := yaml.Marshal(&targetNode)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(targetFile, updatedData, 0644)
}

// Recursively get the value of the nested key in the YAML node
func getNestedKeyValue(node *yaml.Node, keyPath string) *yaml.Node {
	keys := strings.Split(keyPath, ".")

	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			value := getNestedKeyValue(child, keyPath)
			if value != nil {
				return value
			}
		}
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if k.Value == keys[0] {
				if len(keys) == 1 {
					return v
				}
				return getNestedKeyValue(v, strings.Join(keys[1:], "."))
			}
		}
	}
	if node.Kind == yaml.SequenceNode {
		for _, child := range node.Content {
			value := getNestedKeyValue(child, keyPath)
			if value != nil {
				return value
			}
		}
	}

	return nil
}

// Recursively update the nested key-value pair in the YAML node
func updateNestedKeyValue(node *yaml.Node, keyPath string, newValue *yaml.Node) {
	keys := strings.Split(keyPath, ".")

	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			updateNestedKeyValue(child, keyPath, newValue)
		}
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if k.Value == keys[0] {
				if len(keys) == 1 {
					*v = *newValue
				} else {
					updateNestedKeyValue(v, strings.Join(keys[1:], "."), newValue)
				}
				return
			}
		}
	}
	if node.Kind == yaml.SequenceNode {
		for _, child := range node.Content {
			updateNestedKeyValue(child, keyPath, newValue)
		}
	}
}
