# UI Mockups and Visual Design

## Terminal Layout Overview

The application uses a responsive design that adapts to different terminal sizes while maintaining usability and visual appeal.

### Minimum Terminal Size: 80x24 characters
### Recommended Terminal Size: 120x30 characters

## Main Menu Screen

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎯 MONTY HALL SIMULATOR 🎯                          │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                     Welcome to the Monty Hall Problem!                      │
│                                                                              │
│              Test your intuition against mathematical probability            │
│                                                                              │
│                                                                              │
│                            ┌─────────────────┐                              │
│                            │  🎮 Play Game   │                              │
│                            └─────────────────┘                              │
│                                                                              │
│                            ┌─────────────────┐                              │
│                            │  📊 Statistics  │                              │
│                            └─────────────────┘                              │
│                                                                              │
│                            ┌─────────────────┐                              │
│                            │  🎓 Tutorial    │                              │
│                            └─────────────────┘                              │
│                                                                              │
│                            ┌─────────────────┐                              │
│                            │  ❓ Help        │                              │
│                            └─────────────────┘                              │
│                                                                              │
│                            ┌─────────────────┐                              │
│                            │  🚪 Quit        │                              │
│                            └─────────────────┘                              │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Use ↑↓ arrows to navigate, Enter to select, Q to quit                       │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Game Screen - Initial Door Selection

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎯 MONTY HALL GAME 🎯                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 1/4: Choose Your Door                                                │
│  ████████████████████████████████████████████████████████████████████████    │
│                                                                              │
│     Behind one of these doors is a 🚗 CAR!                                  │
│     Behind the other two doors are 🐐 GOATS!                                │
│                                                                              │
│                                                                              │
│        ┌─────────┐         ┌─────────┐         ┌─────────┐                  │
│        │    1    │         │    2    │         │    3    │                  │
│        │  ┌───┐  │         │  ┌───┐  │         │  ┌───┐  │                  │
│        │  │ ○ │  │         │  │ ○ │  │         │  │ ○ │  │                  │
│        │  └───┘  │         │  └───┘  │         │  └───┘  │                  │
│        │         │         │         │         │         │                  │
│        └─────────┘         └─────────┘         └─────────┘                  │
│             ▲                                                               │
│         SELECTED                                                            │
│                                                                              │
│                                                                              │
│  Which door do you choose? (1, 2, or 3)                                     │
│                                                                              │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ ←→ Select door, Enter to confirm, ESC for menu, H for help                  │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Game Screen - Host Reveals Door

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎯 MONTY HALL GAME 🎯                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 2/4: Host Reveals a Door                                             │
│  ████████████████████████████████████████████████████████████████████████    │
│                                                                              │
│     You chose Door 1. Now I'll show you what's behind one of the others...  │
│                                                                              │
│                                                                              │
│        ┌─────────┐         ┌─────────┐         ┌─────────┐                  │
│        │    1    │         │    2    │         │    3    │                  │
│        │  ┌───┐  │         │    🐐    │         │  ┌───┐  │                  │
│        │  │ ○ │  │         │  ┌───┐  │         │  │ ○ │  │                  │
│        │  └───┘  │         │  │   │  │         │  └───┘  │                  │
│        │         │         │  └───┘  │         │         │                  │
│        └─────────┘         └─────────┘         └─────────┘                  │
│         YOUR CHOICE           REVEALED                                       │
│                                                                              │
│                                                                              │
│  As you can see, Door 2 has a goat! 🐐                                      │
│                                                                              │
│  Now you have a choice: STICK with Door 1, or SWITCH to Door 3?             │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Press SPACE to continue...                                                   │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Game Screen - Final Decision

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎯 MONTY HALL GAME 🎯                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 3/4: Make Your Final Decision                                        │
│  ████████████████████████████████████████████████████████████████████████    │
│                                                                              │
│     Will you STICK with your original choice or SWITCH to the other door?   │
│                                                                              │
│                                                                              │
│        ┌─────────┐                             ┌─────────┐                  │
│        │    1    │         ┌─────────┐         │    3    │                  │
│        │  ┌───┐  │         │    2    │         │  ┌───┐  │                  │
│        │  │ ○ │  │         │    🐐    │         │  │ ○ │  │                  │
│        │  └───┘  │         │  GOAT   │         │  └───┘  │                  │
│        │         │         │         │         │         │                  │
│        └─────────┘         └─────────┘         └─────────┘                  │
│         YOUR CHOICE           REVEALED            AVAILABLE                  │
│                                                                              │
│                                                                              │
│                    ┌─────────────┐  ┌─────────────┐                         │
│                    │  📌 STICK   │  │  🔄 SWITCH  │                         │
│                    │   (Door 1)  │  │   (Door 3)  │                         │
│                    └─────────────┘  └─────────────┘                         │
│                                                                              │
│  💡 Remember: Switching gives you a 2/3 chance of winning!                  │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ ←→ Select choice, Enter to confirm, P for probability explanation           │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Game Screen - Results (Win)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎯 MONTY HALL GAME 🎯                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 4/4: Results                                                         │
│  ████████████████████████████████████████████████████████████████████████    │
│                                                                              │
│                          🎉 CONGRATULATIONS! 🎉                             │
│                            YOU WON THE CAR!                                 │
│                                                                              │
│        ┌─────────┐                             ┌─────────┐                  │
│        │    1    │         ┌─────────┐         │    3    │                  │
│        │    🐐    │         │    2    │         │    🚗    │                  │
│        │  ┌───┐  │         │    🐐    │         │  ┌───┐  │                  │
│        │  │   │  │         │  GOAT   │         │  │CAR│  │                  │
│        │  └───┘  │         │         │         │  └───┘  │                  │
│        └─────────┘         └─────────┘         └─────────┘                  │
│         ORIGINAL              REVEALED            YOUR FINAL                 │
│          CHOICE                                    CHOICE                    │
│                                                                              │
│  Strategy Used: SWITCH                                                       │
│  Result: WIN! 🏆                                                            │
│                                                                              │
│  You made the smart choice! Switching gives you a 2/3 probability of        │
│  winning, while sticking only gives you 1/3.                                │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Enter to play again, S for statistics, M for menu, Q to quit                │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Statistics Screen

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          📊 GAME STATISTICS 📊                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─── Overall Statistics ────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Total Games Played: 127                                              │   │
│  │  Total Wins: 84                                                       │   │
│  │  Total Losses: 43                                                     │   │
│  │  Overall Win Rate: 66.1%                                              │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─── Strategy Comparison ───────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  STICK Strategy:                    SWITCH Strategy:                  │   │
│  │  ├─ Games: 58                       ├─ Games: 69                      │   │
│  │  ├─ Wins:  19 (32.8%)               ├─ Wins:  65 (94.2%)              │   │
│  │  └─ Losses: 39                      └─ Losses: 4                      │   │
│  │                                                                        │   │
│  │  Win Rate Comparison:                                                  │   │
│  │  STICK:  ████████░░░░░░░░░░░░ 32.8%                                   │   │
│  │  SWITCH: ████████████████████ 94.2%                                   │   │
│  │                                                                        │   │
│  │  Theoretical: STICK 33.3%, SWITCH 66.7%                              │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─── Session Statistics ────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Session Started: 2024-08-03 14:30:22                                 │   │
│  │  Games This Session: 15                                               │   │
│  │  Session Win Rate: 73.3% (11 wins, 4 losses)                         │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ R to reset stats, E to export, Enter to return to menu                      │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Tutorial Screen - Introduction

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎓 TUTORIAL - STEP 1/7 🎓                          │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                        Welcome to the Monty Hall Problem!                   │
│                                                                              │
│  The Monty Hall problem is one of the most famous probability puzzles       │
│  in mathematics. It's named after Monty Hall, the host of the TV game       │
│  show "Let's Make a Deal."                                                  │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                        │ │
│  │  The Setup:                                                            │ │
│  │                                                                        │ │
│  │  • There are three doors                                              │ │
│  │  • Behind one door is a valuable prize (a car 🚗)                     │ │
│  │  • Behind the other two doors are less valuable prizes (goats 🐐)     │ │
│  │  • You don't know which door has which prize                          │ │
│  │                                                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  The question is: Should you stick with your original choice, or switch     │
│  to the other unopened door when given the chance?                          │
│                                                                              │
│  Most people think it doesn't matter - that it's a 50/50 choice.           │
│  But mathematics tells us something very different...                       │
│                                                                              │
│  Let's find out what that is!                                               │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Press SPACE to continue, ESC to skip tutorial                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Tutorial Screen - Interactive Game

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                          🎓 TUTORIAL - STEP 4/7 🎓                          │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                           Let's Play a Practice Game!                       │
│                                                                              │
│  Now you'll play through a complete game. I'll guide you through each       │
│  step and explain what's happening.                                         │
│                                                                              │
│                                                                              │
│        ┌─────────┐         ┌─────────┐         ┌─────────┐                  │
│        │    1    │         │    2    │         │    3    │                  │
│        │  ┌───┐  │         │  ┌───┐  │         │  ┌───┐  │                  │
│        │  │ ○ │  │         │  │ ○ │  │         │  │ ○ │  │                  │
│        │  └───┘  │         │  └───┘  │         │  └───┘  │                  │
│        │         │         │         │         │         │                  │
│        └─────────┘         └─────────┘         └─────────┘                  │
│                                                                              │
│                                                                              │
│  💡 Tutorial Tip:                                                           │
│  Right now, each door has a 1/3 (33.3%) chance of having the car.          │
│  Your initial choice will have a 1/3 chance of being correct.               │
│                                                                              │
│  Choose any door to begin (1, 2, or 3):                                     │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ ←→ Select door, Enter to confirm, B for back, N for next                    │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Help Screen

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              ❓ HELP SYSTEM ❓                               │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─── Quick Start ───────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  1. Choose "Play Game" from the main menu                             │   │
│  │  2. Select a door (1, 2, or 3) using arrow keys                       │   │
│  │  3. Watch as the host reveals a goat behind another door               │   │
│  │  4. Decide whether to stick with your choice or switch                 │   │
│  │  5. See the results and learn from the outcome!                       │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─── Keyboard Controls ────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Navigation:        ↑↓←→ Arrow keys                                   │   │
│  │  Select/Confirm:    Enter or Space                                    │   │
│  │  Back/Cancel:       Escape                                            │   │
│  │  Quit Application:  Q                                                 │   │
│  │  Help:              H or ?                                            │   │
│  │  Statistics:        S (from game screen)                              │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─── Topics ────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  📖 Game Rules          🧮 Probability Math                           │   │
│  │  📊 Statistics          🎯 Strategy Tips                              │   │
│  │  🎓 Tutorial            ⚙️  Settings                                  │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Select a topic with arrow keys, Enter to view, ESC to return to menu        │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Responsive Design - Small Terminal (80x24)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                        🎯 MONTY HALL 🎯                                     │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Phase 1/4: Choose Your Door                                                 │
│ ████████████████████████████████████████████████████████████████████████     │
│                                                                              │
│ Behind one door: 🚗 CAR | Behind two doors: 🐐 GOATS                        │
│                                                                              │
│    ┌───┐      ┌───┐      ┌───┐                                              │
│    │ 1 │      │ 2 │      │ 3 │                                              │
│    │ ○ │      │ ○ │      │ ○ │                                              │
│    └───┘      └───┘      └───┘                                              │
│      ▲                                                                      │
│                                                                              │
│ Which door? (1, 2, or 3)                                                    │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ ←→ Select, Enter confirm, ESC menu, H help                                  │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Advanced Color Scheme & Visual Effects

### 24-Bit RGB Color Support
The application leverages full 24-bit color depth (16.7 million colors) through Lipgloss for rich visual experiences.

### Default Theme with Gradients
- **Primary**: Go Blue (#00ADD8) with gradient to (#0099CC)
- **Secondary**: Light Blue (#5DC9E2) with shimmer effect
- **Success**: Green gradient (#00D084 → #00B574 → #009A64)
- **Warning**: Amber gradient (#FFD23F → #FFC107 → #FF9800)
- **Error**: Red gradient (#FF6B6B → #F44336 → #D32F2F)
- **Background**: Deep gradient (#1A1A1A → #2D2D2D → #1A1A1A)
- **Text**: Pure White (#FFFFFF) with subtle glow effects
- **Accent**: Pink gradient (#CE3262 → #E91E63 → #AD1457)

### Advanced Visual Effects
- **Door Glow**: Pulsing border effects for selected doors
- **Prize Shimmer**: Golden gradient animation for car reveals
- **Goat Bounce**: Subtle animation for goat reveals  
- **Progress Gradients**: Smooth color transitions in progress bars
- **Shadow Effects**: Subtle depth with darker background gradients
- **Highlight Animations**: Color transitions on hover/selection

### High Contrast Theme
- **Primary**: Bright White (#FFFFFF)
- **Secondary**: Light Gray (#CCCCCC)
- **Success**: Bright Green (#00FF00)
- **Warning**: Bright Yellow (#FFFF00)
- **Error**: Bright Red (#FF0000)
- **Background**: Black (#000000)
- **Text**: White (#FFFFFF)
- **Accent**: Bright Cyan (#00FFFF)

### Color-Blind Safe Theme
- **Primary**: Blue (#0173B2)
- **Secondary**: Light Blue (#56B4E9)
- **Success**: Green (#029E73)
- **Warning**: Orange (#D55E00)
- **Error**: Red (#CC79A7)
- **Background**: Dark Gray (#1A1A1A)
- **Text**: White (#FFFFFF)
- **Accent**: Purple (#9467BD)

## Animation Sequences

### Door Opening Animation (4 frames, 200ms each)

**Frame 1: Door Slightly Ajar**
```
┌─────────┐
│    2   /│
│  ┌───┐/ │
│  │ ○ │  │
│  └───┘  │
│         │
└─────────┘
```

**Frame 2: Door Half Open**
```
┌─────────┐
│    2  / │
│  ┌──/┐  │
│  │ / │  │
│  └───┘  │
│    🐐    │
└─────────┘
```

**Frame 3: Door Mostly Open**
```
┌─────────┐
│   /2    │
│  /───┐  │
│ /  ○ │  │
│/  └──┘  │
│   🐐    │
└─────────┘
```

**Frame 4: Door Fully Open**
```
┌─────────┐
│ /  2    │
│/  ───┐  │
│    ○ │  │
│   └──┘  │
│   🐐    │
└─────────┘
```

## Accessibility Features

### Screen Reader Support
- All visual elements have text descriptions
- Navigation instructions are clearly announced
- Game state changes are verbally described
- Statistical information is read in logical order

### High Contrast Mode
- Increased contrast ratios for better visibility
- Bold text for important information
- Clear visual separation between elements
- Reduced reliance on color alone for information

### Reduced Motion Mode
- Disable door opening animations
- Instant state transitions
- Static progress indicators
- Simplified visual effects

### Large Text Mode
- Increased font sizes where possible
- Simplified layouts to accommodate larger text
- Reduced information density
- Clear visual hierarchy