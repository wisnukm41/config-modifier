import argparse
import os
import yaml

# Argument parsing setup
arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('-k', '--key')
arg_parser.add_argument('-v', '--value')
parsed_args = arg_parser.parse_args()
input_key = parsed_args.key
input_value = parsed_args.value

# Get the current working directory
current_directory = os.getcwd()
yaml_file_paths = []

# Walk through the directory to find YAML files
for dirpath, dirnames, filenames in os.walk(current_directory):
    for filename in filenames:
        if filename.endswith('.yaml') or filename.endswith('.yml'):
            yaml_file_paths.append(os.path.join(dirpath, filename))

# Dictionary to store the contents of the YAML files
yaml_data = {}

# Load the contents of the YAML files
for filepath in yaml_file_paths:
    with open(filepath, 'r') as file:
        try:
            data = yaml.safe_load(file)
            yaml_data[filepath] = data
        except:
            print('Error loading YAML file')

# Function to replace the key-value pair in the YAML content
def update_key_value(filepath, data, key, value):
    if key is None or value is None:
        print('Please provide both key and value to update')
        return

    updated_data = modify_value(data, key, value)
    with open(filepath, 'w') as file:
        yaml.safe_dump(updated_data, file)

# Recursive function to update the key-value pair
def modify_value(data_structure, target_key, new_value):
    if isinstance(data_structure, dict):
        for k, v in data_structure.items():
            if k == target_key:
                data_structure[k] = new_value
            elif isinstance(v, (dict, list)):
                modify_value(v, target_key, new_value)
    elif isinstance(data_structure, list):
        for item in data_structure:
            modify_value(item, target_key, new_value)
    return data_structure

# Try to update the key-value pair in each YAML file
try:
    for filepath, data in yaml_data.items():
        update_key_value(filepath, data, input_key, input_value)
        print(f'Successfully updated key `{input_key}` to value `{input_value}` in {filepath}')
except:
    print('Failed to update files')