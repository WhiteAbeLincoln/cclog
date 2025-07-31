# TDD Development Guidelines

## TDD Definition (t-wada's approach)
1. **Test List**: Write comprehensive test scenarios list
2. **Select One**: Choose ONE scenario from the list
3. **Red**: Write concrete, executable test code and confirm it fails
4. **Green**: Change production code to make the test (and all previous tests) pass
5. **Refactor**: Improve implementation design while keeping tests green
6. **Repeat**: Return to step 2 until test list is empty

## TDD Flow in cclog Project
1. **Always write tests before implementation**
2. **Run tests frequently**: `go test ./...`
3. **Ensure all tests pass before committing**
4. **Refactor only when tests are green**

## Testing Patterns in cclog
- **Unit Tests**: For individual functions and methods
- **Integration Tests**: For component interactions
- **TUI Tests**: Using teatest framework for UI components
- **Table-Driven Tests**: For multiple test cases
- **Test Data**: Use `testdata/sample.jsonl` for realistic scenarios

## Test Coverage Areas
- Message unmarshaling and data integrity
- File and directory parsing with error cases
- Markdown formatting with various message types
- Content extraction from different message structures
- TUI interactions and state management
- CLI argument parsing and validation

## Key Testing Commands
- `go test ./...` - Run all tests
- `go test -v ./...` - Verbose output
- `go test -cover ./...` - With coverage
- `go test ./pkg/types/` - Specific package