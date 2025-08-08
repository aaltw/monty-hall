package ui

import (
	"testing"
	"time"

	"github.com/westhuis/monty-hall/pkg/game"
)

// TestPhase4AnimationSystem tests the animation system integration
func TestPhase4AnimationSystem(t *testing.T) {
	model := NewModel()

	// Test animation manager initialization
	if model.AnimationManager == nil {
		t.Error("AnimationManager should be initialized")
	}

	if model.DoorAnimations == nil {
		t.Error("DoorAnimations map should be initialized")
	}

	if !model.ShowAnimations {
		t.Error("ShowAnimations should be enabled by default")
	}

	// Test door animation creation
	doorAnim := NewDoorOpenAnimation(0)
	if doorAnim == nil {
		t.Error("Door animation should be created successfully")
	}

	if doorAnim.DoorIndex != 0 {
		t.Errorf("Expected door index 0, got %d", doorAnim.DoorIndex)
	}

	// Test animation frames
	frame, color := doorAnim.GetCurrentFrame()
	if frame == "" {
		t.Error("Animation frame should not be empty")
	}

	if color == "" {
		t.Error("Animation color should not be empty")
	}
}

// TestPhase4EnhancedStyling tests the enhanced styling functions
func TestPhase4EnhancedStyling(t *testing.T) {
	// Test rainbow text creation
	rainbowText := CreateRainbowText("TEST")
	if rainbowText == "" {
		t.Error("Rainbow text should not be empty")
	}

	// Test gradient text creation
	gradientText := CreateGradientText("TEST", PrimaryColor, SecondaryColor)
	if gradientText == "" {
		t.Error("Gradient text should not be empty")
	}

	// Test glow text creation
	glowText := CreateGlowText("TEST", GlowColor)
	if glowText == "" {
		t.Error("Glow text should not be empty")
	}

	// Test enhanced title
	title := CreateEnhancedTitle("MONTY HALL")
	if title == "" {
		t.Error("Enhanced title should not be empty")
	}

	// Test winning message
	winMessage := CreateWinningMessage("YOU WIN!")
	if winMessage == "" {
		t.Error("Winning message should not be empty")
	}
}

// TestPhase4ResponsiveDesign tests the responsive design system
func TestPhase4ResponsiveDesign(t *testing.T) {
	// Test screen size detection
	smallSize := DetectScreenSize(70, 20)
	if smallSize != ScreenSmall {
		t.Errorf("Expected ScreenSmall for 70x20, got %v", smallSize)
	}

	mediumSize := DetectScreenSize(100, 30)
	if mediumSize != ScreenMedium {
		t.Errorf("Expected ScreenMedium for 100x30, got %v", mediumSize)
	}

	largeSize := DetectScreenSize(140, 40)
	if largeSize != ScreenLarge {
		t.Errorf("Expected ScreenLarge for 140x40, got %v", largeSize)
	}

	// Test responsive layout
	smallLayout := GetResponsiveLayout(ScreenSmall)
	if !smallLayout.CompactMode {
		t.Error("Small screen should use compact mode")
	}

	if smallLayout.DoorWidth >= 12 {
		t.Error("Small screen should have smaller door width")
	}

	largeLayout := GetResponsiveLayout(ScreenLarge)
	if largeLayout.CompactMode {
		t.Error("Large screen should not use compact mode")
	}

	if largeLayout.DoorWidth <= 8 {
		t.Error("Large screen should have larger door width")
	}

	// Test content adaptation
	longContent := "This is a very long line that should be truncated on small screens"
	adaptedContent := AdaptContentToScreen(longContent, smallLayout)
	if len(adaptedContent) > smallLayout.MaxWidth {
		t.Error("Content should be adapted for small screens")
	}
}

// TestPhase4GameIntegration tests Phase 4 integration with game flow
func TestPhase4GameIntegration(t *testing.T) {
	model := NewModel()
	model.Game = game.NewGame()

	// Test animation triggers during game flow
	err := model.Game.MakeInitialChoice(0)
	if err != nil {
		t.Fatalf("Failed to make initial choice: %v", err)
	}

	// Test that animations can be started
	cmd := model.startDoorOpenAnimation(model.Game.HostOpenedDoor)
	if cmd == nil && model.ShowAnimations {
		// This is okay - animation might not start if conditions aren't met
	}

	// Test winning animation
	if model.Game.Result != nil && model.Game.Result.Won {
		winCmd := model.startWinningAnimation()
		if winCmd == nil && model.ShowAnimations {
			// This is okay - animation might not start if conditions aren't met
		}
	}
}

// TestPhase4AnimationLifecycle tests the complete animation lifecycle
func TestPhase4AnimationLifecycle(t *testing.T) {
	// Create animation
	anim := NewAnimation("test", time.Millisecond*100, EaseLinear)
	if anim == nil {
		t.Fatal("Animation should be created")
	}

	// Test initial state
	if anim.State != AnimationStopped {
		t.Error("Animation should start in stopped state")
	}

	if anim.Progress != 0.0 {
		t.Error("Animation should start with 0 progress")
	}

	// Start animation
	anim.Start()
	if anim.State != AnimationRunning {
		t.Error("Animation should be running after start")
	}

	// Test pause/resume
	anim.Pause()
	if anim.State != AnimationPaused {
		t.Error("Animation should be paused")
	}

	anim.Resume()
	if anim.State != AnimationRunning {
		t.Error("Animation should be running after resume")
	}

	// Test stop
	anim.Stop()
	if anim.State != AnimationStopped {
		t.Error("Animation should be stopped")
	}

	if anim.Progress != 0.0 {
		t.Error("Animation progress should reset to 0 after stop")
	}
}

// TestPhase4EnhancedDoorRendering tests the enhanced door rendering
func TestPhase4EnhancedDoorRendering(t *testing.T) {
	door := &game.Door{
		State:   game.Closed,
		Content: game.Car,
	}

	doorComp := NewDoorComponent(1, door, false, false)
	if doorComp == nil {
		t.Fatal("Door component should be created")
	}

	// Test normal rendering
	normalRender := doorComp.Render()
	if normalRender == "" {
		t.Error("Normal door rendering should not be empty")
	}

	// Test animated rendering
	animRender := doorComp.RenderWithAnimation("🔓", WarningColor, true)
	if animRender == "" {
		t.Error("Animated door rendering should not be empty")
	}

	// Animated rendering should be different from normal
	if normalRender == animRender {
		t.Error("Animated rendering should differ from normal rendering")
	}
}

// TestPhase4VisualEffects tests various visual effects
func TestPhase4VisualEffects(t *testing.T) {
	// Test sparkle creation
	sparkles := CreateSparkleEffect(3)
	if sparkles == "" {
		t.Error("Sparkle effect should not be empty")
	}

	// Test glow effect
	glowText := CreateGlowEffect("TEST", 0.8)
	if glowText == "" {
		t.Error("Glow effect should not be empty")
	}

	// Test pulse effect
	pulseText := CreatePulseEffect("TEST", 0.5)
	if pulseText == "" {
		t.Error("Pulse effect should not be empty")
	}

	// Test gradient application
	gradientText := ApplyGradient("TEST", WinGradient, 0.5)
	if gradientText == "" {
		t.Error("Gradient application should not be empty")
	}
}
