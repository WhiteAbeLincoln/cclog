# Task Completion Checklist

## When a Development Task is Completed

### Code Quality Checks (MUST RUN)
1. **Format Code**: `make fmt` or `go fmt ./...`
2. **Run Linter**: `make vet` or `go vet ./...`
3. **Run All Tests**: `make test` or `go test ./...`
4. **Test Coverage**: `make test-coverage` or `go test -cover ./...`

### Build Verification
1. **Build Application**: `make build` or `go build -o cclog ./cmd/cclog/`
2. **Verify Binary Works**: `./cclog --help`

### TDD Workflow (CRITICAL)
Following t-wada's TDD practices:
1. **Red**: Write failing test first
2. **Green**: Write minimal code to make test pass
3. **Refactor**: Improve code while keeping tests green
4. **Repeat**: Continue until test list is empty

### Testing Requirements
- All new code MUST have corresponding tests
- Test files follow `*_test.go` naming convention
- TUI components tested with teatest framework
- Use `testdata/sample.jsonl` for realistic test scenarios

### Integration Testing
- Test both CLI and TUI modes
- Verify JSONL parsing with various file structures
- Test message filtering and Markdown output
- Verify clipboard and editor integration

### Final Verification
- Ensure all tests pass: `go test ./...`
- No lint errors: `go vet ./...`
- Code properly formatted: `go fmt ./...`
- Binary builds successfully and runs