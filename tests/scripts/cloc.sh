# shellcheck disable=SC2046
cloc --exclude-dir=$(tr '\n' ',' < .clocignore) .