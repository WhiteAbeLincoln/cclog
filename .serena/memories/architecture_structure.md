# Codebase Architecture and Structure

## Project Layout (Go Standard)
```
├── cmd/cclog/           # Main application entry point
├── internal/            # Private application packages
│   ├── cli/            # Command-line interface, argument parsing, TUI entry
│   ├── parser/         # JSONL file parsing logic
│   └── formatter/      # Message filtering and Markdown conversion
├── pkg/                # Public packages (reusable)
│   ├── types/          # Core data structures (Message, ConversationLog)
│   └── filepicker/     # Interactive TUI implementation
├── testdata/           # Test data files (sample.jsonl)
├── specs/              # Specifications or documentation
└── .devcontainer/      # Development container configuration
```

## Core Data Flow
1. **JSONL Parsing** (`internal/parser`) - Reads conversation log files
2. **Type System** (`pkg/types`) - Message structures and conversation logs
3. **Message Filtering** (`internal/formatter/filter`) - Removes noise and system messages
4. **Markdown Formatting** (`internal/formatter/markdown`) - Converts to readable Markdown
5. **CLI Interface** (`internal/cli`) - Command-line interface and TUI orchestration
6. **TUI System** (`pkg/filepicker`) - Interactive file browser with live preview

## Key Components
- **Message Type System**: Handles complex JSONL structure from Claude conversations
- **Parser Strategy**: Line-by-line JSONL parsing with buffer expansion (up to 1MB)
- **Message Filtering**: Intelligent filtering of system messages, API errors, interrupted requests
- **Markdown Generation**: Time-sorted processing with timezone conversion
- **TUI Architecture**: File browser, live preview, conversation metadata, clipboard integration