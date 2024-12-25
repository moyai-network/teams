if ! [ -x "$(command -v golangci-lint)" ]; then
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
fi

golangci-lint run || exit 1