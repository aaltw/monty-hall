# Testing Strategy

## Testing Philosophy

The Monty Hall application requires comprehensive testing to ensure:
1. **Mathematical Accuracy**: Game logic correctly implements probability theory
2. **User Experience**: UI behaves predictably and intuitively
3. **Data Integrity**: Statistics are accurate and persistent
4. **Cross-Platform Compatibility**: Works reliably across different environments
5. **Performance**: Responsive under various conditions

## Test Pyramid Structure

```
                    ┌─────────────────┐
                    │   E2E Tests     │ ← 5% (Critical user journeys)
                    └─────────────────┘
                ┌─────────────────────────┐
                │  Integration Tests      │ ← 15% (Component interactions)
                └─────────────────────────┘
            ┌─────────────────────────────────┐
            │      Unit Tests                 │ ← 80% (Individual functions)
            └─────────────────────────────────┘
```

## Unit Testing

### Game Logic Tests (`pkg/game/`)

#### Core Game Mechanics
```go
// TestNewGame verifies game initialization
func TestNewGame(t *testing.T) {
    game := NewGame()
    
    assert.NotNil(t, game)
    assert.Equal(t, PhaseInitialized, game.Phase)
    assert.Len(t, game.Doors, 3)
    
    // Verify exactly one prize
    prizeCount := 0
    for _, door := range game.Doors {
        if door.HasPrize {
            prizeCount++
        }
    }
    assert.Equal(t, 1, prizeCount)
}

// TestDoorSelection verifies player choice logic
func TestDoorSelection(t *testing.T) {
    tests := []struct {
        name       string
        doorNumber int
        wantErr    bool
    }{
        {"Valid door 1", 1, false},
        {"Valid door 2", 2, false},
        {"Valid door 3", 3, false},
        {"Invalid door 0", 0, true},
        {"Invalid door 4", 4, true},
        {"Invalid door -1", -1, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            game := NewGame()
            err := game.SelectDoor(tt.doorNumber)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, PhasePlayerChoice, game.Phase)
                assert.Equal(t, tt.doorNumber-1, game.PlayerChoice)
            }
        })
    }
}

// TestHostRevealLogic verifies host behavior
func TestHostRevealLogic(t *testing.T) {
    // Test all possible prize locations and player choices
    for prizeLocation := 0; prizeLocation < 3; prizeLocation++ {
        for playerChoice := 0; playerChoice < 3; playerChoice++ {
            t.Run(fmt.Sprintf("prize_%d_choice_%d", prizeLocation, playerChoice), func(t *testing.T) {
                game := createGameWithSetup(prizeLocation, playerChoice)
                
                revealedDoor, err := game.HostRevealDoor()
                assert.NoError(t, err)
                
                // Host must not reveal player's door
                assert.NotEqual(t, playerChoice+1, revealedDoor)
                
                // Host must not reveal prize door
                assert.NotEqual(t, prizeLocation+1, revealedDoor)
                
                // Revealed door must contain goat
                assert.False(t, game.Doors[revealedDoor-1].HasPrize)
                assert.True(t, game.Doors[revealedDoor-1].IsOpen)
            })
        }
    }
}
```

#### Probability Validation
```go
// TestProbabilityDistribution verifies statistical correctness
func TestProbabilityDistribution(t *testing.T) {
    const numGames = 10000
    const tolerance = 0.05 // 5% tolerance
    
    stickWins := 0
    switchWins := 0
    
    for i := 0; i < numGames; i++ {
        // Test stick strategy
        game := NewGame()
        game.SelectDoor(1) // Always choose door 1
        game.HostRevealDoor()
        result, _ := game.MakeFinalChoice(StrategyStick)
        if result.Won {
            stickWins++
        }
        
        // Test switch strategy
        game = NewGame()
        game.SelectDoor(1) // Always choose door 1
        game.HostRevealDoor()
        result, _ = game.MakeFinalChoice(StrategySwitch)
        if result.Won {
            switchWins++
        }
    }
    
    stickRate := float64(stickWins) / numGames
    switchRate := float64(switchWins) / numGames
    
    // Verify theoretical probabilities
    assert.InDelta(t, 1.0/3.0, stickRate, tolerance, "Stick strategy should win ~33.3%")
    assert.InDelta(t, 2.0/3.0, switchRate, tolerance, "Switch strategy should win ~66.7%")
}

// TestRandomnessUniformity verifies prize placement is uniform
func TestRandomnessUniformity(t *testing.T) {
    const numGames = 10000
    const tolerance = 0.05
    
    doorCounts := make([]int, 3)
    
    for i := 0; i < numGames; i++ {
        game := NewGame()
        doorCounts[game.PrizeLocation]++
    }
    
    expectedCount := numGames / 3
    for i, count := range doorCounts {
        rate := float64(count) / numGames
        assert.InDelta(t, 1.0/3.0, rate, tolerance, 
            "Door %d should have prize ~33.3% of the time", i+1)
    }
}
```

#### Edge Cases and Error Conditions
```go
// TestInvalidPhaseTransitions verifies state machine integrity
func TestInvalidPhaseTransitions(t *testing.T) {
    game := NewGame()
    
    // Cannot reveal door before player choice
    _, err := game.HostRevealDoor()
    assert.Error(t, err)
    
    // Cannot make final choice before host reveal
    _, err = game.MakeFinalChoice(StrategyStick)
    assert.Error(t, err)
    
    // Cannot select door twice
    game.SelectDoor(1)
    err = game.SelectDoor(2)
    assert.Error(t, err)
}

// TestGameValidation verifies game state validation
func TestGameValidation(t *testing.T) {
    tests := []struct {
        name    string
        setup   func() *Game
        wantErr bool
    }{
        {
            name: "Valid game",
            setup: func() *Game {
                return NewGame()
            },
            wantErr: false,
        },
        {
            name: "No prize",
            setup: func() *Game {
                game := NewGame()
                for _, door := range game.Doors {
                    door.HasPrize = false
                }
                return game
            },
            wantErr: true,
        },
        {
            name: "Multiple prizes",
            setup: func() *Game {
                game := NewGame()
                game.Doors[0].HasPrize = true
                game.Doors[1].HasPrize = true
                return game
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            game := tt.setup()
            err := game.Validate()
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Statistics Tests (`pkg/stats/`)

#### Statistics Tracking
```go
// TestStatisticsTracking verifies accurate game recording
func TestStatisticsTracking(t *testing.T) {
    stats := NewStatistics()
    
    // Record some games
    results := []*game.Result{
        {Won: true, Strategy: game.StrategyStick},
        {Won: false, Strategy: game.StrategyStick},
        {Won: true, Strategy: game.StrategySwitch},
        {Won: true, Strategy: game.StrategySwitch},
    }
    
    for _, result := range results {
        stats.RecordGame(result)
    }
    
    assert.Equal(t, 4, stats.TotalGames)
    assert.Equal(t, 1, stats.StickWins)
    assert.Equal(t, 1, stats.StickLosses)
    assert.Equal(t, 2, stats.SwitchWins)
    assert.Equal(t, 0, stats.SwitchLosses)
    
    assert.Equal(t, 0.5, stats.GetStickWinRate())
    assert.Equal(t, 1.0, stats.GetSwitchWinRate())
}

// TestConfidenceIntervals verifies statistical calculations
func TestConfidenceIntervals(t *testing.T) {
    tests := []struct {
        wins       int
        total      int
        confidence float64
        wantLower  float64
        wantUpper  float64
        tolerance  float64
    }{
        {50, 100, 0.95, 0.402, 0.598, 0.01},
        {67, 100, 0.95, 0.578, 0.762, 0.01},
        {33, 100, 0.95, 0.238, 0.422, 0.01},
    }
    
    for _, tt := range tests {
        t.Run(fmt.Sprintf("%d_%d_%.2f", tt.wins, tt.total, tt.confidence), func(t *testing.T) {
            lower, upper := CalculateConfidenceInterval(tt.wins, tt.total, tt.confidence)
            
            assert.InDelta(t, tt.wantLower, lower, tt.tolerance)
            assert.InDelta(t, tt.wantUpper, upper, tt.tolerance)
        })
    }
}
```

#### Data Persistence
```go
// TestStatisticsPersistence verifies save/load functionality
func TestStatisticsPersistence(t *testing.T) {
    tempDir := t.TempDir()
    filePath := filepath.Join(tempDir, "test_stats.json")
    
    pm := NewPersistenceManager(filePath)
    
    // Create and save statistics
    originalStats := NewStatistics()
    originalStats.RecordGame(&game.Result{Won: true, Strategy: game.StrategySwitch})
    originalStats.RecordGame(&game.Result{Won: false, Strategy: game.StrategyStick})
    
    err := pm.SaveStatistics(originalStats)
    assert.NoError(t, err)
    
    // Load statistics
    loadedStats, err := pm.LoadStatistics()
    assert.NoError(t, err)
    
    // Verify data integrity
    assert.Equal(t, originalStats.TotalGames, loadedStats.TotalGames)
    assert.Equal(t, originalStats.SwitchWins, loadedStats.SwitchWins)
    assert.Equal(t, originalStats.StickLosses, loadedStats.StickLosses)
}
```

### UI Component Tests (`pkg/ui/`)

#### Component Rendering
```go
// TestDoorComponentRendering verifies visual representation
func TestDoorComponentRendering(t *testing.T) {
    doors := [3]*game.Door{
        {Number: 1, IsOpen: false, HasPrize: false, State: game.DoorClosed},
        {Number: 2, IsOpen: true, HasPrize: false, State: game.DoorOpenWithGoat},
        {Number: 3, IsOpen: true, HasPrize: true, State: game.DoorOpenWithPrize},
    }
    
    component := NewDoorComponent(doors)
    rendered := component.Render(80)
    
    // Verify door numbers are present
    assert.Contains(t, rendered, "1")
    assert.Contains(t, rendered, "2")
    assert.Contains(t, rendered, "3")
    
    // Verify goat and car representations
    assert.Contains(t, rendered, "🐐")
    assert.Contains(t, rendered, "🚗")
}

// TestResponsiveLayout verifies layout adaptation
func TestResponsiveLayout(t *testing.T) {
    tests := []struct {
        width    int
        height   int
        expected ScreenSize
    }{
        {70, 20, ScreenSmall},
        {80, 24, ScreenMedium},
        {120, 30, ScreenLarge},
    }
    
    for _, tt := range tests {
        t.Run(fmt.Sprintf("%dx%d", tt.width, tt.height), func(t *testing.T) {
            size := DetectScreenSize(tt.width, tt.height)
            assert.Equal(t, tt.expected, size)
        })
    }
}
```

#### Input Handling
```go
// TestKeyboardInput verifies input processing
func TestKeyboardInput(t *testing.T) {
    app := NewApp()
    
    tests := []struct {
        key      tea.KeyMsg
        expected tea.Cmd
    }{
        {tea.KeyMsg{Type: tea.KeyLeft}, nil},
        {tea.KeyMsg{Type: tea.KeyRight}, nil},
        {tea.KeyMsg{Type: tea.KeyEnter}, nil},
        {tea.KeyMsg{Type: tea.KeyEsc}, nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.key.String(), func(t *testing.T) {
            _, cmd := app.Update(tt.key)
            // Verify command is appropriate for key
            // (specific assertions depend on implementation)
            _ = cmd
        })
    }
}
```

## Integration Testing

### Game Flow Integration
```go
// TestCompleteGameFlow verifies end-to-end game functionality
func TestCompleteGameFlow(t *testing.T) {
    app := NewApp()
    
    // Start new game
    app.screen = ScreenGame
    app.game = game.NewGame()
    
    // Player selects door
    _, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
    assert.NotNil(t, cmd)
    
    // Confirm selection
    _, cmd = app.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.NotNil(t, cmd)
    
    // Host reveals door (automatic)
    assert.Equal(t, game.PhaseHostReveal, app.game.Phase)
    
    // Player makes final choice
    _, cmd = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}) // Switch
    assert.NotNil(t, cmd)
    
    // Game should be complete
    assert.Equal(t, game.PhaseComplete, app.game.Phase)
    assert.NotNil(t, app.game.Result)
}

// TestStatisticsIntegration verifies stats tracking with game play
func TestStatisticsIntegration(t *testing.T) {
    app := NewApp()
    initialGames := app.stats.TotalGames
    
    // Play a complete game
    playCompleteGame(app, game.StrategySwitch)
    
    // Verify statistics updated
    assert.Equal(t, initialGames+1, app.stats.TotalGames)
    
    if app.game.Result.Won {
        assert.Equal(t, 1, app.stats.SwitchWins)
    } else {
        assert.Equal(t, 1, app.stats.SwitchLosses)
    }
}
```

### Data Persistence Integration
```go
// TestPersistenceIntegration verifies data survives app restart
func TestPersistenceIntegration(t *testing.T) {
    tempDir := t.TempDir()
    
    // Create app with temporary storage
    app1 := NewAppWithStorage(tempDir)
    
    // Play some games
    for i := 0; i < 5; i++ {
        playCompleteGame(app1, game.StrategySwitch)
    }
    
    originalStats := app1.stats
    
    // Simulate app restart
    app2 := NewAppWithStorage(tempDir)
    
    // Verify statistics persisted
    assert.Equal(t, originalStats.TotalGames, app2.stats.TotalGames)
    assert.Equal(t, originalStats.SwitchWins, app2.stats.SwitchWins)
}
```

## End-to-End Testing

### User Journey Tests
```go
// TestNewUserJourney simulates first-time user experience
func TestNewUserJourney(t *testing.T) {
    app := NewApp()
    
    // User starts at main menu
    assert.Equal(t, ScreenMenu, app.screen)
    
    // User selects tutorial
    app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
    assert.Equal(t, ScreenTutorial, app.screen)
    
    // User completes tutorial
    completeTutorial(app)
    
    // User plays first game
    app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
    assert.Equal(t, ScreenGame, app.screen)
    
    playCompleteGame(app, game.StrategySwitch)
    
    // User views statistics
    app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
    assert.Equal(t, ScreenStats, app.screen)
    
    // Verify first game recorded
    assert.Equal(t, 1, app.stats.TotalGames)
}

// TestEducationalEffectiveness verifies learning outcomes
func TestEducationalEffectiveness(t *testing.T) {
    app := NewApp()
    
    // Simulate user playing many games with different strategies
    for i := 0; i < 50; i++ {
        if i%2 == 0 {
            playCompleteGame(app, game.StrategyStick)
        } else {
            playCompleteGame(app, game.StrategySwitch)
        }
    }
    
    // Verify statistical convergence demonstrates the principle
    stickRate := app.stats.GetStickWinRate()
    switchRate := app.stats.GetSwitchWinRate()
    
    // With enough games, switch should clearly outperform stick
    assert.Greater(t, switchRate, stickRate, 
        "Switch strategy should outperform stick strategy")
    
    // Rates should be reasonably close to theoretical values
    assert.InDelta(t, 1.0/3.0, stickRate, 0.15, "Stick rate should approach 33%")
    assert.InDelta(t, 2.0/3.0, switchRate, 0.15, "Switch rate should approach 67%")
}
```

## Property-Based Testing

### Game Invariants
```go
// TestGameInvariants verifies properties that must always hold
func TestGameInvariants(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // Generate random game sequence
        initialDoor := rapid.IntRange(1, 3).Draw(t, "initial_door")
        strategy := rapid.SampledFrom([]game.Strategy{
            game.StrategyStick, 
            game.StrategySwitch,
        }).Draw(t, "strategy")
        
        // Play game
        g := game.NewGame()
        g.SelectDoor(initialDoor)
        g.HostRevealDoor()
        result, err := g.MakeFinalChoice(strategy)
        
        require.NoError(t, err)
        
        // Verify invariants
        assert.True(t, g.Phase == game.PhaseComplete, "Game must be complete")
        assert.NotNil(t, result, "Result must exist")
        
        // Exactly one door has prize
        prizeCount := 0
        for _, door := range g.Doors {
            if door.HasPrize {
                prizeCount++
            }
        }
        assert.Equal(t, 1, prizeCount, "Exactly one door must have prize")
        
        // Host revealed door cannot have prize
        assert.False(t, g.Doors[g.HostRevealed].HasPrize, 
            "Host revealed door cannot have prize")
        
        // Host revealed door cannot be player's initial choice
        assert.NotEqual(t, g.HostRevealed, g.PlayerChoice,
            "Host cannot reveal player's chosen door")
    })
}

// TestStatisticalConvergence verifies probability convergence
func TestStatisticalConvergence(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        numGames := rapid.IntRange(1000, 5000).Draw(t, "num_games")
        
        stats := stats.NewStatistics()
        
        for i := 0; i < numGames; i++ {
            g := game.NewGame()
            g.SelectDoor(1) // Always choose door 1 for consistency
            g.HostRevealDoor()
            
            strategy := game.StrategyStick
            if i%2 == 1 {
                strategy = game.StrategySwitch
            }
            
            result, _ := g.MakeFinalChoice(strategy)
            stats.RecordGame(result)
        }
        
        // With large sample sizes, rates should converge
        if numGames >= 2000 {
            stickRate := stats.GetStickWinRate()
            switchRate := stats.GetSwitchWinRate()
            
            // Allow for statistical variation
            tolerance := 0.1
            assert.InDelta(t, 1.0/3.0, stickRate, tolerance,
                "Stick rate should converge to 33%")
            assert.InDelta(t, 2.0/3.0, switchRate, tolerance,
                "Switch rate should converge to 67%")
        }
    })
}
```

## Performance Testing

### Benchmark Tests
```go
// BenchmarkGameCreation measures game initialization performance
func BenchmarkGameCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        game.NewGame()
    }
}

// BenchmarkCompleteGameFlow measures full game performance
func BenchmarkCompleteGameFlow(b *testing.B) {
    for i := 0; i < b.N; i++ {
        g := game.NewGame()
        g.SelectDoor(1)
        g.HostRevealDoor()
        g.MakeFinalChoice(game.StrategySwitch)
    }
}

// BenchmarkUIRendering measures rendering performance
func BenchmarkUIRendering(b *testing.B) {
    app := NewApp()
    app.game = game.NewGame()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        app.View()
    }
}

// BenchmarkStatisticsCalculation measures stats performance
func BenchmarkStatisticsCalculation(b *testing.B) {
    stats := createLargeStatistics(10000) // 10k games
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        stats.GetStickWinRate()
        stats.GetSwitchWinRate()
        CalculateConfidenceInterval(stats.SwitchWins, 
            stats.SwitchWins+stats.SwitchLosses, 0.95)
    }
}
```

### Memory Usage Tests
```go
// TestMemoryUsage verifies reasonable memory consumption
func TestMemoryUsage(t *testing.T) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // Create many games
    games := make([]*game.Game, 1000)
    for i := range games {
        games[i] = game.NewGame()
    }
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    // Memory usage should be reasonable
    memoryPerGame := (m2.Alloc - m1.Alloc) / uint64(len(games))
    assert.Less(t, memoryPerGame, uint64(1024), // Less than 1KB per game
        "Memory usage per game should be minimal")
}
```

## Cross-Platform Testing

### Platform-Specific Tests
```go
// TestTerminalCompatibility verifies terminal support
func TestTerminalCompatibility(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Run("Windows Terminal", testWindowsTerminal)
        t.Run("Command Prompt", testCommandPrompt)
        t.Run("PowerShell", testPowerShell)
    } else if runtime.GOOS == "darwin" {
        t.Run("Terminal.app", testMacTerminal)
        t.Run("iTerm2", testITerm2)
    } else if runtime.GOOS == "linux" {
        t.Run("GNOME Terminal", testGnomeTerminal)
        t.Run("Konsole", testKonsole)
        t.Run("xterm", testXterm)
    }
}

// TestColorSupport verifies color handling
func TestColorSupport(t *testing.T) {
    tests := []struct {
        colorTerm string
        expected  bool
    }{
        {"truecolor", true},
        {"256color", true},
        {"16color", true},
        {"", false},
        {"monochrome", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.colorTerm, func(t *testing.T) {
            os.Setenv("COLORTERM", tt.colorTerm)
            defer os.Unsetenv("COLORTERM")
            
            hasColor := detectColorSupport()
            assert.Equal(t, tt.expected, hasColor)
        })
    }
}
```

## Test Data and Fixtures

### Test Helpers
```go
// createGameWithSetup creates a game with specific configuration
func createGameWithSetup(prizeLocation, playerChoice int) *game.Game {
    g := &game.Game{
        Phase:         game.PhasePlayerChoice,
        PlayerChoice:  playerChoice,
        PrizeLocation: prizeLocation,
        Created:       time.Now(),
    }
    
    // Initialize doors
    for i := 0; i < 3; i++ {
        g.Doors[i] = &game.Door{
            Number:   i + 1,
            Index:    i,
            HasPrize: i == prizeLocation,
            State:    game.DoorClosed,
        }
    }
    
    if playerChoice >= 0 && playerChoice < 3 {
        g.Doors[playerChoice].State = game.DoorSelected
    }
    
    return g
}

// playCompleteGame simulates a full game with specified strategy
func playCompleteGame(app *App, strategy game.Strategy) {
    app.game = game.NewGame()
    app.game.SelectDoor(1)
    app.game.HostRevealDoor()
    result, _ := app.game.MakeFinalChoice(strategy)
    app.stats.RecordGame(result)
}

// createLargeStatistics creates statistics with many games
func createLargeStatistics(numGames int) *stats.Statistics {
    s := stats.NewStatistics()
    
    for i := 0; i < numGames; i++ {
        strategy := game.StrategyStick
        if i%3 < 2 { // 2/3 switch, 1/3 stick
            strategy = game.StrategySwitch
        }
        
        won := false
        if strategy == game.StrategySwitch && rand.Float64() < 2.0/3.0 {
            won = true
        } else if strategy == game.StrategyStick && rand.Float64() < 1.0/3.0 {
            won = true
        }
        
        result := &game.Result{
            Won:      won,
            Strategy: strategy,
        }
        s.RecordGame(result)
    }
    
    return s
}
```

## Test Coverage Requirements

### Coverage Targets
- **Overall Coverage**: ≥ 90%
- **Game Logic**: ≥ 95%
- **Statistics**: ≥ 95%
- **UI Components**: ≥ 85%
- **Integration**: ≥ 80%

### Coverage Exclusions
- Generated code
- Third-party dependencies
- Platform-specific code that cannot be tested in CI
- Debug/development-only code

## Continuous Integration

### CI Pipeline Tests
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: [1.21, 1.22, 1.23]
    
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Run benchmarks
      run: go test -bench=. -benchmem ./...
    
    - name: Check coverage
      run: |
        go tool cover -func=coverage.out
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Quality Gates
- All tests must pass
- Coverage must meet minimum thresholds
- No race conditions detected
- Benchmarks must not regress significantly
- Linting must pass without errors

## Test Maintenance

### Test Organization
- Group related tests in the same file
- Use descriptive test names
- Include setup and teardown helpers
- Document complex test scenarios

### Test Data Management
- Use table-driven tests for multiple scenarios
- Generate test data programmatically when possible
- Keep test data minimal and focused
- Use temporary directories for file operations

### Flaky Test Prevention
- Avoid time-dependent tests
- Use deterministic random seeds for reproducibility
- Mock external dependencies
- Set appropriate timeouts
- Clean up resources properly