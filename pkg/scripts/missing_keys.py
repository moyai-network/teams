import json

def read_json(file_path):
    """Read JSON file and return the data."""
    with open(file_path, 'r') as file:
        return json.load(file)

def find_missing_keys(json1, json2):
    """Find keys that are in json1 but missing in json2."""
    keys1 = set(json1['translations'])
    keys2 = set(json2['translations'])
    missing_keys = keys1 - keys2
    return missing_keys

def main():
    # Paths to your JSON files
    json_file1 = 'assets/translations/en.json'
    json_file2 = 'assets/translations/es.json'

    # Read the JSON files
    json_data1 = read_json(json_file1)
    json_data2 = read_json(json_file2)

    # Find missing keys
    missing_keys = find_missing_keys(json_data1, json_data2)

    # Convert the set to a string
    missing_keys_str = '\n'.join(sorted(missing_keys))

    # Print the missing keys
    print("Keys missing from the second JSON file:", missing_keys_str)

if __name__ == "__main__":
    main()