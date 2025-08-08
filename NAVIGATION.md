# Navigation Guide

## 🎮 Consistent Button Layout

All screens now follow a standardized button order for better user experience:

1. **Enter** - Primary action (Play/Select/Confirm)
2. **s** - Secondary action (Statistics/Switch)
3. **r** - Reset/Restart (where applicable)
4. **q** - Quit/Menu

## 📱 Screen Navigation

### Main Menu
- **Enter**: Select highlighted option
- **↑↓**: Navigate menu options
- **q**: Quit application

### Game Screen (Initial Choice)
- **Enter**: Select highlighted door
- **s**: View statistics
- **←→**: Navigate between doors
- **q**: Return to main menu

### Game Screen (Final Choice)
- **Enter**: Stay with original choice
- **s**: Switch to other door
- **q**: Return to main menu

### Game Screen (Results)
- **Enter**: Play again
- **s**: View statistics
- **q**: Return to main menu

### Statistics Screen
- **Enter**: Start new game
- **r**: Reset statistics
- **q**: Return to main menu

### Help Screen
- **ESC**: Return to previous screen
- **q**: Return to previous screen

## 🔄 Navigation Flow

```
Main Menu ←→ Game Screen ←→ Statistics Screen
    ↓              ↓              ↓
Help Screen    Game Results   Reset Stats
```

### Key Navigation Paths:
- **Menu → Game**: Select "Play Game" and press Enter
- **Game → Stats**: Press 's' at any time during game
- **Stats → Game**: Press Enter to start new game
- **Any Screen → Menu**: Press 'q'
- **Any Screen → Help**: Press 'h'

## 🔧 Technical Improvements

### Fixed Button Order Issue:
- **Problem**: Arrow key navigation was causing button order to change randomly
- **Root Cause**: Go maps don't preserve insertion order, causing random iteration
- **Solution**: Replaced `map[string]string` with ordered `[]KeyBinding` slice
- **Result**: Button order is now consistent and predictable across all screens

### Fixed Context-Aware 'q' Key:
- **Problem**: 'q' key always quit the application, even when footer said "Menu"
- **Root Cause**: Global key handler always called `tea.Quit` for 'q' key
- **Solution**: Made 'q' key context-aware based on current screen
- **Result**: 
  - From main menu: Quits application
  - From any other screen: Returns to main menu
  - From help: Returns to previous screen

## ✨ Navigation Improvements

### Fixed Issues:
1. ✅ **Statistics to Play**: Enter key now starts new game from stats
2. ✅ **Play to Statistics**: 's' key works in all game phases (except when switching)
3. ✅ **Consistent Button Order**: Same order across all screens
4. ✅ **Clear Instructions**: Footers show available actions in logical order

### Smart Context Switching:
- **'s' key behavior**:
  - In Final Choice phase: Switch doors
  - In all other phases: View statistics
- **Enter key behavior**:
  - In menus: Select option
  - In game: Confirm selection/choice
  - In stats: Start new game
  - In results: Play again

The navigation is now intuitive and consistent across the entire application! 🎯