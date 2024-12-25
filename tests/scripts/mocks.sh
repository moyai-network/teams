if ! [ -x "$(command -v mockgen)" ]; then
  go install go.uber.org/mock/mockgen@latest
fi

mockgen -destination tests/mocks/ports.go -package mocks -source ./internal/ports/ports.go