# Essential Development Commands

## Build Commands
```bash
go build -o cclog ./cmd/cclog/    # Build application
make build                        # Alternative build via Makefile
make install                      # Install to GOPATH/bin
make build-all                    # Build for multiple platforms
```

## Testing Commands
```bash
go test ./...                     # Run all tests
go test -cover ./...              # Run tests with coverage
go test -v ./...                  # Verbose test output
go test ./pkg/types/              # Test specific package
make test                         # Run all tests via Makefile
make test-coverage                # Run tests with coverage via Makefile
```

## Code Quality Commands
```bash
make fmt                          # Format code (go fmt ./...)
make vet                          # Run linter (go vet ./...)
go mod tidy                       # Clean up dependencies
make deps                         # Install and tidy dependencies
```

## Development Commands
```bash
./cclog                           # Run built application
make run                          # Build and run (starts TUI mode)
make clean                        # Clean build artifacts
```

## Platform-Specific (macOS)
- Uses standard Unix commands (`ls`, `cd`, `grep`, `find`)
- Editor integration via `$EDITOR` environment variable
- Clipboard integration works natively on macOS