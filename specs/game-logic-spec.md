# Game Logic Specification

## Core Game Mechanics

### Game State Model

The Monty Hall game follows a strict state machine with well-defined transitions and validation rules.

```go
type Game struct {
    ID            string    `json:"id"`
    Doors         [3]*Door  `json:"doors"`
    PlayerChoice  int       `json:"player_choice"`   // 0, 1, or 2 (internal), 1, 2, or 3 (user-facing)
    HostRevealed  int       `json:"host_revealed"`   // Door index revealed by host
    FinalChoice   int       `json:"final_choice"`    // Player's final door choice
    PrizeLocation int       `json:"prize_location"`  // Door index with the prize
    Phase         Phase     `json:"phase"`           // Current game phase
    Strategy      Strategy  `json:"strategy"`        // Stick or Switch
    Result        *Result   `json:"result"`          // Game outcome (nil until complete)
    Created       time.Time `json:"created"`
    Completed     time.Time `json:"completed"`
}

type Door struct {
    Number   int       `json:"number"`    // 1, 2, or 3 (user-facing)
    Index    int       `json:"index"`     // 0, 1, or 2 (internal)
    IsOpen   bool      `json:"is_open"`   // Whether door has been opened
    HasPrize bool      `json:"has_prize"` // Whether door contains the prize
    State    DoorState `json:"state"`     // Current visual state
}

type Phase int
const (
    PhaseInitialized Phase = iota  // Game created, ready for player choice
    PhasePlayerChoice               // Player has selected initial door
    PhaseHostReveal                 // Host has revealed a goat door
    PhaseFinalDecision             // Player has made final choice (stick/switch)
    PhaseComplete                  // Game finished, result determined
)

type DoorState int
const (
    DoorClosed DoorState = iota    // Door is closed
    DoorOpenWithPrize              // Door is open, showing prize
    DoorOpenWithGoat               // Door is open, showing goat
    DoorSelected                   // Door is selected but not opened
)

type Strategy int
const (
    StrategyStick Strategy = iota   // Player sticks with original choice
    StrategySwitch                  // Player switches to other door
)

type Result struct {
    Won           bool      `json:"won"`
    Strategy      Strategy  `json:"strategy"`
    InitialChoice int       `json:"initial_choice"`
    FinalChoice   int       `json:"final_choice"`
    PrizeLocation int       `json:"prize_location"`
    GameDuration  time.Duration `json:"game_duration"`
}
```

### Game Rules and Constraints

#### Rule 1: Prize Placement
- Exactly one door contains the prize (car)
- Exactly two doors contain goats
- Prize placement is random and uniform across all doors
- Prize location is determined at game initialization

#### Rule 2: Player Initial Choice
- Player must select exactly one door initially
- Choice must be valid (door 1, 2, or 3)
- Player cannot change this choice until the final decision phase
- Initial choice does not open the door

#### Rule 3: Host Behavior
- Host must reveal exactly one door after player's initial choice
- Host must reveal a door that contains a goat (never the prize)
- Host cannot reveal the door chosen by the player
- Host's choice is deterministic given the constraints above

#### Rule 4: Final Decision
- Player can choose to stick with original choice or switch
- If switching, player must choose the remaining unopened door
- Final decision immediately determines the game outcome

#### Rule 5: Game Outcome
- Player wins if their final choice contains the prize
- Player loses if their final choice contains a goat
- Game ends immediately after final choice is made

### State Transitions

```
PhaseInitialized
    ↓ (player selects door)
PhasePlayerChoice
    ↓ (host reveals goat door)
PhaseHostReveal
    ↓ (player makes final decision)
PhaseFinalDecision
    ↓ (result calculated)
PhaseComplete
```

### Game Logic Implementation

#### Game Initialization
```go
func NewGame() *Game {
    game := &Game{
        ID:      generateGameID(),
        Phase:   PhaseInitialized,
        Created: time.Now(),
    }
    
    // Initialize doors
    for i := 0; i < 3; i++ {
        game.Doors[i] = &Door{
            Number: i + 1,
            Index:  i,
            IsOpen: false,
            State:  DoorClosed,
        }
    }
    
    // Randomly place prize
    game.PrizeLocation = rand.Intn(3)
    game.Doors[game.PrizeLocation].HasPrize = true
    
    return game
}
```

#### Player Door Selection
```go
func (g *Game) SelectDoor(doorNumber int) error {
    // Validate game state
    if g.Phase != PhaseInitialized {
        return ErrInvalidPhase
    }
    
    // Validate door number
    if doorNumber < 1 || doorNumber > 3 {
        return ErrInvalidDoor
    }
    
    // Convert to internal index
    doorIndex := doorNumber - 1
    
    // Set player choice
    g.PlayerChoice = doorIndex
    g.Doors[doorIndex].State = DoorSelected
    g.Phase = PhasePlayerChoice
    
    return nil
}
```

#### Host Door Revelation
```go
func (g *Game) HostRevealDoor() (int, error) {
    // Validate game state
    if g.Phase != PhasePlayerChoice {
        return 0, ErrInvalidPhase
    }
    
    // Find door to reveal (must be goat, not player's choice)
    revealIndex := g.findDoorToReveal()
    if revealIndex == -1 {
        return 0, ErrNoValidDoorToReveal
    }
    
    // Reveal the door
    g.HostRevealed = revealIndex
    g.Doors[revealIndex].IsOpen = true
    g.Doors[revealIndex].State = DoorOpenWithGoat
    g.Phase = PhaseHostReveal
    
    return revealIndex + 1, nil // Convert to user-facing number
}

func (g *Game) findDoorToReveal() int {
    for i := 0; i < 3; i++ {
        // Skip player's chosen door
        if i == g.PlayerChoice {
            continue
        }
        
        // Skip door with prize
        if g.Doors[i].HasPrize {
            continue
        }
        
        // This door has a goat and wasn't chosen by player
        return i
    }
    return -1 // Should never happen in valid game
}
```

#### Final Decision Processing
```go
func (g *Game) MakeFinalChoice(strategy Strategy) (*Result, error) {
    // Validate game state
    if g.Phase != PhaseHostReveal {
        return nil, ErrInvalidPhase
    }
    
    // Determine final door based on strategy
    var finalDoor int
    switch strategy {
    case StrategyStick:
        finalDoor = g.PlayerChoice
    case StrategySwitch:
        finalDoor = g.findSwitchDoor()
    default:
        return nil, ErrInvalidStrategy
    }
    
    // Set final choice and strategy
    g.FinalChoice = finalDoor
    g.Strategy = strategy
    
    // Calculate result
    won := g.Doors[finalDoor].HasPrize
    
    // Create result
    result := &Result{
        Won:           won,
        Strategy:      strategy,
        InitialChoice: g.PlayerChoice + 1, // Convert to user-facing
        FinalChoice:   finalDoor + 1,      // Convert to user-facing
        PrizeLocation: g.PrizeLocation + 1, // Convert to user-facing
        GameDuration:  time.Since(g.Created),
    }
    
    // Update game state
    g.Result = result
    g.Phase = PhaseComplete
    g.Completed = time.Now()
    
    // Open all doors for final reveal
    g.revealAllDoors()
    
    return result, nil
}

func (g *Game) findSwitchDoor() int {
    for i := 0; i < 3; i++ {
        // Skip player's original choice
        if i == g.PlayerChoice {
            continue
        }
        
        // Skip door revealed by host
        if i == g.HostRevealed {
            continue
        }
        
        // This is the switch door
        return i
    }
    return -1 // Should never happen in valid game
}

func (g *Game) revealAllDoors() {
    for i := 0; i < 3; i++ {
        g.Doors[i].IsOpen = true
        if g.Doors[i].HasPrize {
            g.Doors[i].State = DoorOpenWithPrize
        } else {
            g.Doors[i].State = DoorOpenWithGoat
        }
    }
}
```

### Probability Mathematics

#### Theoretical Probabilities
- **Initial Choice Correct**: P = 1/3 ≈ 0.3333
- **Initial Choice Incorrect**: P = 2/3 ≈ 0.6667
- **Stick Strategy Wins**: P = 1/3 ≈ 0.3333
- **Switch Strategy Wins**: P = 2/3 ≈ 0.6667

#### Mathematical Proof

**Proof by Cases:**

Let's say the player initially chooses Door 1. There are three equally likely scenarios:

1. **Prize behind Door 1** (probability 1/3):
   - Host reveals Door 2 or 3 (both have goats)
   - Sticking wins, switching loses

2. **Prize behind Door 2** (probability 1/3):
   - Host must reveal Door 3 (has goat)
   - Sticking loses, switching wins (Door 2)

3. **Prize behind Door 3** (probability 1/3):
   - Host must reveal Door 2 (has goat)
   - Sticking loses, switching wins (Door 3)

**Result:**
- Sticking wins in 1 out of 3 cases: P(stick wins) = 1/3
- Switching wins in 2 out of 3 cases: P(switch wins) = 2/3

#### Probability Calculation Functions
```go
func CalculateTheoreticalProbabilities() map[Strategy]float64 {
    return map[Strategy]float64{
        StrategyStick:   1.0 / 3.0,
        StrategySwitch: 2.0 / 3.0,
    }
}

func CalculateObservedProbabilities(stats *Statistics) map[Strategy]float64 {
    result := make(map[Strategy]float64)
    
    stickTotal := stats.StickWins + stats.StickLosses
    if stickTotal > 0 {
        result[StrategyStick] = float64(stats.StickWins) / float64(stickTotal)
    }
    
    switchTotal := stats.SwitchWins + stats.SwitchLosses
    if switchTotal > 0 {
        result[StrategySwitch] = float64(stats.SwitchWins) / float64(switchTotal)
    }
    
    return result
}

func CalculateConfidenceInterval(wins, total int, confidence float64) (float64, float64) {
    if total == 0 {
        return 0, 0
    }
    
    p := float64(wins) / float64(total)
    
    // Calculate z-score for given confidence level
    var z float64
    switch confidence {
    case 0.90:
        z = 1.645
    case 0.95:
        z = 1.96
    case 0.99:
        z = 2.576
    default:
        z = 1.96 // Default to 95%
    }
    
    // Calculate margin of error
    margin := z * math.Sqrt(p*(1-p)/float64(total))
    
    return math.Max(0, p-margin), math.Min(1, p+margin)
}
```

### Validation and Error Handling

#### Error Types
```go
var (
    ErrInvalidPhase           = errors.New("invalid game phase for this operation")
    ErrInvalidDoor           = errors.New("invalid door number (must be 1, 2, or 3)")
    ErrInvalidStrategy       = errors.New("invalid strategy")
    ErrGameAlreadyComplete   = errors.New("game is already complete")
    ErrNoValidDoorToReveal   = errors.New("no valid door available for host to reveal")
    ErrGameNotInitialized    = errors.New("game not properly initialized")
)
```

#### Validation Functions
```go
func (g *Game) Validate() error {
    // Check basic structure
    if g.Doors == nil || len(g.Doors) != 3 {
        return ErrGameNotInitialized
    }
    
    // Validate prize placement
    prizeCount := 0
    for _, door := range g.Doors {
        if door.HasPrize {
            prizeCount++
        }
    }
    if prizeCount != 1 {
        return errors.New("exactly one door must have the prize")
    }
    
    // Validate phase-specific constraints
    switch g.Phase {
    case PhasePlayerChoice:
        if g.PlayerChoice < 0 || g.PlayerChoice > 2 {
            return errors.New("invalid player choice")
        }
    case PhaseHostReveal:
        if g.HostRevealed < 0 || g.HostRevealed > 2 {
            return errors.New("invalid host revealed door")
        }
        if g.HostRevealed == g.PlayerChoice {
            return errors.New("host cannot reveal player's chosen door")
        }
        if g.Doors[g.HostRevealed].HasPrize {
            return errors.New("host cannot reveal door with prize")
        }
    case PhaseComplete:
        if g.Result == nil {
            return errors.New("complete game must have result")
        }
    }
    
    return nil
}

func ValidateDoorNumber(doorNumber int) error {
    if doorNumber < 1 || doorNumber > 3 {
        return ErrInvalidDoor
    }
    return nil
}

func ValidateStrategy(strategy Strategy) error {
    if strategy != StrategyStick && strategy != StrategySwitch {
        return ErrInvalidStrategy
    }
    return nil
}
```

### Random Number Generation

#### Secure Randomization
```go
import (
    "crypto/rand"
    "math/big"
)

func SecureRandomInt(max int) (int, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
    if err != nil {
        return 0, err
    }
    return int(n.Int64()), nil
}

func PlacePrizeRandomly() (int, error) {
    return SecureRandomInt(3)
}
```

#### Deterministic Testing
```go
type DeterministicRandom struct {
    sequence []int
    index    int
}

func NewDeterministicRandom(sequence []int) *DeterministicRandom {
    return &DeterministicRandom{
        sequence: sequence,
        index:    0,
    }
}

func (dr *DeterministicRandom) Intn(n int) int {
    if dr.index >= len(dr.sequence) {
        dr.index = 0 // Wrap around
    }
    result := dr.sequence[dr.index] % n
    dr.index++
    return result
}
```

### Game Serialization

#### JSON Serialization
```go
func (g *Game) MarshalJSON() ([]byte, error) {
    type Alias Game
    return json.Marshal(&struct {
        *Alias
        Version string `json:"version"`
    }{
        Alias:   (*Alias)(g),
        Version: "1.0",
    })
}

func (g *Game) UnmarshalJSON(data []byte) error {
    type Alias Game
    aux := &struct {
        *Alias
        Version string `json:"version"`
    }{
        Alias: (*Alias)(g),
    }
    
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    
    // Validate version compatibility
    if aux.Version != "1.0" {
        return errors.New("unsupported game version")
    }
    
    return g.Validate()
}
```

### Performance Considerations

#### Memory Efficiency
- Use array instead of slice for doors (fixed size)
- Minimize allocations during gameplay
- Reuse game objects when possible

#### CPU Efficiency
- Avoid unnecessary calculations
- Cache computed values
- Use efficient algorithms for door selection

#### Concurrency Safety
```go
type SafeGame struct {
    mu   sync.RWMutex
    game *Game
}

func (sg *SafeGame) SelectDoor(doorNumber int) error {
    sg.mu.Lock()
    defer sg.mu.Unlock()
    return sg.game.SelectDoor(doorNumber)
}

func (sg *SafeGame) GetState() GameState {
    sg.mu.RLock()
    defer sg.mu.RUnlock()
    return sg.game.GetCurrentState()
}
```

### Testing Strategy

#### Unit Tests
- Test each game phase transition
- Validate all error conditions
- Test probability distributions
- Test serialization/deserialization

#### Property-Based Tests
- Verify game rules hold across random inputs
- Test statistical convergence over many games
- Validate state machine invariants

#### Integration Tests
- Test complete game flows
- Verify UI integration
- Test statistics tracking accuracy