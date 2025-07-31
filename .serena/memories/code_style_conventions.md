# Code Style and Conventions

## Go Style Guide
- Follows standard Go formatting (`go fmt`)
- Standard Go naming conventions (PascalCase for exported, camelCase for unexported)
- Proper package organization with clear separation of concerns

## Code Structure Patterns
- **Layered Architecture**: Clear separation between CLI, parsing, formatting, and TUI layers
- **Standard Project Layout**: Follows Go's conventional project structure
- **Package Boundaries**: `internal/` for private code, `pkg/` for reusable packages
- **Error Handling**: Proper Go error handling patterns throughout

## Testing Conventions
- Test files named `*_test.go` following Go conventions
- Comprehensive test coverage for all packages
- Table-driven tests where appropriate
- TUI testing using Bubble Tea's teatest framework
- Test data in `testdata/` directory for realistic scenarios

## Documentation Style
- Package-level documentation for each package
- Function documentation for exported functions
- README.md with comprehensive usage examples
- CLAUDE.md with development guidelines and architecture overview

## Import Organization
- Standard library imports first
- Third-party imports second
- Local package imports last
- Proper grouping and ordering