import json

def order_translations(json_file):
    with open(json_file, 'r') as f:
        data = json.load(f)
    
    if 'translations' in data:
        translations = data['translations']
        ordered_translations = dict(sorted(translations.items()))
        data['translations'] = ordered_translations
        with open(json_file, 'w') as f:
            json.dump(data, f, indent=4)
        print("Translations ordered alphabetically and saved successfully.")
    else:
        print("Error: 'translations' key not found in the JSON file.")

# Usage
json_file = 'assets/translations/fr.json'  # Change this to your JSON file path
order_translations(json_file)