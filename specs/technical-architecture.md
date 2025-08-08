# Technical Architecture

## System Architecture Overview

The Monty Hall application follows a clean architecture pattern with clear separation between game logic, user interface, and data persistence layers.

```
┌─────────────────────────────────────────────────────────────┐
│                    Terminal Interface                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Game Screen   │  │  Stats Screen   │  │ Tutorial UI  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                     UI Controller                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  Input Handler  │  │  State Manager  │  │ View Router  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Game Engine   │  │ Stats Calculator│  │ Probability  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Data Layer                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Game State    │  │   Statistics    │  │ Config Data  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

```
monty-hall/
├── cmd/
│   └── monty-hall/
│       └── main.go                 # Application entry point
├── pkg/
│   ├── game/                       # Core game logic
│   │   ├── game.go                 # Game state and rules
│   │   ├── door.go                 # Door representation
│   │   ├── host.go                 # Host behavior logic
│   │   └── probability.go          # Probability calculations
│   ├── stats/                      # Statistics tracking
│   │   ├── tracker.go              # Statistics collection
│   │   ├── calculator.go           # Statistical calculations
│   │   └── persistence.go          # Save/load statistics
│   ├── ui/                         # Terminal user interface
│   │   ├── app.go                  # Main application controller
│   │   ├── screens/                # Individual screen components
│   │   │   ├── game.go             # Game play screen
│   │   │   ├── stats.go            # Statistics display
│   │   │   ├── tutorial.go         # Tutorial screens
│   │   │   └── menu.go             # Main menu
│   │   ├── components/             # Reusable UI components
│   │   │   ├── doors.go            # Door visualization
│   │   │   ├── progress.go         # Progress indicators
│   │   │   └── help.go             # Help text components
│   │   └── styles/                 # Visual styling
│   │       ├── colors.go           # Color definitions
│   │       └── layout.go           # Layout constants
│   └── config/                     # Configuration management
│       ├── config.go               # Configuration structure
│       └── defaults.go             # Default values
├── internal/                       # Private application code
│   ├── models/                     # Shared data structures
│   │   ├── game_state.go           # Game state models
│   │   └── ui_state.go             # UI state models
│   └── utils/                      # Utility functions
│       ├── random.go               # Random number generation
│       └── validation.go           # Input validation
├── testdata/                       # Test fixtures and data
└── docs/                          # Documentation
```

## Core Components

### Game Engine (`pkg/game/`)

**Purpose**: Implements the core Monty Hall game logic and rules.

**Key Types**:
```go
type Game struct {
    doors         [3]*Door
    playerChoice  int
    hostRevealed  int
    gamePhase     Phase
    prizeLocation int
}

type Door struct {
    number   int
    isOpen   bool
    hasPrize bool
    state    DoorState
}

type Phase int
const (
    PhaseInitialChoice Phase = iota
    PhaseHostReveal
    PhaseFinalDecision
    PhaseGameOver
)
```

**Responsibilities**:
- Initialize new games with random prize placement
- Track player choices and game state
- Implement host behavior (reveal goat door)
- Determine game outcomes
- Validate moves and state transitions

### Statistics Tracker (`pkg/stats/`)

**Purpose**: Collects and analyzes game statistics over time.

**Key Types**:
```go
type Statistics struct {
    TotalGames    int
    StickWins     int
    StickLosses   int
    SwitchWins    int
    SwitchLosses  int
    SessionStats  SessionStats
}

type SessionStats struct {
    GamesPlayed   int
    StartTime     time.Time
    StickStrategy int
    SwitchStrategy int
}
```

**Responsibilities**:
- Track wins/losses for each strategy
- Calculate win percentages and confidence intervals
- Persist statistics to local storage
- Provide statistical analysis and insights

### UI Controller (`pkg/ui/`)

**Purpose**: Manages the terminal user interface using Bubbletea framework.

**Key Types**:
```go
type App struct {
    game     *game.Game
    stats    *stats.Statistics
    screen   Screen
    state    UIState
}

type Screen int
const (
    ScreenMenu Screen = iota
    ScreenGame
    ScreenStats
    ScreenTutorial
    ScreenHelp
)
```

**Responsibilities**:
- Handle keyboard input and navigation
- Render game state visually
- Manage screen transitions
- Coordinate between game logic and display

## Data Flow

### Game Play Flow
1. **Initialization**: User starts new game
2. **Door Selection**: Player chooses initial door
3. **Host Reveal**: System reveals goat door
4. **Final Decision**: Player chooses to stick or switch
5. **Result**: Game outcome determined and displayed
6. **Statistics Update**: Results recorded for analysis

### State Management
- **Game State**: Managed by game engine, immutable transitions
- **UI State**: Managed by UI controller, reactive to user input
- **Statistics State**: Persistent across sessions, updated after each game

## Technology Stack

### Core Dependencies
- **Go 1.23+**: Latest Go features including range-over-func iterators
- **Bubbletea**: Modern TUI framework for terminal applications with full event handling
- **Lipgloss**: Advanced styling with 24-bit RGB color support, gradients, and layouts
- **Bubbles**: Pre-built components (progress bars, text inputs, viewports)
- **Harmonica**: Animation and easing functions for smooth transitions
- **Cobra** (optional): CLI command structure if needed

### Development Tools
- **golangci-lint**: Code linting and quality checks
- **goimports**: Import organization and formatting
- **testify**: Enhanced testing capabilities
- **go mod**: Dependency management

### File Formats
- **JSON**: Statistics persistence and configuration
- **Markdown**: Documentation and help content

## Design Patterns

### Model-View-Controller (MVC)
- **Model**: Game logic and statistics (`pkg/game/`, `pkg/stats/`)
- **View**: Terminal UI components (`pkg/ui/screens/`, `pkg/ui/components/`)
- **Controller**: UI application controller (`pkg/ui/app.go`)

### Observer Pattern
- UI components observe game state changes
- Statistics tracker observes game completion events
- Reactive updates without tight coupling

### Strategy Pattern
- Different rendering strategies for various terminal capabilities
- Pluggable statistics calculation methods
- Configurable game rule variations

## Error Handling

### Error Categories
1. **User Input Errors**: Invalid selections, out-of-range values
2. **System Errors**: File I/O failures, terminal capability issues
3. **Logic Errors**: Invalid game state transitions

### Error Handling Strategy
- **Graceful Degradation**: Continue operation with reduced functionality
- **User Feedback**: Clear error messages and recovery suggestions
- **Logging**: Detailed error logging for debugging
- **Recovery**: Automatic recovery from transient errors

## Performance Considerations

### Memory Management
- Minimal memory allocation during gameplay
- Efficient string building for UI rendering
- Garbage collection friendly data structures

### Rendering Optimization
- Differential rendering (only update changed areas)
- Efficient terminal escape sequence usage
- Responsive input handling

### Startup Performance
- Lazy loading of non-essential components
- Fast configuration loading
- Minimal initialization overhead

## Security Considerations

### Data Privacy
- No network communication
- Local file storage only
- No sensitive data collection

### Input Validation
- Sanitize all user input
- Validate game state transitions
- Prevent invalid operations

## Testing Strategy

### Unit Testing
- Comprehensive game logic testing
- Statistics calculation validation
- UI component testing

### Integration Testing
- End-to-end game flow testing
- File persistence testing
- Cross-platform compatibility testing

### Property-Based Testing
- Statistical convergence testing
- Random game sequence validation
- Probability distribution verification

## Deployment and Distribution

### Build Process
- Cross-compilation for multiple platforms
- Static binary generation
- Automated build pipeline

### Distribution
- GitHub releases with binaries
- Package manager integration (Homebrew, etc.)
- Docker container option

### Installation
- Single binary deployment
- No external dependencies
- Automatic configuration initialization