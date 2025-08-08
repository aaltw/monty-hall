package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// Enhanced color palette with gradients and effects
var (
	// Primary colors
	PrimaryColor   = lipgloss.Color("#00ADD8") // Go blue
	SecondaryColor = lipgloss.Color("#00D084") // Success green
	AccentColor    = lipgloss.Color("#FF6B6B") // Attention red
	WarningColor   = lipgloss.Color("#FFA726") // Warning orange

	// Neutral colors
	TextColor       = lipgloss.Color("#FFFFFF") // White text
	MutedColor      = lipgloss.Color("#888888") // Muted gray
	BorderColor     = lipgloss.Color("#444444") // Border gray
	BackgroundColor = lipgloss.Color("#1A1A1A") // Dark background

	// Game-specific colors
	CarColor      = lipgloss.Color("#FFD700") // Gold for car
	GoatColor     = lipgloss.Color("#8B4513") // Brown for goat
	DoorColor     = lipgloss.Color("#8B4513") // Wood brown for doors
	SelectedColor = lipgloss.Color("#00ADD8") // Highlight color

	// Enhanced visual colors
	GlowColor      = lipgloss.Color("#00FFFF") // Cyan glow
	SparkleColor   = lipgloss.Color("#FFFF00") // Yellow sparkle
	ShadowColor    = lipgloss.Color("#000000") // Black shadow
	HighlightColor = lipgloss.Color("#FFFFFF") // White highlight

	// Gradient color sets
	WinGradient = []lipgloss.Color{
		lipgloss.Color("#FFD700"), // Gold
		lipgloss.Color("#FFA500"), // Orange
		lipgloss.Color("#FF6347"), // Tomato
	}

	DoorGradient = []lipgloss.Color{
		lipgloss.Color("#8B4513"), // Saddle brown
		lipgloss.Color("#A0522D"), // Sienna
		lipgloss.Color("#CD853F"), // Peru
	}

	MenuGradient = []lipgloss.Color{
		lipgloss.Color("#1A1A1A"), // Dark
		lipgloss.Color("#2A2A2A"), // Medium
		lipgloss.Color("#3A3A3A"), // Light
	}
)

// Base styles
var (
	// Container styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2).
			Margin(1, 0)

	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Margin(1, 0)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Align(lipgloss.Center)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Interactive styles
	MenuItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(SelectedColor).
				Bold(true).
				Padding(0, 2).
				Background(lipgloss.Color("#2A2A2A"))

	// Flat, minimalistic menu buttons for Phase 3
	MenuButtonStyle = lipgloss.NewStyle().
			Width(24).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(TextColor).
			Margin(0, 0).
			Padding(1, 2)

	SelectedMenuButtonStyle = MenuButtonStyle.
				Foreground(PrimaryColor).
				Background(lipgloss.Color("#2A2A2A")).
				Bold(true)

	// Door styles - no width/height constraints to prevent Unicode collapse
	DoorStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(DoorColor).
			Background(lipgloss.Color("#2D1B0E")).
			Padding(0, 1)

	SelectedDoorStyle = DoorStyle.
				BorderForeground(SelectedColor).
				Background(lipgloss.Color("#1A3A3A")).
				Bold(true)

	OpenDoorStyle = DoorStyle.
			BorderForeground(SecondaryColor).
			Background(lipgloss.Color("#1A2A1A"))

	WinningDoorStyle = DoorStyle.
				BorderForeground(CarColor).
				Background(lipgloss.Color("#2A2A1A")).
				Bold(true)

	// Statistics styles
	StatsHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor).
				Underline(true)

	StatsValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(SecondaryColor)

	StatsLabelStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// Progress bar styles
	ProgressBarStyle = lipgloss.NewStyle().
				Width(30).
				Height(1).
				Background(lipgloss.Color("#333333"))

	ProgressFillStyle = lipgloss.NewStyle().
				Background(PrimaryColor)

	// Enhanced visual effect styles
	GlowStyle = lipgloss.NewStyle().
			Foreground(GlowColor).
			Bold(true)

	SparkleStyle = lipgloss.NewStyle().
			Foreground(SparkleColor).
			Bold(true)

	WinningStyle = lipgloss.NewStyle().
			Foreground(CarColor).
			Background(lipgloss.Color("#2A1A00")).
			Bold(true).
			Blink(true)

	// Animation-ready door styles
	DoorClosedStyle = DoorStyle.
			BorderForeground(DoorColor).
			Background(lipgloss.Color("#2D1B0E"))

	DoorOpeningStyle = DoorStyle.
				BorderForeground(WarningColor).
				Background(lipgloss.Color("#2A2A1A")).
				Bold(true)

	DoorRevealedStyle = DoorStyle.
				BorderForeground(SecondaryColor).
				Background(lipgloss.Color("#1A2A1A")).
				Bold(true)

	// Particle effect styles
	ParticleStyle = lipgloss.NewStyle().
			Foreground(SparkleColor)

	// Typewriter effect style
	TypewriterStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// Pulse effect styles
	PulseBaseStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	PulseActiveStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true)
)

// Layout helpers
func CenterHorizontal(content string, width int) string {
	return lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, content)
}

func CenterVertical(content string, height int) string {
	return lipgloss.Place(1, height, lipgloss.Center, lipgloss.Center, content)
}

func Center(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// SafeCenter centers content without truncating it if it's wider than the available width
func SafeCenter(content string, width int) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return content
	}

	// Find the maximum line width
	maxLineWidth := 0
	for _, line := range lines {
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxLineWidth {
			maxLineWidth = lineWidth
		}
	}

	// If content is wider than available width, don't center (avoid truncation)
	if maxLineWidth >= width {
		return content
	}

	// Otherwise, center normally
	return lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, content)
}

// Responsive layout
func GetLayoutWidth(terminalWidth int) int {
	if terminalWidth < 80 {
		return terminalWidth - 4
	} else if terminalWidth < 120 {
		return 80
	}
	return 100
}

func GetLayoutHeight(terminalHeight int) int {
	if terminalHeight < 24 {
		return terminalHeight - 2
	}
	return terminalHeight - 4
}

// Animation helpers
func PulseStyle(baseStyle lipgloss.Style, intensity float64) lipgloss.Style {
	// Simple pulse effect by adjusting brightness
	if intensity > 0.5 {
		return baseStyle.Copy().Bold(true)
	}
	return baseStyle
}

// Game phase styling
func GetPhaseStyle(phase string) lipgloss.Style {
	switch phase {
	case "Initial Choice":
		return lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	case "Final Choice":
		return lipgloss.NewStyle().Foreground(WarningColor).Bold(true)
	case "Game Over":
		return lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(TextColor)
	}
}

// Utility functions for consistent spacing
func Spacer(height int) string {
	return lipgloss.NewStyle().Height(height).Render("")
}

func Divider(width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Foreground(BorderColor).
		Render(lipgloss.NewStyle().Width(width).Render("─"))
}

// Enhanced visual effect utilities
func CreateGlowEffect(text string, intensity float64) string {
	if intensity <= 0 {
		return text
	}

	// Simple glow effect using color intensity
	glowStyle := lipgloss.NewStyle().Foreground(GlowColor)
	if intensity > 0.5 {
		glowStyle = glowStyle.Bold(true)
	}

	return glowStyle.Render(text)
}

func CreateSparkleEffect(count int) string {
	sparkles := []string{"✨", "⭐", "💫", "🌟", "✦", "✧"}
	if count <= 0 || count > len(sparkles) {
		return ""
	}

	result := ""
	for i := 0; i < count; i++ {
		if i < len(sparkles) {
			result += SparkleStyle.Render(sparkles[i])
		}
	}
	return result
}

func InterpolateColor(color1, color2 lipgloss.Color, progress float64) lipgloss.Color {
	// Simple color interpolation (for basic effects)
	// In a more advanced implementation, you'd parse RGB values and interpolate
	if progress <= 0.5 {
		return color1
	}
	return color2
}

func CreatePulseEffect(text string, progress float64) string {
	// Create pulsing effect based on progress (0.0 to 1.0)
	intensity := (1 + math.Sin(progress*math.Pi*4)) / 2 // Oscillate between 0 and 1

	if intensity > 0.7 {
		return PulseActiveStyle.Render(text)
	} else if intensity > 0.3 {
		return lipgloss.NewStyle().Foreground(PrimaryColor).Render(text)
	}
	return PulseBaseStyle.Render(text)
}

func CreateShadowEffect(text string, shadowChar rune) string {
	// Simple shadow effect by adding shadow characters
	shadow := lipgloss.NewStyle().Foreground(ShadowColor).Render(string(shadowChar))
	return text + shadow
}

// Gradient utilities
func ApplyGradient(text string, colors []lipgloss.Color, progress float64) string {
	if len(colors) == 0 {
		return text
	}

	if len(colors) == 1 {
		return lipgloss.NewStyle().Foreground(colors[0]).Render(text)
	}

	// Select color based on progress
	index := int(progress * float64(len(colors)-1))
	if index >= len(colors) {
		index = len(colors) - 1
	}

	return lipgloss.NewStyle().Foreground(colors[index]).Render(text)
}

// Animation state styling
func GetAnimationStyle(state string, progress float64) lipgloss.Style {
	switch state {
	case "opening":
		return DoorOpeningStyle
	case "revealed":
		return DoorRevealedStyle
	case "winning":
		return WinningStyle
	case "glowing":
		return GlowStyle
	default:
		return DoorClosedStyle
	}
}

// Enhanced gradient and glow effects for Phase 4

// CreateRainbowText creates text with rainbow gradient effect
func CreateRainbowText(text string) string {
	rainbowColors := []lipgloss.Color{
		lipgloss.Color("#FF0000"), // Red
		lipgloss.Color("#FF7F00"), // Orange
		lipgloss.Color("#FFFF00"), // Yellow
		lipgloss.Color("#00FF00"), // Green
		lipgloss.Color("#0000FF"), // Blue
		lipgloss.Color("#4B0082"), // Indigo
		lipgloss.Color("#9400D3"), // Violet
	}

	if len(text) == 0 {
		return text
	}

	result := ""
	for i, char := range text {
		colorIndex := i % len(rainbowColors)
		style := lipgloss.NewStyle().Foreground(rainbowColors[colorIndex])
		result += style.Render(string(char))
	}

	return result
}

// CreateGradientText creates text with a gradient between two colors
func CreateGradientText(text string, startColor, endColor lipgloss.Color) string {
	if len(text) == 0 {
		return text
	}

	result := ""
	for i, char := range text {
		// Simple gradient interpolation (alternating for now)
		var color lipgloss.Color
		if i%2 == 0 {
			color = startColor
		} else {
			color = endColor
		}

		style := lipgloss.NewStyle().Foreground(color)
		result += style.Render(string(char))
	}

	return result
}

// CreateGlowText creates text with a glow effect
func CreateGlowText(text string, glowColor lipgloss.Color) string {
	glowStyle := lipgloss.NewStyle().
		Foreground(glowColor).
		Bold(true)

	return glowStyle.Render(text)
}

// CreatePulsingText creates text that pulses with different intensities
func CreatePulsingText(text string, baseColor, pulseColor lipgloss.Color, intensity float64) string {
	// Use intensity to determine which color to use
	var color lipgloss.Color
	if intensity > 0.5 {
		color = pulseColor
	} else {
		color = baseColor
	}

	style := lipgloss.NewStyle().Foreground(color)
	if intensity > 0.7 {
		style = style.Bold(true)
	}

	return style.Render(text)
}

// Enhanced title styling for Phase 4
func CreateEnhancedTitle(text string) string {
	return CreateGradientText(text, PrimaryColor, SecondaryColor)
}

// Enhanced winning message with effects
func CreateWinningMessage(text string) string {
	return CreateGlowText(CreateRainbowText(text), CarColor)
}

// Enhanced door number with glow effect
func CreateGlowingDoorNumber(number int, isSelected bool) string {
	text := fmt.Sprintf("%d", number)
	if isSelected {
		return CreateGlowText(text, SelectedColor)
	}
	return text
}

// Responsive design enhancements for Phase 4

// ScreenSize represents different terminal size categories
type ScreenSize int

const (
	ScreenSmall  ScreenSize = iota // < 80 columns
	ScreenMedium                   // 80-120 columns
	ScreenLarge                    // > 120 columns
)

// ResponsiveLayout contains layout parameters for different screen sizes
type ResponsiveLayout struct {
	DoorWidth   int
	DoorSpacing int
	ShowDetails bool
	CompactMode bool
	MaxWidth    int
	Padding     int
}

// DetectScreenSize determines the screen size category
func DetectScreenSize(width, height int) ScreenSize {
	if width < 80 {
		return ScreenSmall
	} else if width < 120 {
		return ScreenMedium
	}
	return ScreenLarge
}

// GetResponsiveLayout returns layout parameters for the given screen size
func GetResponsiveLayout(size ScreenSize) ResponsiveLayout {
	switch size {
	case ScreenSmall:
		return ResponsiveLayout{
			DoorWidth:   8,
			DoorSpacing: 1,
			ShowDetails: false,
			CompactMode: true,
			MaxWidth:    70,
			Padding:     1,
		}
	case ScreenMedium:
		return ResponsiveLayout{
			DoorWidth:   12,
			DoorSpacing: 2,
			ShowDetails: true,
			CompactMode: false,
			MaxWidth:    100,
			Padding:     2,
		}
	case ScreenLarge:
		return ResponsiveLayout{
			DoorWidth:   14,
			DoorSpacing: 3,
			ShowDetails: true,
			CompactMode: false,
			MaxWidth:    120,
			Padding:     3,
		}
	}
	return ResponsiveLayout{}
}

// AdaptContentToScreen adapts content based on screen size
func AdaptContentToScreen(content string, layout ResponsiveLayout) string {
	if layout.CompactMode {
		// Simplify content for small screens
		lines := strings.Split(content, "\n")
		var adaptedLines []string

		for _, line := range lines {
			if len(line) > layout.MaxWidth {
				// Truncate long lines
				line = line[:layout.MaxWidth-3] + "..."
			}
			adaptedLines = append(adaptedLines, line)
		}

		return strings.Join(adaptedLines, "\n")
	}

	return content
}

// GetResponsiveDoorStyle returns door styling based on screen size
func GetResponsiveDoorStyle(size ScreenSize, isSelected, isCursor bool) lipgloss.Style {
	layout := GetResponsiveLayout(size)

	// Remove width constraint to prevent Unicode collapse - door content handles its own width
	baseStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, layout.Padding)

	if isCursor {
		return baseStyle.Copy().
			BorderForeground(SelectedColor).
			Background(lipgloss.Color("#1A3A3A")).
			Bold(true)
	} else if isSelected {
		return baseStyle.Copy().
			BorderForeground(SecondaryColor).
			Background(lipgloss.Color("#1A2A1A"))
	}

	return baseStyle.Copy().
		BorderForeground(DoorColor).
		Background(lipgloss.Color("#2D1B0E"))
}

// CreateASCIIBanner creates a large ASCII art banner for the title
func CreateASCIIBanner(width int) string {
	// Large ASCII banner for wide terminals (120+ chars)
	largeBanner := []string{
		"███╗   ███╗ ██████╗ ███╗   ██╗████████╗██╗   ██╗    ██╗  ██╗ █████╗ ██╗     ██╗     ",
		"████╗ ████║██╔═══██╗████╗  ██║╚══██╔══╝╚██╗ ██╔╝    ██║  ██║██╔══██╗██║     ██║     ",
		"██╔████╔██║██║   ██║██╔██╗ ██║   ██║    ╚████╔╝     ███████║███████║██║     ██║     ",
		"██║╚██╔╝██║██║   ██║██║╚██╗██║   ██║     ╚██╔╝      ██╔══██║██╔══██║██║     ██║     ",
		"██║ ╚═╝ ██║╚██████╔╝██║ ╚████║   ██║      ██║       ██║  ██║██║  ██║███████╗███████╗",
		"╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝      ╚═╝       ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝",
	}

	// Medium ASCII banner for medium terminals (80-119 chars)
	mediumBanner := []string{
		"███╗   ███╗ ██████╗ ███╗   ██╗████████╗██╗   ██╗",
		"████╗ ████║██╔═══██╗████╗  ██║╚══██╔══╝╚██╗ ██╔╝",
		"██╔████╔██║██║   ██║██╔██╗ ██║   ██║    ╚████╔╝ ",
		"██║╚██╔╝██║██║   ██║██║╚██╗██║   ██║     ╚██╔╝  ",
		"██║ ╚═╝ ██║╚██████╔╝██║ ╚████║   ██║      ██║   ",
		"╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝      ╚═╝   ",
		"",
		"██╗  ██╗ █████╗ ██╗     ██╗     ",
		"██║  ██║██╔══██╗██║     ██║     ",
		"███████║███████║██║     ██║     ",
		"██╔══██║██╔══██║██║     ██║     ",
		"██║  ██║██║  ██║███████╗███████╗",
		"╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝",
	}

	// Small banner for narrow terminals (< 80 chars)
	smallBanner := []string{
		"╔═╗╔═╗╔╗╔╔╦╗╦ ╦  ╦ ╦╔═╗╦  ╦  ",
		"║║║║ ║║║║ ║ ╚╦╝  ╠═╣╠═╣║  ║  ",
		"╩ ╩╚═╝╝╚╝ ╩  ╩   ╩ ╩╩ ╩╩═╝╩═╝",
	}

	var selectedBanner []string
	var style lipgloss.Style

	if width >= 120 {
		selectedBanner = largeBanner
		style = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Align(lipgloss.Center)
	} else if width >= 80 {
		selectedBanner = mediumBanner
		style = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Align(lipgloss.Center)
	} else {
		selectedBanner = smallBanner
		style = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Align(lipgloss.Center)
	}

	// Apply styling to each line
	var styledLines []string
	for _, line := range selectedBanner {
		styledLines = append(styledLines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledLines...)
}

// CreateGameBanner creates a banner for the game screen
func CreateGameBanner(width int) string {
	if width >= 100 {
		// Large game banner
		banner := []string{
			"╔═╗╔═╗╔╦╗╔═╗  ╔╦╗╦╔╦╗╔═╗",
			"║ ╦╠═╣║║║║╣    ║ ║║║║║╣ ",
			"╚═╝╩ ╩╩ ╩╚═╝   ╩ ╩╩ ╩╚═╝",
		}

		style := lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Align(lipgloss.Center)

		var styledLines []string
		for _, line := range banner {
			styledLines = append(styledLines, style.Render(line))
		}

		return lipgloss.JoinVertical(lipgloss.Center, styledLines...)
	}

	// Fallback to regular header for smaller screens
	return HeaderStyle.Render("MONTY HALL GAME")
}

// CreateStatsBanner creates a banner for the statistics screen
func CreateStatsBanner(width int) string {
	if width >= 100 {
		// Large stats banner
		banner := []string{
			"╔═╗╔╦╗╔═╗╔╦╗╦╔═╗╔╦╗╦╔═╗╔═╗",
			"╚═╗ ║ ╠═╣ ║ ║╚═╗ ║ ║║  ╚═╗",
			"╚═╝ ╩ ╩ ╩ ╩ ╩╚═╝ ╩ ╩╚═╝╚═╝",
		}

		style := lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true).
			Align(lipgloss.Center)

		var styledLines []string
		for _, line := range banner {
			styledLines = append(styledLines, style.Render(line))
		}

		return lipgloss.JoinVertical(lipgloss.Center, styledLines...)
	}

	// Fallback to regular header for smaller screens
	return HeaderStyle.Render("STATISTICS")
}
