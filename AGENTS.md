# Agent Development Guidelines

## Build/Test/Lint Commands
- **Build**: `make build` (creates `monty-hall` binary)
- **Test**: `make test` (all tests) or `go test -v -run TestName ./pkg/package/` (single test)
- **Lint**: `make lint` (requires golangci-lint)
- **Format**: `make fmt` and `make vet`
- **Quality Check**: `make check` (fmt + vet + test)
- **IMPORTANT**: Never run `make run` or `./monty-hall` - interactive apps freeze CLI agents

## Code Style Guidelines
- **Imports**: Group standard library first, then third-party packages
- **Naming**: PascalCase for exported items, camelCase for internal variables
- **Types**: Use custom types extensively; enums with `const` blocks and `iota`
- **Error Handling**: Standard `if err != nil` checks, return errors up call stack
- **Structure**: Follow standard Go project layout (`cmd/`, `pkg/`)

## Unicode String Handling (CRITICAL)
**ALWAYS use the correct string length function for Unicode text alignment:**

- **`len(string)`** → Counts **bytes** (❌ WRONG for Unicode alignment)
- **`len([]rune(string))`** → Counts **characters** (✅ CORRECT for Unicode alignment)  
- **`runewidth.StringWidth(string)`** → Counts **visual width** (🏆 BEST for terminal display)

**For terminal UI alignment, follow this pattern:**
```go
// For visual width calculation (best for terminal display)
textWidth := runewidth.StringWidth(text)

// For character-based operations (good for string manipulation)
textRunes := []rune(text)
charCount := len(textRunes)

// NEVER use len(string) for Unicode text alignment - causes misalignment
```

**This codebase uses `runewidth.StringWidth()` in door components - always follow that pattern for consistent Unicode handling.**

## Project Structure
- `cmd/monty-hall/`: Main application entry point
- `pkg/game/`: Core game logic and types
- `pkg/stats/`: Statistics collection and persistence
- `pkg/ui/`: Bubble Tea terminal interface
- `specs/`: Project specifications and documentation