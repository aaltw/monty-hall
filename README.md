# Monty Hall Terminal Simulator

A beautiful, interactive terminal application that demonstrates the famous Monty Hall probability problem. Built with Go and featuring a modern TUI (Terminal User Interface) with rich colors, animations, and comprehensive statistics tracking.

## 🎯 The Monty Hall Problem

The Monty Hall problem is a probability puzzle based on the TV game show "Let's Make a Deal":

1. **Setup**: Three doors hide prizes - one car (valuable) and two goats (less valuable)
2. **Initial Choice**: You choose one door
3. **Host Reveal**: The host opens one of the remaining doors, always revealing a goat
4. **Final Decision**: You can stick with your original choice or switch to the other unopened door

**The Counter-Intuitive Result**: Switching gives you a 2/3 chance of winning, while staying gives only 1/3!

## ✨ Features

### 🎮 Interactive Gameplay
- Beautiful ASCII art doors with visual feedback
- Smooth navigation with keyboard controls
- Real-time game state visualization
- Animated door opening sequences

### 📊 Comprehensive Statistics
- Win/loss tracking for both strategies (switch vs stay)
- Progress bars showing win rates
- Streak tracking (current and best)
- Statistical convergence visualization
- Persistent data storage across sessions

### 🎨 Modern Terminal UI
- Rich color scheme with 24-bit RGB support
- Responsive design for different terminal sizes
- Consistent styling with lipgloss
- Professional ASCII banner and layouts

### 📚 Educational Content
- Built-in help system explaining the problem
- Mathematical insights and probability theory
- Real-time demonstration of statistical convergence
- Clear visual feedback for learning

## 🚀 Installation

### Prerequisites
- Go 1.21 or later
- Terminal with color support (recommended)

### Build from Source
```bash
git clone https://github.com/westhuis/monty-hall.git
cd monty-hall
make build
```

### Run the Application
```bash
./monty-hall
```

## 🎮 How to Play

### Controls
- **Arrow Keys / hjkl**: Navigate menus and options
- **Enter / Space**: Select options
- **1, 2, 3**: Directly select doors
- **s**: Switch choice (during final decision)
- **h**: Toggle help
- **q**: Quit application
- **r**: Reset statistics

### Game Flow
1. **Main Menu**: Choose to play, view statistics, or get help
2. **Initial Choice**: Select one of three doors
3. **Host Reveal**: Watch as the host opens a door with a goat
4. **Final Decision**: Choose to switch or stay with your original choice
5. **Results**: See the outcome and updated statistics

## 📈 Understanding the Statistics

The application tracks detailed statistics to demonstrate the mathematical principle:

- **Stay Strategy**: Should win ~33.3% of games (1/3 probability)
- **Switch Strategy**: Should win ~66.7% of games (2/3 probability)

As you play more games, you'll see the actual results converge to these theoretical probabilities, proving the counter-intuitive nature of the problem.

## 🏗️ Architecture

The application follows clean architecture principles:

```
├── cmd/monty-hall/     # Application entry point
├── pkg/
│   ├── game/          # Core game logic and rules
│   ├── stats/         # Statistics tracking and persistence
│   └── ui/            # Terminal user interface
└── specs/             # Project specifications
```

### Key Components
- **Game Engine**: Implements Monty Hall rules with proper validation
- **Statistics Tracker**: Comprehensive data collection and analysis
- **UI Controller**: Modern TUI built with Bubble Tea framework
- **Persistence Layer**: JSON-based data storage

## 🧪 Testing

The project includes comprehensive tests:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./pkg/game/
go test -v ./pkg/stats/
go test -v ./pkg/ui/
```

### Test Coverage
- **Game Logic**: 90.8% coverage
- **Statistics**: 83.6% coverage  
- **UI Components**: 100% coverage

## 🛠️ Development

### Build Commands
```bash
make build          # Build the application
make test           # Run all tests
make lint           # Run linter (requires golangci-lint)
make fmt            # Format code
make vet            # Vet code
make check          # Run fmt, vet, and test
make clean          # Clean build artifacts
```

### Code Quality
- Follows Go best practices and idioms
- Comprehensive error handling
- Clean separation of concerns
- Extensive documentation
- Type-safe design with custom types

## 📊 Statistics Features

### Tracked Metrics
- Total games played
- Win/loss ratios by strategy
- Current and best win streaks
- Game duration tracking
- Daily statistics
- Historical game records

### Visual Analytics
- Progress bars for win rates
- Statistical cards with key metrics
- Theoretical vs actual comparison
- Convergence insights

## 🎨 Visual Design

### Color Scheme
- **Primary**: Go blue (#00ADD8)
- **Success**: Green (#00D084)
- **Warning**: Orange (#FFA726)
- **Error**: Red (#FF6B6B)
- **Accent**: Various semantic colors

### Components
- Styled doors with state indicators
- Progress bars for statistics
- Cards for metric display
- Responsive layouts
- Consistent spacing and typography

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📝 License

This project is open source and available under the MIT License.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Inspired by the classic Monty Hall problem from probability theory

## 📚 Educational Value

This application serves as an excellent educational tool for:
- **Probability Theory**: Demonstrates counter-intuitive probability
- **Statistical Analysis**: Shows convergence to theoretical values
- **Go Programming**: Example of clean architecture and TUI development
- **Software Testing**: Comprehensive test coverage examples

---

**Try it yourself and see if you can beat the odds! Remember: when in doubt, always switch! 🚗**