package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/westhuis/monty-hall/pkg/game"
)

// Door component with enhanced ASCII art (Phase 3)
type DoorComponent struct {
	Number   int
	State    game.DoorState
	Content  game.DoorContent
	Selected bool
	Cursor   bool
	Width    int
	Height   int
}

// NewDoorComponent creates a new door component
func NewDoorComponent(number int, door *game.Door, selected, cursor bool) *DoorComponent {
	return &DoorComponent{
		Number:   number,
		State:    door.State,
		Content:  door.Content,
		Selected: selected,
		Cursor:   cursor,
		Width:    16, // Increased from 14 to 16 for better knob positioning
		Height:   10,
	}
}

// NewResponsiveDoorComponent creates a door component with responsive sizing
func NewResponsiveDoorComponent(number int, door *game.Door, selected, cursor bool, terminalWidth int) *DoorComponent {
	// Calculate door width based on terminal size
	doorWidth := 14 // Default width
	if terminalWidth >= 120 {
		doorWidth = 18 // Wider doors for large terminals
	} else if terminalWidth >= 100 {
		doorWidth = 16 // Medium doors for medium terminals
	}

	return &DoorComponent{
		Number:   number,
		State:    door.State,
		Content:  door.Content,
		Selected: selected,
		Cursor:   cursor,
		Width:    doorWidth,
		Height:   10,
	}
}

// Render renders the door with appropriate styling (Phase 3)
func (d *DoorComponent) Render() string {
	var style lipgloss.Style
	var content string

	// Simplified door styling - selected/cursor doors just get highlighted border
	if d.Cursor || d.Selected {
		// Cursor or selected door gets highlighted border
		style = SelectedDoorStyle.UnsetWidth().UnsetHeight()
	} else if d.State == game.Opened && d.Content == game.Car {
		// Winning door gets special styling
		style = WinningDoorStyle.UnsetWidth().UnsetHeight()
	} else {
		// Normal door
		style = DoorStyle.UnsetWidth().UnsetHeight()
	}

	// Generate the door content
	switch d.State {
	case game.Closed, game.Selected:
		// Both closed and selected doors show the same closed door content
		content = d.renderClosedDoor()
	case game.Opened:
		if d.Content == game.Car {
			content = d.renderCarDoor()
		} else {
			content = d.renderGoatDoor()
		}
	}

	return style.Render(content)
}

// RenderWithAnimation renders the door with animation support (Phase 4)
func (d *DoorComponent) RenderWithAnimation(animFrame string, animColor lipgloss.Color, isAnimating bool) string {
	var style lipgloss.Style
	var content string

	// Simplified door styling with animation support
	if isAnimating {
		// Use animation-specific styling
		style = DoorOpeningStyle.BorderForeground(animColor).UnsetWidth().UnsetHeight()
	} else if d.Cursor || d.Selected {
		// Cursor or selected door gets highlighted border
		style = SelectedDoorStyle.UnsetWidth().UnsetHeight()
	} else if d.State == game.Opened && d.Content == game.Car {
		// Winning door gets special styling
		style = WinningDoorStyle.UnsetWidth().UnsetHeight()
	} else {
		// Normal door
		style = DoorStyle.UnsetWidth().UnsetHeight()
	}

	// Generate the door content with animation support
	if isAnimating {
		content = d.renderAnimatedDoor(animFrame)
	} else {
		switch d.State {
		case game.Closed, game.Selected:
			// Both closed and selected doors show the same closed door content
			content = d.renderClosedDoor()
		case game.Opened:
			if d.Content == game.Car {
				content = d.renderCarDoor()
			} else {
				content = d.renderGoatDoor()
			}
		}
	}

	return style.Render(content)
}

// renderClosedDoor renders a closed door with responsive width
func (d *DoorComponent) renderClosedDoor() string {
	// Create door frame based on width
	topLine := "┌" + strings.Repeat("─", d.Width-2) + "┐"
	bottomLine := "└" + strings.Repeat("─", d.Width-2) + "┘"

	// Center content within the door width
	doorLabel := d.centerText("DOOR", d.Width-2)
	numberLabel := d.centerText(fmt.Sprintf("%d", d.Number), d.Width-2)
	closedLabel := d.centerText("CLOSED", d.Width-2)
	emptyLine := d.centerText("", d.Width-2)

	// Simple single-character doorknob positioned on the right with proper spacing
	knobLine := d.rightAlignText(" ●", d.Width-2) // Add space before knob for proper spacing
	emptyLine1 := d.centerText("", d.Width-2)
	emptyLine2 := d.centerText("", d.Width-2)

	doorArt := fmt.Sprintf(`%s
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
%s`, topLine, doorLabel, numberLabel, emptyLine1, knobLine, emptyLine2, emptyLine, closedLabel, bottomLine)

	// Add status indicator
	statusWidth := d.Width
	if d.Cursor {
		status := d.centerText("▶ SELECT ◀", statusWidth)
		doorArt += "\n" + status
	} else if d.Selected {
		status := d.centerText("★ CHOSEN ★", statusWidth)
		doorArt += "\n" + status
	} else {
		status := strings.Repeat(" ", statusWidth)
		doorArt += "\n" + status
	}

	return doorArt
}

// centerText centers text within a given width using proper Unicode width calculation
func (d *DoorComponent) centerText(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return runewidth.Truncate(text, width, "")
	}
	padding := (width - textWidth) / 2
	leftPad := strings.Repeat(" ", padding)
	rightPad := strings.Repeat(" ", width-textWidth-padding)
	return leftPad + text + rightPad
}

// rightAlignText aligns text to the right within a given width using proper Unicode width calculation
func (d *DoorComponent) rightAlignText(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return runewidth.Truncate(text, width, "")
	}
	leftPad := strings.Repeat(" ", width-textWidth)
	return leftPad + text
}

// renderCarDoor renders an open door with a car
func (d *DoorComponent) renderCarDoor() string {
	// Create door frame based on width
	topLine := "┌" + strings.Repeat("─", d.Width-2) + "┐"
	bottomLine := "└" + strings.Repeat("─", d.Width-2) + "┘"

	// Center content within the door width
	doorLabel := d.centerText("DOOR", d.Width-2)
	numberLabel := d.centerText(fmt.Sprintf("%d", d.Number), d.Width-2)

	// Create car ASCII art based on door width
	var carLines []string
	if d.Width >= 18 {
		// Larger car for wide doors
		carLines = []string{
			d.centerText("┌───────┐", d.Width-2),
			d.centerText("│ ░░░░░ │", d.Width-2),
			d.centerText("│░█████░│", d.Width-2),
			d.centerText("│ ░░░░░ │", d.Width-2),
			d.centerText("└───────┘", d.Width-2),
		}
	} else if d.Width >= 16 {
		// Medium car
		carLines = []string{
			d.centerText("┌─────┐", d.Width-2),
			d.centerText("│ ░░░ │", d.Width-2),
			d.centerText("│░███░│", d.Width-2),
			d.centerText("│ ░░░ │", d.Width-2),
			d.centerText("└─────┘", d.Width-2),
		}
	} else {
		// Compact car
		carLines = []string{
			d.centerText("┌───┐", d.Width-2),
			d.centerText("│░█░│", d.Width-2),
			d.centerText("└───┘", d.Width-2),
			d.centerText("", d.Width-2),
			d.centerText("", d.Width-2),
		}
	}

	doorArt := fmt.Sprintf(`%s
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
%s`, topLine, doorLabel, numberLabel, carLines[0], carLines[1], carLines[2], carLines[3], carLines[4], bottomLine)

	// Add status indicator (same height as other doors)
	statusWidth := d.Width
	if d.Cursor {
		status := d.centerText("▶ SELECT ◀", statusWidth)
		doorArt += "\n" + status
	} else if d.Selected {
		status := d.centerText("★ WIN! ★", statusWidth)
		doorArt += "\n" + status
	} else {
		status := strings.Repeat(" ", statusWidth)
		doorArt += "\n" + status
	}

	return doorArt
}

// renderGoatDoor renders an open door with a goat
func (d *DoorComponent) renderGoatDoor() string {
	// Create door frame based on width
	topLine := "┌" + strings.Repeat("─", d.Width-2) + "┐"
	bottomLine := "└" + strings.Repeat("─", d.Width-2) + "┘"

	// Center content within the door width
	doorLabel := d.centerText("DOOR", d.Width-2)
	numberLabel := d.centerText(fmt.Sprintf("%d", d.Number), d.Width-2)
	goatLabel := d.centerText("GOAT", d.Width-2)
	emptyLine := d.centerText("", d.Width-2)

	// Create goat ASCII art based on door width
	var goatLines []string
	if d.Width >= 18 {
		// Larger goat for wide doors
		goatLines = []string{
			d.centerText("(\\     /)", d.Width-2),
			d.centerText("( ^._.^ )", d.Width-2),
			d.centerText("o_(\")(\")_o", d.Width-2),
		}
	} else if d.Width >= 16 {
		// Medium goat
		goatLines = []string{
			d.centerText("(\\   /)", d.Width-2),
			d.centerText("( ._. )", d.Width-2),
			d.centerText("o_(\")(\")", d.Width-2),
		}
	} else {
		// Compact goat
		goatLines = []string{
			d.centerText("(\\ /)", d.Width-2),
			d.centerText("(._. )", d.Width-2),
			d.centerText("o_(\")", d.Width-2),
		}
	}

	doorArt := fmt.Sprintf(`%s
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
│%s│
%s`, topLine, doorLabel, numberLabel, goatLines[0], goatLines[1], goatLines[2], emptyLine, goatLabel, bottomLine)

	// Add status indicator (same pattern as other doors)
	statusWidth := d.Width
	if d.Cursor {
		status := d.centerText("▶ SELECT ◀", statusWidth)
		doorArt += "\n" + status
	} else if d.Selected {
		status := d.centerText("★ LOSE ★", statusWidth)
		doorArt += "\n" + status
	} else {
		status := d.centerText("OPENED", statusWidth)
		doorArt += "\n" + status
	}

	return doorArt
}

// DoorsRow renders all three doors in a row
func RenderDoorsRow(doors []*game.Door, playerChoice, hostOpened, cursor int, showAll bool) string {
	var doorComponents []string

	for i, door := range doors {
		selected := i == playerChoice
		isCursor := i == cursor

		// Override door state for display purposes
		displayDoor := &game.Door{
			State:   door.State,
			Content: door.Content,
		}

		// Show host opened door
		if i == hostOpened && hostOpened != -1 {
			displayDoor.State = game.Opened
		}

		// Show all doors if game is over
		if showAll {
			displayDoor.State = game.Opened
		}

		doorComp := NewDoorComponent(i+1, displayDoor, selected, isCursor)
		doorComponents = append(doorComponents, doorComp.Render())
	}

	// Join doors horizontally with center alignment to prevent collapse
	return lipgloss.JoinHorizontal(lipgloss.Center, doorComponents...)
}

// ProgressBar component
type ProgressBar struct {
	Current int
	Total   int
	Width   int
	Label   string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(current, total, width int, label string) *ProgressBar {
	return &ProgressBar{
		Current: current,
		Total:   total,
		Width:   width,
		Label:   label,
	}
}

// Render renders the progress bar
func (p *ProgressBar) Render() string {
	if p.Total == 0 {
		return ""
	}

	percentage := float64(p.Current) / float64(p.Total)
	filled := int(percentage * float64(p.Width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)

	style := ProgressBarStyle.Width(p.Width)
	progressBar := style.Render(bar)

	label := fmt.Sprintf("%s: %d/%d (%.1f%%)", p.Label, p.Current, p.Total, percentage*100)

	return lipgloss.JoinVertical(lipgloss.Left, label, progressBar)
}

// StatsCard component for displaying statistics
type StatsCard struct {
	Title  string
	Value  string
	Detail string
	Color  lipgloss.Color
}

// NewStatsCard creates a new stats card
func NewStatsCard(title, value, detail string, color lipgloss.Color) *StatsCard {
	return &StatsCard{
		Title:  title,
		Value:  value,
		Detail: detail,
		Color:  color,
	}
}

// Render renders the stats card
func (s *StatsCard) Render() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(s.Color).
		Bold(true).
		Align(lipgloss.Center)

	valueStyle := lipgloss.NewStyle().
		Foreground(TextColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginTop(1)

	detailStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Align(lipgloss.Center).
		MarginTop(1)

	cardStyle := lipgloss.NewStyle().
		Width(20).
		Height(6).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.Color).
		Padding(1).
		Align(lipgloss.Center, lipgloss.Center)

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(s.Title),
		valueStyle.Render(s.Value),
		detailStyle.Render(s.Detail),
	)

	return cardStyle.Render(content)
}

// GamePhaseIndicator shows the current game phase
type GamePhaseIndicator struct {
	Phase       game.GamePhase
	Description string
}

// NewGamePhaseIndicator creates a new phase indicator
func NewGamePhaseIndicator(phase game.GamePhase) *GamePhaseIndicator {
	var description string
	switch phase {
	case game.InitialChoice:
		description = "Choose your door"
	case game.FinalChoice:
		description = "Switch or stay?"
	case game.GameOver:
		description = "Game complete!"
	default:
		description = "Ready to play"
	}

	return &GamePhaseIndicator{
		Phase:       phase,
		Description: description,
	}
}

// Render renders the phase indicator
func (g *GamePhaseIndicator) Render() string {
	phaseStyle := GetPhaseStyle(g.Description)

	indicator := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(phaseStyle.GetForeground()).
		Padding(0, 2).
		Margin(1, 0)

	content := fmt.Sprintf("Phase: %s", g.Description)
	return indicator.Render(phaseStyle.Render(content))
}

// HelpBox component for displaying help text
type HelpBox struct {
	Title   string
	Content []string
	Width   int
}

// NewHelpBox creates a new help box
func NewHelpBox(title string, content []string, width int) *HelpBox {
	return &HelpBox{
		Title:   title,
		Content: content,
		Width:   width,
	}
}

// Render renders the help box
func (h *HelpBox) Render() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	contentStyle := lipgloss.NewStyle().
		Foreground(TextColor).
		Width(h.Width - 4).
		Align(lipgloss.Left)

	boxStyle := lipgloss.NewStyle().
		Width(h.Width).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	var lines []string
	lines = append(lines, titleStyle.Render(h.Title))

	for _, line := range h.Content {
		lines = append(lines, contentStyle.Render(line))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return boxStyle.Render(content)
}

// Banner component for the main title
func RenderBanner() string {
	banner := `
 ███╗   ███╗ ██████╗ ███╗   ██╗████████╗██╗   ██╗    ██╗  ██╗ █████╗ ██╗     ██╗     
 ████╗ ████║██╔═══██╗████╗  ██║╚══██╔══╝╚██╗ ██╔╝    ██║  ██║██╔══██╗██║     ██║     
 ██╔████╔██║██║   ██║██╔██╗ ██║   ██║    ╚████╔╝     ███████║███████║██║     ██║     
 ██║╚██╔╝██║██║   ██║██║╚██╗██║   ██║     ╚██╔╝      ██╔══██║██╔══██║██║     ██║     
 ██║ ╚═╝ ██║╚██████╔╝██║ ╚████║   ██║      ██║       ██║  ██║██║  ██║███████╗███████╗
 ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝      ╚═╝       ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝`

	bannerStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(2)

	return bannerStyle.Render(banner)
}

// MenuButton component for modern button-style menu items
type MenuButton struct {
	Text     string
	Selected bool
	Width    int
}

// NewMenuButton creates a new menu button
func NewMenuButton(text string, selected bool) *MenuButton {
	return &MenuButton{
		Text:     text,
		Selected: selected,
		Width:    24,
	}
}

// Render renders the menu button
func (m *MenuButton) Render() string {
	var style lipgloss.Style

	if m.Selected {
		style = SelectedMenuButtonStyle
	} else {
		style = MenuButtonStyle
	}

	return style.Width(m.Width).Render(m.Text)
}

// KeyBinding represents a key and its description
type KeyBinding struct {
	Key  string
	Desc string
}

// Footer component with key bindings in order
func RenderFooter(bindings []KeyBinding) string {
	var items []string

	for _, binding := range bindings {
		keyStyle := lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(MutedColor)

		item := fmt.Sprintf("%s %s", keyStyle.Render(binding.Key), descStyle.Render(binding.Desc))
		items = append(items, item)
	}

	content := strings.Join(items, " • ")

	footerStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Align(lipgloss.Center).
		MarginTop(2).
		Padding(1, 0).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(BorderColor)

	return footerStyle.Render(content)
}

// renderAnimatedDoor renders a door during animation
func (d *DoorComponent) renderAnimatedDoor(animFrame string) string {
	doorArt := fmt.Sprintf(`┌─────────┐
│  DOOR   │
│    %d    │
│         │
│    %s    │
│ OPENING │
│         │
│   ...   │
└─────────┘`, d.Number, animFrame)

	if d.Cursor {
		doorArt += "\n ▶ SELECT ◀"
	} else if d.Selected {
		doorArt += "\n  ★ CHOSEN ★"
	} else {
		doorArt += "\n            "
	}

	return doorArt
}

// RenderDoorsRowWithAnimation renders doors with animation support (Phase 4)
func RenderDoorsRowWithAnimation(doors []*game.Door, playerChoice, hostOpened, cursor int, showAll bool, model *Model) string {
	var doorComponents []string

	for i, door := range doors {
		selected := i == playerChoice
		isCursor := i == cursor

		// Override door state for display purposes
		displayDoor := &game.Door{
			State:   door.State,
			Content: door.Content,
		}

		// Show host opened door
		if i == hostOpened && hostOpened != -1 {
			displayDoor.State = game.Opened
		}

		// Show all doors if game is over
		if showAll {
			displayDoor.State = game.Opened
		}

		doorComp := NewDoorComponent(i+1, displayDoor, selected, isCursor)

		// Check if this door has an active animation
		if model.ShowAnimations {
			animFrame, animColor := model.getDoorAnimationFrame(i)
			isAnimating := model.DoorAnimations[i] != nil && model.DoorAnimations[i].IsRunning()

			if isAnimating {
				doorComponents = append(doorComponents, doorComp.RenderWithAnimation(animFrame, animColor, true))
			} else {
				doorComponents = append(doorComponents, doorComp.Render())
			}
		} else {
			doorComponents = append(doorComponents, doorComp.Render())
		}
	}

	// Join doors horizontally with center alignment to prevent collapse
	return lipgloss.JoinHorizontal(lipgloss.Center, doorComponents...)
}

// RenderDoorsRowResponsive renders doors with responsive sizing
func RenderDoorsRowResponsive(doors []*game.Door, playerChoice, hostOpened, cursor int, showAll bool, terminalWidth int) string {
	var doorComponents []string

	for i, door := range doors {
		selected := i == playerChoice
		isCursor := i == cursor

		// Override door state for display purposes
		displayDoor := &game.Door{
			State:   door.State,
			Content: door.Content,
		}

		// Show host opened door
		if i == hostOpened && hostOpened != -1 {
			displayDoor.State = game.Opened
		}

		// Show all doors if game is over
		if showAll {
			displayDoor.State = game.Opened
		}

		doorComp := NewResponsiveDoorComponent(i+1, displayDoor, selected, isCursor, terminalWidth)
		doorComponents = append(doorComponents, doorComp.Render())
	}

	// Join doors horizontally with appropriate spacing
	spacing := 2
	if terminalWidth >= 120 {
		spacing = 4 // More spacing for wide terminals
	} else if terminalWidth >= 100 {
		spacing = 3
	}

	// Add spacing between doors
	var spacedComponents []string
	for i, comp := range doorComponents {
		spacedComponents = append(spacedComponents, comp)
		if i < len(doorComponents)-1 {
			spacedComponents = append(spacedComponents, strings.Repeat(" ", spacing))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, spacedComponents...)
}

// ResetConfirmationPopover component for confirming statistics reset
type ResetConfirmationPopover struct {
	ConfirmationNumbers [4]int
	UserInputNumbers    [4]int
	CurrentInputIndex   int
	Width               int
}

// NewResetConfirmationPopover creates a new reset confirmation popover
func NewResetConfirmationPopover(confirmationNumbers [4]int, userInputNumbers [4]int, currentInputIndex int, width int) *ResetConfirmationPopover {
	return &ResetConfirmationPopover{
		ConfirmationNumbers: confirmationNumbers,
		UserInputNumbers:    userInputNumbers,
		CurrentInputIndex:   currentInputIndex,
		Width:               width,
	}
}

// Render renders the reset confirmation popover
func (r *ResetConfirmationPopover) Render() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	warningStyle := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	instructionStyle := lipgloss.NewStyle().
		Foreground(TextColor).
		Align(lipgloss.Center).
		MarginBottom(1)

	numbersStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	inputStyle := lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	footerStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Align(lipgloss.Center)

	boxStyle := lipgloss.NewStyle().
		Width(r.Width).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(WarningColor).
		Padding(2, 3).
		Align(lipgloss.Center)

	// Format confirmation numbers
	confirmationText := fmt.Sprintf("%d  %d  %d  %d",
		r.ConfirmationNumbers[0],
		r.ConfirmationNumbers[1],
		r.ConfirmationNumbers[2],
		r.ConfirmationNumbers[3])

	// Format user input with cursor
	var inputParts []string
	for i := 0; i < 4; i++ {
		if i == r.CurrentInputIndex {
			// Current input position with cursor
			if r.UserInputNumbers[i] == 0 {
				inputParts = append(inputParts, "_")
			} else {
				inputParts = append(inputParts, fmt.Sprintf("%d", r.UserInputNumbers[i]))
			}
		} else if r.UserInputNumbers[i] == 0 {
			inputParts = append(inputParts, "_")
		} else {
			inputParts = append(inputParts, fmt.Sprintf("%d", r.UserInputNumbers[i]))
		}
	}
	inputText := strings.Join(inputParts, "  ")

	// Add cursor indicator using proper Unicode width calculation
	if r.CurrentInputIndex < 4 {
		cursorPos := r.CurrentInputIndex * 3 // Account for spacing (each number + 2 spaces)
		inputRunes := []rune(inputText)
		if cursorPos < len(inputRunes) {
			// Replace the character at cursor position with cursor indicator
			inputRunes[cursorPos] = '▶'
			inputText = string(inputRunes)
		}
	}

	var lines []string
	lines = append(lines, titleStyle.Render("⚠️  RESET STATISTICS  ⚠️"))
	lines = append(lines, warningStyle.Render("This will permanently delete all game data!"))
	lines = append(lines, instructionStyle.Render("To confirm, enter these 4 numbers:"))
	lines = append(lines, numbersStyle.Render(confirmationText))
	lines = append(lines, instructionStyle.Render("Your input:"))
	lines = append(lines, inputStyle.Render(inputText))
	lines = append(lines, footerStyle.Render("Use number keys 1-9, Backspace to delete, ESC to cancel"))

	content := lipgloss.JoinVertical(lipgloss.Center, lines...)
	return boxStyle.Render(content)
}
