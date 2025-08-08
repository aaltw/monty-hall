# Phase-by-Phase Development Plan

## Phase 1: Core Game Logic Foundation (2-3 days)

### Objectives
- Implement the fundamental Monty Hall game mechanics
- Create robust game state management
- Establish testing foundation

### Deliverables

#### 1.1 Game State Management
**File**: `pkg/game/game.go`
```go
type Game struct {
    doors         [3]*Door
    playerChoice  int
    hostRevealed  int
    gamePhase     Phase
    prizeLocation int
    created       time.Time
}

// Core methods
func NewGame() *Game
func (g *Game) SelectDoor(doorNum int) error
func (g *Game) HostRevealDoor() (int, error)
func (g *Game) MakeFinalChoice(strategy Strategy) (*Result, error)
func (g *Game) GetCurrentState() GameState
```

#### 1.2 Door Management
**File**: `pkg/game/door.go`
```go
type Door struct {
    Number   int
    IsOpen   bool
    HasPrize bool
    State    DoorState
}

type DoorState int
const (
    DoorClosed DoorState = iota
    DoorOpenWithPrize
    DoorOpenWithGoat
)
```

#### 1.3 Host Logic
**File**: `pkg/game/host.go`
```go
type Host struct{}

func (h *Host) ChooseDoorToReveal(game *Game) int
func (h *Host) CanRevealDoor(doorNum int, game *Game) bool
```

#### 1.4 Game Rules Engine
**File**: `pkg/game/rules.go`
```go
type Strategy int
const (
    StrategyStick Strategy = iota
    StrategySwitch
)

type Result struct {
    Won           bool
    Strategy      Strategy
    InitialChoice int
    FinalChoice   int
    PrizeLocation int
}

func ValidateGameState(game *Game) error
func CalculateResult(game *Game, strategy Strategy) *Result
```

#### 1.5 Testing Suite
**Files**: `pkg/game/*_test.go`
- Unit tests for all game logic
- Property-based testing for randomization
- Edge case validation
- Statistical validation over 10,000+ games

### Acceptance Criteria
- [ ] Game correctly implements Monty Hall rules
- [ ] Random prize placement is uniform
- [ ] Host always reveals a goat door
- [ ] Host never reveals the player's chosen door
- [ ] Statistical results converge to 1/3 stick, 2/3 switch over many games
- [ ] All edge cases handled gracefully
- [ ] 100% test coverage for game logic
- [ ] **Agent Code Review**: Core game logic reviewed for correctness and Go best practices

### CLI Testing Tool
**File**: `cmd/test-game/main.go`
```go
// Simple CLI tool to test game mechanics
func main() {
    // Run statistical validation
    // Test edge cases
    // Verify probability distributions
}
```

---

## Phase 2: Statistics and Data Tracking (1-2 days)

### Objectives
- Implement comprehensive statistics tracking
- Add data persistence capabilities
- Create statistical analysis tools

### Deliverables

#### 2.1 Statistics Engine
**File**: `pkg/stats/tracker.go`
```go
type Statistics struct {
    TotalGames      int                 `json:"total_games"`
    StickWins       int                 `json:"stick_wins"`
    StickLosses     int                 `json:"stick_losses"`
    SwitchWins      int                 `json:"switch_wins"`
    SwitchLosses    int                 `json:"switch_losses"`
    SessionStats    SessionStats        `json:"session_stats"`
    GameHistory     []GameRecord        `json:"game_history,omitempty"`
    LastUpdated     time.Time           `json:"last_updated"`
}

type SessionStats struct {
    GamesPlayed     int       `json:"games_played"`
    StartTime       time.Time `json:"start_time"`
    StickWins       int       `json:"stick_wins"`
    SwitchWins      int       `json:"switch_wins"`
}

func NewStatistics() *Statistics
func (s *Statistics) RecordGame(result *game.Result)
func (s *Statistics) GetStickWinRate() float64
func (s *Statistics) GetSwitchWinRate() float64
func (s *Statistics) Reset()
func (s *Statistics) ResetSession()
```

#### 2.2 Statistical Calculations
**File**: `pkg/stats/calculator.go`
```go
type StatisticalAnalysis struct {
    StickWinRate        float64
    SwitchWinRate       float64
    StickConfidenceInt  ConfidenceInterval
    SwitchConfidenceInt ConfidenceInterval
    SampleSize          int
    IsSignificant       bool
}

type ConfidenceInterval struct {
    Lower      float64
    Upper      float64
    Confidence float64 // e.g., 0.95 for 95%
}

func CalculateWinRates(stats *Statistics) StatisticalAnalysis
func CalculateConfidenceInterval(wins, total int, confidence float64) ConfidenceInterval
func IsStatisticallySignificant(stats *Statistics) bool
```

#### 2.3 Data Persistence
**File**: `pkg/stats/persistence.go`
```go
type PersistenceManager struct {
    filePath string
}

func NewPersistenceManager(filePath string) *PersistenceManager
func (pm *PersistenceManager) SaveStatistics(stats *Statistics) error
func (pm *PersistenceManager) LoadStatistics() (*Statistics, error)
func (pm *PersistenceManager) BackupStatistics() error
func (pm *PersistenceManager) ExportCSV(stats *Statistics, filePath string) error
```

#### 2.4 Game History
**File**: `pkg/stats/history.go`
```go
type GameRecord struct {
    Timestamp     time.Time `json:"timestamp"`
    Strategy      Strategy  `json:"strategy"`
    Won           bool      `json:"won"`
    InitialChoice int       `json:"initial_choice"`
    FinalChoice   int       `json:"final_choice"`
    PrizeLocation int       `json:"prize_location"`
}

func (s *Statistics) AddGameRecord(result *game.Result)
func (s *Statistics) GetRecentGames(count int) []GameRecord
func (s *Statistics) GetGamesByStrategy(strategy Strategy) []GameRecord
```

### Acceptance Criteria
- [ ] Statistics accurately track all game outcomes
- [ ] Win rates calculated correctly
- [ ] Confidence intervals computed properly
- [ ] Data persists across application restarts
- [ ] Statistics can be reset and exported
- [ ] Session statistics separate from all-time statistics
- [ ] Statistical significance detection works
- [ ] File I/O errors handled gracefully
- [ ] **Agent Code Review**: Statistics calculations and persistence logic validated

---

## Phase 3: Basic Terminal UI (3-4 days)

### Objectives
- Create functional terminal user interface
- Implement basic navigation and interaction
- Establish UI framework and patterns

### Deliverables

#### 3.1 UI Framework Setup
**File**: `pkg/ui/app.go`
```go
type App struct {
    game        *game.Game
    stats       *stats.Statistics
    persistence *stats.PersistenceManager
    screen      Screen
    state       UIState
    width       int
    height      int
}

type Screen int
const (
    ScreenMenu Screen = iota
    ScreenGame
    ScreenStats
    ScreenTutorial
    ScreenHelp
    ScreenQuit
)

func NewApp() *App
func (a *App) Init() tea.Cmd
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (a *App) View() string
```

#### 3.2 Game Screen
**File**: `pkg/ui/screens/game.go`
```go
type GameScreen struct {
    game     *game.Game
    selected int
    message  string
    phase    game.Phase
}

func NewGameScreen(game *game.Game) *GameScreen
func (gs *GameScreen) Update(msg tea.Msg) (*GameScreen, tea.Cmd)
func (gs *GameScreen) View(width, height int) string
func (gs *GameScreen) renderDoors() string
func (gs *GameScreen) renderInstructions() string
func (gs *GameScreen) renderGamePhase() string
```

#### 3.3 Door Visualization
**File**: `pkg/ui/components/doors.go`
```go
type DoorComponent struct {
    doors    [3]*game.Door
    selected int
    revealed int
}

func NewDoorComponent(doors [3]*game.Door) *DoorComponent
func (dc *DoorComponent) Render(width int) string
func (dc *DoorComponent) renderDoor(door *game.Door, index int, selected bool) string

// ASCII art for doors
const (
    ClosedDoorArt = `
┌─────────┐
│    %d    │
│  ┌───┐  │
│  │ ○ │  │
│  └───┘  │
│         │
└─────────┘`

    OpenDoorArt = `
┌─────────┐
│    %d    │
│    %s    │
│         │
│         │
│         │
└─────────┘`
)
```

#### 3.4 Statistics Screen
**File**: `pkg/ui/screens/stats.go`
```go
type StatsScreen struct {
    stats    *stats.Statistics
    analysis stats.StatisticalAnalysis
}

func NewStatsScreen(stats *stats.Statistics) *StatsScreen
func (ss *StatsScreen) Update(msg tea.Msg) (*StatsScreen, tea.Cmd)
func (ss *StatsScreen) View(width, height int) string
func (ss *StatsScreen) renderOverallStats() string
func (ss *StatsScreen) renderSessionStats() string
func (ss *StatsScreen) renderWinRates() string
```

#### 3.5 Menu System
**File**: `pkg/ui/screens/menu.go`
```go
type MenuScreen struct {
    options  []MenuOption
    selected int
}

type MenuOption struct {
    Label  string
    Action func() tea.Cmd
    Key    string
}

func NewMenuScreen() *MenuScreen
func (ms *MenuScreen) Update(msg tea.Msg) (*MenuScreen, tea.Cmd)
func (ms *MenuScreen) View(width, height int) string
```

#### 3.6 Input Handling
**File**: `pkg/ui/input.go`
```go
type KeyMap struct {
    Up     key.Binding
    Down   key.Binding
    Left   key.Binding
    Right  key.Binding
    Enter  key.Binding
    Space  key.Binding
    Escape key.Binding
    Quit   key.Binding
    Help   key.Binding
}

func DefaultKeyMap() KeyMap
func HandleGameInput(msg tea.KeyMsg, game *game.Game) tea.Cmd
func HandleMenuInput(msg tea.KeyMsg, menu *MenuScreen) tea.Cmd
```

### Acceptance Criteria
- [ ] Application starts and displays main menu
- [ ] User can navigate between screens using keyboard
- [ ] Game can be played from start to finish
- [ ] Doors are visually represented clearly
- [ ] Statistics are displayed accurately
- [ ] Input validation prevents invalid actions
- [ ] Application handles terminal resize gracefully
- [ ] Help information is accessible
- [ ] Application can be quit cleanly
- [ ] **Agent Code Review**: UI architecture and Bubbletea integration patterns reviewed

---

## Phase 4: Enhanced Visual Experience (2-3 days)

### Objectives
- Polish the visual interface with colors and styling
- Add animations and smooth transitions
- Improve overall user experience

### Deliverables

#### 4.1 Advanced Visual Styling System
**File**: `pkg/ui/styles/colors.go`
```go
type ColorScheme struct {
    Primary     GradientColor
    Secondary   GradientColor
    Success     GradientColor
    Warning     GradientColor
    Error       GradientColor
    Background  GradientColor
    Text        lipgloss.Color
    Accent      GradientColor
}

type GradientColor struct {
    Start  lipgloss.Color
    Middle lipgloss.Color
    End    lipgloss.Color
    Angle  int
}

var (
    DefaultColors = ColorScheme{
        Primary: GradientColor{
            Start:  lipgloss.Color("#00ADD8"),
            Middle: lipgloss.Color("#0099CC"),
            End:    lipgloss.Color("#0088BB"),
            Angle:  45,
        },
        Success: GradientColor{
            Start: lipgloss.Color("#00D084"),
            End:   lipgloss.Color("#009A64"),
            Angle: 90,
        },
        Background: GradientColor{
            Start: lipgloss.Color("#1A1A1A"),
            End:   lipgloss.Color("#2D2D2D"),
            Angle: 180,
        },
        Text: lipgloss.Color("#FFFFFF"),
    }
)

func GetColorScheme() ColorScheme
func ApplyGradient(style lipgloss.Style, gradient GradientColor) lipgloss.Style
func CreateGlowEffect(color lipgloss.Color, intensity float64) lipgloss.Style
func CreatePulseAnimation(baseStyle lipgloss.Style) lipgloss.Style
```

#### 4.2 Layout System
**File**: `pkg/ui/styles/layout.go`
```go
type Layout struct {
    Width       int
    Height      int
    Padding     int
    Margin      int
    BorderStyle lipgloss.Border
}

var (
    DoorStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Margin(1)
    
    SelectedDoorStyle = DoorStyle.Copy().
        BorderForeground(lipgloss.Color("#00ADD8")).
        Bold(true)
    
    WinningDoorStyle = DoorStyle.Copy().
        BorderForeground(lipgloss.Color("#00D084")).
        Background(lipgloss.Color("#003300"))
)

func CalculateLayout(width, height int) Layout
func CenterHorizontal(content string, width int) string
func CenterVertical(content string, height int) string
```

#### 4.3 Advanced Animation System
**File**: `pkg/ui/animations/animator.go`
```go
type Animation struct {
    frames     []AnimationFrame
    duration   time.Duration
    current    int
    ticker     *time.Ticker
    easing     EasingFunction
    loop       bool
}

type AnimationFrame struct {
    Content string
    Style   lipgloss.Style
    Colors  []lipgloss.Color
}

type EasingFunction func(t float64) float64

type DoorOpenAnimation struct {
    Animation
    door        *game.Door
    revealed    bool
    glowEffect  bool
    sparkles    []Sparkle
}

type PulseAnimation struct {
    Animation
    baseColor   lipgloss.Color
    intensity   float64
    frequency   time.Duration
}

type ShimmerAnimation struct {
    Animation
    colors      []lipgloss.Color
    direction   int
    speed       time.Duration
}

func NewDoorOpenAnimation(door *game.Door) *DoorOpenAnimation
func NewPulseAnimation(color lipgloss.Color) *PulseAnimation
func NewShimmerAnimation(colors []lipgloss.Color) *ShimmerAnimation
func (a *Animation) Start() tea.Cmd
func (a *Animation) Stop()
func (a *Animation) NextFrame() AnimationFrame
func (a *Animation) IsComplete() bool
func EaseInOut(t float64) float64
func EaseBounce(t float64) float64

// Enhanced animation frames with 24-bit color transitions
var DoorOpenFrames = []AnimationFrame{
    // Frame 1: Door slightly ajar with glow
    // Frame 2: Door half open with color transition
    // Frame 3: Door fully open with reveal effect
    // Frame 4: Content revealed with sparkle effects
}
```

#### 4.4 Enhanced Door Rendering
**File**: `pkg/ui/components/enhanced_doors.go`
```go
type EnhancedDoorComponent struct {
    doors       [3]*game.Door
    selected    int
    animations  map[int]*DoorOpenAnimation
    showPrizes  bool
    glowEffect  bool
}

func (edc *EnhancedDoorComponent) RenderWithEffects(width int) string
func (edc *EnhancedDoorComponent) AddGlowEffect(doorIndex int)
func (edc *EnhancedDoorComponent) StartOpenAnimation(doorIndex int) tea.Cmd
func (edc *EnhancedDoorComponent) renderPrize(prize string) string
func (edc *EnhancedDoorComponent) renderGoat() string
func (edc *EnhancedDoorComponent) renderCar() string

// Enhanced ASCII art with colors
const (
    CarArt = `
    🚗
  ┌─────┐
  │ CAR │
  └─────┘`

    GoatArt = `
    🐐
  ┌─────┐
  │GOAT │
  └─────┘`
)
```

#### 4.5 Progress Indicators
**File**: `pkg/ui/components/progress.go`
```go
type ProgressBar struct {
    current int
    total   int
    width   int
    style   lipgloss.Style
}

type GameProgress struct {
    phase       game.Phase
    totalPhases int
    description string
}

func NewProgressBar(total, width int) *ProgressBar
func (pb *ProgressBar) SetProgress(current int)
func (pb *ProgressBar) Render() string

func NewGameProgress(phase game.Phase) *GameProgress
func (gp *GameProgress) Render() string
func (gp *GameProgress) GetPhaseDescription(phase game.Phase) string
```

#### 4.6 Responsive Design
**File**: `pkg/ui/responsive.go`
```go
type ScreenSize int
const (
    ScreenSmall  ScreenSize = iota // < 80 cols
    ScreenMedium                   // 80-120 cols
    ScreenLarge                    // > 120 cols
)

type ResponsiveLayout struct {
    size        ScreenSize
    doorWidth   int
    doorSpacing int
    showDetails bool
}

func DetectScreenSize(width, height int) ScreenSize
func GetResponsiveLayout(size ScreenSize) ResponsiveLayout
func AdaptContentToScreen(content string, layout ResponsiveLayout) string
```

### Acceptance Criteria
- [ ] Consistent color scheme throughout application
- [ ] Smooth door opening animations
- [ ] Visual feedback for user selections
- [ ] Responsive layout for different terminal sizes
- [ ] Clear visual distinction between game phases
- [ ] Attractive prize and goat representations
- [ ] Progress indicators show game advancement
- [ ] Graceful degradation on terminals without color support
- [ ] **Agent Code Review**: Visual effects implementation and performance optimization reviewed

---

## Phase 5: Educational Features (2-3 days)

### Objectives
- Add comprehensive tutorial system
- Implement educational explanations
- Create interactive learning modes

### Deliverables

#### 5.1 Tutorial System
**File**: `pkg/ui/screens/tutorial.go`
```go
type TutorialScreen struct {
    steps       []TutorialStep
    currentStep int
    game        *game.Game
    interactive bool
}

type TutorialStep struct {
    Title       string
    Description string
    Action      TutorialAction
    Interactive bool
    Validation  func(*game.Game) bool
}

type TutorialAction int
const (
    ActionExplain TutorialAction = iota
    ActionDemonstrate
    ActionInteract
    ActionValidate
)

func NewTutorialScreen() *TutorialScreen
func (ts *TutorialScreen) NextStep() bool
func (ts *TutorialScreen) PreviousStep() bool
func (ts *TutorialScreen) ExecuteCurrentStep() tea.Cmd
func (ts *TutorialScreen) RenderStep() string
```

#### 5.2 Educational Content
**File**: `pkg/education/content.go`
```go
type EducationalContent struct {
    sections map[string]Section
}

type Section struct {
    Title       string
    Content     string
    Examples    []Example
    Interactive bool
}

type Example struct {
    Description string
    Scenario    game.Game
    Explanation string
}

var (
    ProbabilityExplanation = Section{
        Title: "Why Switching Works",
        Content: `
The key insight is that your initial choice has a 1/3 probability of being correct.
This means the other two doors combined have a 2/3 probability.
When the host eliminates one losing door, that 2/3 probability transfers
entirely to the remaining unopened door.
        `,
    }
    
    MathematicalProof = Section{
        Title: "Mathematical Proof",
        Content: `
Let's prove this mathematically:
- P(initial choice correct) = 1/3
- P(initial choice wrong) = 2/3
- If initial choice is wrong, switching always wins
- Therefore: P(switching wins) = P(initial choice wrong) = 2/3
        `,
    }
)

func LoadEducationalContent() *EducationalContent
func (ec *EducationalContent) GetSection(name string) Section
func (ec *EducationalContent) RenderSection(section Section, width int) string
```

#### 5.3 Interactive Demonstrations
**File**: `pkg/education/demonstrations.go`
```go
type Demonstration struct {
    name        string
    description string
    scenarios   []DemoScenario
    autoPlay    bool
}

type DemoScenario struct {
    setup       func() *game.Game
    actions     []DemoAction
    explanation string
}

type DemoAction struct {
    description string
    execute     func(*game.Game) error
    highlight   string
}

func NewProbabilityDemo() *Demonstration
func NewStatisticalDemo() *Demonstration
func (d *Demonstration) Run() tea.Cmd
func (d *Demonstration) Pause()
func (d *Demonstration) Resume()
func (d *Demonstration) Reset()
```

#### 5.4 Learning Modes
**File**: `pkg/ui/screens/learning.go`
```go
type LearningMode int
const (
    ModeTutorial LearningMode = iota
    ModePractice
    ModeChallenge
    ModeDemonstration
)

type LearningScreen struct {
    mode         LearningMode
    challenge    *Challenge
    progress     LearningProgress
    achievements []Achievement
}

type Challenge struct {
    name        string
    description string
    target      ChallengeTarget
    progress    int
    completed   bool
}

type ChallengeTarget struct {
    gamesRequired    int
    winRateRequired  float64
    strategyRequired game.Strategy
}

func NewLearningScreen(mode LearningMode) *LearningScreen
func (ls *LearningScreen) StartChallenge(challenge *Challenge)
func (ls *LearningScreen) UpdateProgress(result *game.Result)
func (ls *LearningScreen) CheckAchievements() []Achievement
```

#### 5.5 Help System
**File**: `pkg/ui/screens/help.go`
```go
type HelpScreen struct {
    topics   map[string]HelpTopic
    current  string
    searchTerm string
}

type HelpTopic struct {
    title    string
    content  string
    keywords []string
    related  []string
}

var HelpTopics = map[string]HelpTopic{
    "gameplay": {
        Title: "How to Play",
        Content: `
1. Choose a door by pressing 1, 2, or 3
2. Watch as the host opens a door with a goat
3. Decide whether to stick with your choice or switch
4. See the result and learn from the outcome
        `,
        Keywords: []string{"play", "game", "rules", "how"},
    },
    "probability": {
        Title: "Understanding the Probability",
        Content: `
The Monty Hall problem demonstrates counterintuitive probability...
        `,
        Keywords: []string{"math", "probability", "statistics", "why"},
    },
}

func NewHelpScreen() *HelpScreen
func (hs *HelpScreen) SearchTopics(term string) []string
func (hs *HelpScreen) DisplayTopic(topic string) string
func (hs *HelpScreen) GetRelatedTopics(topic string) []string
```

#### 5.6 Achievement System
**File**: `pkg/education/achievements.go`
```go
type Achievement struct {
    id          string
    name        string
    description string
    icon        string
    unlocked    bool
    progress    float64
    requirement AchievementRequirement
}

type AchievementRequirement struct {
    gamesPlayed     int
    winRate         float64
    strategy        game.Strategy
    consecutiveWins int
}

var Achievements = []Achievement{
    {
        ID: "first_game",
        Name: "First Steps",
        Description: "Play your first game",
        Icon: "🎯",
        Requirement: AchievementRequirement{GamesPlayed: 1},
    },
    {
        ID: "switcher",
        Name: "The Switcher",
        Description: "Win 10 games using the switch strategy",
        Icon: "🔄",
        Requirement: AchievementRequirement{
            Strategy: game.StrategySwitch,
            WinRate: 1.0,
            GamesPlayed: 10,
        },
    },
}

func CheckAchievements(stats *stats.Statistics) []Achievement
func UnlockAchievement(id string) Achievement
func GetProgress(achievement Achievement, stats *stats.Statistics) float64
```

### Acceptance Criteria
- [ ] Complete tutorial guides new users through the concept
- [ ] Educational content explains the mathematics clearly
- [ ] Interactive demonstrations show probability in action
- [ ] Multiple learning modes cater to different preferences
- [ ] Help system provides comprehensive assistance
- [ ] Achievement system motivates continued learning
- [ ] Content is accurate and pedagogically sound
- [ ] **Agent Code Review**: Educational content accuracy and tutorial flow logic validated

---

## Phase 6: Advanced Features and Polish (1-2 days)

### Objectives
- Add final polish and advanced features
- Optimize performance and user experience
- Prepare for release

### Deliverables

#### 6.1 Advanced Statistics Visualization
**File**: `pkg/ui/components/charts.go`
```go
type ASCIIChart struct {
    data   []float64
    labels []string
    width  int
    height int
    style  ChartStyle
}

type ChartStyle struct {
    barChar    rune
    emptyChar  rune
    colors     []lipgloss.Color
    showValues bool
}

func NewBarChart(data []float64, labels []string) *ASCIIChart
func (c *ASCIIChart) Render() string
func (c *ASCIIChart) SetStyle(style ChartStyle)

type ConvergenceChart struct {
    stickRates   []float64
    switchRates  []float64
    gameNumbers  []int
    theoretical  map[string]float64
}

func NewConvergenceChart(stats *stats.Statistics) *ConvergenceChart
func (cc *ConvergenceChart) Render(width, height int) string
func (cc *ConvergenceChart) AddDataPoint(gameNum int, stickRate, switchRate float64)
```

#### 6.2 Configuration System
**File**: `pkg/config/config.go`
```go
type Config struct {
    UI        UIConfig        `json:"ui"`
    Game      GameConfig      `json:"game"`
    Stats     StatsConfig     `json:"stats"`
    Education EducationConfig `json:"education"`
}

type UIConfig struct {
    ColorScheme    string `json:"color_scheme"`
    AnimationSpeed int    `json:"animation_speed"`
    ShowTutorial   bool   `json:"show_tutorial"`
    AutoSave       bool   `json:"auto_save"`
}

type GameConfig struct {
    AutoAdvance      bool          `json:"auto_advance"`
    ConfirmChoices   bool          `json:"confirm_choices"`
    ShowProbability  bool          `json:"show_probability"`
    DefaultStrategy  game.Strategy `json:"default_strategy"`
}

func LoadConfig() (*Config, error)
func (c *Config) Save() error
func (c *Config) Reset() error
func DefaultConfig() *Config
```

#### 6.3 Export and Import Features
**File**: `pkg/stats/export.go`
```go
type ExportFormat int
const (
    FormatCSV ExportFormat = iota
    FormatJSON
    FormatMarkdown
)

type ExportOptions struct {
    Format      ExportFormat
    IncludeHistory bool
    DateRange   DateRange
    Anonymize   bool
}

type DateRange struct {
    Start time.Time
    End   time.Time
}

func ExportStatistics(stats *Statistics, options ExportOptions) ([]byte, error)
func ImportStatistics(data []byte, format ExportFormat) (*Statistics, error)
func GenerateReport(stats *Statistics) string
```

#### 6.4 Performance Optimization
**File**: `pkg/ui/optimization.go`
```go
type RenderCache struct {
    cache    map[string]CachedRender
    maxSize  int
    hits     int
    misses   int
}

type CachedRender struct {
    content   string
    timestamp time.Time
    hash      uint64
}

func NewRenderCache(maxSize int) *RenderCache
func (rc *RenderCache) Get(key string) (string, bool)
func (rc *RenderCache) Set(key, content string)
func (rc *RenderCache) Clear()
func (rc *RenderCache) Stats() (hits, misses int)

// Differential rendering
type ScreenDiff struct {
    previous string
    current  string
    changes  []Change
}

type Change struct {
    line   int
    column int
    old    string
    new    string
}

func CalculateDiff(previous, current string) ScreenDiff
func ApplyDiff(diff ScreenDiff) string
```

#### 6.5 Accessibility Features
**File**: `pkg/ui/accessibility.go`
```go
type AccessibilityOptions struct {
    HighContrast     bool
    LargeText        bool
    ReducedMotion    bool
    ScreenReader     bool
    ColorBlindSafe   bool
}

func DetectAccessibilityNeeds() AccessibilityOptions
func ApplyAccessibilityOptions(options AccessibilityOptions)
func GetHighContrastColors() ColorScheme
func GetColorBlindSafeColors() ColorScheme
func ConvertToScreenReaderText(content string) string
```

#### 6.6 Error Recovery and Logging
**File**: `pkg/utils/recovery.go`
```go
type ErrorRecovery struct {
    logger    *log.Logger
    handlers  map[error]RecoveryHandler
    fallbacks []FallbackAction
}

type RecoveryHandler func(error) error
type FallbackAction func() error

func NewErrorRecovery() *ErrorRecovery
func (er *ErrorRecovery) RegisterHandler(err error, handler RecoveryHandler)
func (er *ErrorRecovery) Recover(err error) error
func (er *ErrorRecovery) LogError(err error, context string)

// Graceful degradation
func HandleTerminalError(err error) error
func HandleFileSystemError(err error) error
func HandleRenderingError(err error) error
```

#### 6.7 Build and Release Automation
**File**: `scripts/build.sh`
```bash
#!/bin/bash
# Cross-platform build script
PLATFORMS="windows/amd64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64"
VERSION=$(git describe --tags --always)

for platform in $PLATFORMS; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="monty-hall-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ $GOOS = "windows" ]; then
        output+=".exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o "dist/$output" ./cmd/monty-hall
done
```

### Acceptance Criteria
- [ ] Advanced statistics with visual charts
- [ ] Comprehensive configuration system
- [ ] Export/import functionality works correctly
- [ ] Performance is optimized for smooth operation
- [ ] Accessibility features support diverse users
- [ ] Error recovery handles edge cases gracefully
- [ ] Build system produces cross-platform binaries
- [ ] Documentation is complete and accurate
- [ ] All tests pass with high coverage
- [ ] Application is ready for public release
- [ ] **Agent Code Review**: Final codebase review for production readiness and security

## Integration Testing

### End-to-End Test Scenarios
1. **Complete Game Flow**: New user → Tutorial → Multiple games → Statistics review
2. **Statistical Validation**: 1000+ games to verify probability convergence
3. **UI Responsiveness**: Test on various terminal sizes and capabilities
4. **Data Persistence**: Statistics survive application restarts
5. **Error Handling**: Graceful handling of file system errors, invalid input
6. **Performance**: Smooth operation under continuous use

### Cross-Platform Testing
- **Windows**: Command Prompt, PowerShell, Windows Terminal
- **macOS**: Terminal.app, iTerm2
- **Linux**: GNOME Terminal, Konsole, xterm

### Accessibility Testing
- **Color blindness**: Verify color-blind safe palette
- **High contrast**: Test high contrast mode
- **Screen readers**: Validate screen reader compatibility
- **Reduced motion**: Test with animations disabled

## Release Criteria

### Quality Gates
- [ ] All unit tests pass (>95% coverage)
- [ ] Integration tests pass on all platforms
- [ ] Performance benchmarks meet requirements
- [ ] Security scan shows no vulnerabilities
- [ ] Documentation is complete and accurate
- [ ] User acceptance testing completed successfully

### Release Artifacts
- [ ] Cross-platform binaries (Windows, macOS, Linux)
- [ ] Installation packages (Homebrew, Chocolatey, etc.)
- [ ] Docker container
- [ ] Source code archive
- [ ] Documentation website
- [ ] Release notes and changelog