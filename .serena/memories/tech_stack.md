# Tech Stack and Dependencies

## Language
- **Go 1.24** - Modern Go with latest features

## Core Dependencies
- **Standard Library Only** for parsing and formatting logic
- **github.com/charmbracelet/bubbletea** - TUI framework (v1.3.5)
- **github.com/charmbracelet/lipgloss** - TUI styling (v1.1.0)
- **github.com/atotto/clipboard** - Cross-platform clipboard support (v0.1.4)
- **golang.org/x/term** - Terminal handling (v0.32.0)

## Testing Dependencies
- **github.com/charmbracelet/x/exp/teatest** - TUI testing framework
- **github.com/philistino/teacup** - Additional TUI testing utilities

## Key Design Decisions
- Uses only Go standard library for core parsing/formatting logic
- TUI built with Charm's Bubble Tea ecosystem
- Cross-platform clipboard integration with atotto/clipboard
- Comprehensive test coverage including TUI interaction tests