package ui

import (
	"math"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/westhuis/monty-hall/pkg/game"
)

// EasingFunction defines the signature for easing functions
type EasingFunction func(t float64) float64

// Common easing functions
var (
	EaseLinear = func(t float64) float64 { return t }

	EaseInOut = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	EaseBounce = func(t float64) float64 {
		if t < 1/2.75 {
			return 7.5625 * t * t
		} else if t < 2/2.75 {
			t -= 1.5 / 2.75
			return 7.5625*t*t + 0.75
		} else if t < 2.5/2.75 {
			t -= 2.25 / 2.75
			return 7.5625*t*t + 0.9375
		} else {
			t -= 2.625 / 2.75
			return 7.5625*t*t + 0.984375
		}
	}

	EaseElastic = func(t float64) float64 {
		if t == 0 || t == 1 {
			return t
		}
		p := 0.3
		s := p / 4
		return math.Pow(2, -10*t)*math.Sin((t-s)*(2*math.Pi)/p) + 1
	}
)

// AnimationState represents the current state of an animation
type AnimationState int

const (
	AnimationStopped AnimationState = iota
	AnimationRunning
	AnimationPaused
	AnimationComplete
)

// Animation represents a single animation with timing and easing
type Animation struct {
	ID         string
	Duration   time.Duration
	StartTime  time.Time
	Easing     EasingFunction
	State      AnimationState
	Progress   float64
	OnUpdate   func(progress float64)
	OnComplete func()
	Loop       bool
	Reverse    bool
}

// NewAnimation creates a new animation with the given parameters
func NewAnimation(id string, duration time.Duration, easing EasingFunction) *Animation {
	return &Animation{
		ID:       id,
		Duration: duration,
		Easing:   easing,
		State:    AnimationStopped,
		Progress: 0.0,
		Loop:     false,
		Reverse:  false,
	}
}

// Start begins the animation
func (a *Animation) Start() {
	a.State = AnimationRunning
	a.StartTime = time.Now()
	a.Progress = 0.0
}

// Pause pauses the animation
func (a *Animation) Pause() {
	if a.State == AnimationRunning {
		a.State = AnimationPaused
	}
}

// Resume resumes a paused animation
func (a *Animation) Resume() {
	if a.State == AnimationPaused {
		a.State = AnimationRunning
		// Adjust start time to account for pause duration
		elapsed := time.Duration(a.Progress * float64(a.Duration))
		a.StartTime = time.Now().Add(-elapsed)
	}
}

// Stop stops the animation
func (a *Animation) Stop() {
	a.State = AnimationStopped
	a.Progress = 0.0
}

// Update updates the animation progress based on current time
func (a *Animation) Update() bool {
	if a.State != AnimationRunning {
		return false
	}

	elapsed := time.Since(a.StartTime)
	rawProgress := float64(elapsed) / float64(a.Duration)

	if rawProgress >= 1.0 {
		a.Progress = 1.0
		a.State = AnimationComplete

		if a.OnComplete != nil {
			a.OnComplete()
		}

		if a.Loop {
			a.Start() // Restart the animation
			return true
		}

		return false // Animation finished
	}

	// Apply easing function
	if a.Easing != nil {
		a.Progress = a.Easing(rawProgress)
	} else {
		a.Progress = rawProgress
	}

	// Handle reverse animation
	if a.Reverse {
		a.Progress = 1.0 - a.Progress
	}

	// Call update callback
	if a.OnUpdate != nil {
		a.OnUpdate(a.Progress)
	}

	return true // Animation continues
}

// IsRunning returns true if the animation is currently running
func (a *Animation) IsRunning() bool {
	return a.State == AnimationRunning
}

// IsComplete returns true if the animation has completed
func (a *Animation) IsComplete() bool {
	return a.State == AnimationComplete
}

// AnimationManager manages multiple animations
type AnimationManager struct {
	animations map[string]*Animation
	ticker     *time.Ticker
	running    bool
}

// NewAnimationManager creates a new animation manager
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: make(map[string]*Animation),
		running:    false,
	}
}

// AddAnimation adds an animation to the manager
func (am *AnimationManager) AddAnimation(animation *Animation) {
	am.animations[animation.ID] = animation
}

// RemoveAnimation removes an animation from the manager
func (am *AnimationManager) RemoveAnimation(id string) {
	delete(am.animations, id)
}

// GetAnimation retrieves an animation by ID
func (am *AnimationManager) GetAnimation(id string) *Animation {
	return am.animations[id]
}

// StartAnimation starts a specific animation
func (am *AnimationManager) StartAnimation(id string) {
	if anim, exists := am.animations[id]; exists {
		anim.Start()
		am.ensureRunning()
	}
}

// StopAnimation stops a specific animation
func (am *AnimationManager) StopAnimation(id string) {
	if anim, exists := am.animations[id]; exists {
		anim.Stop()
	}
}

// StopAllAnimations stops all running animations
func (am *AnimationManager) StopAllAnimations() {
	for _, anim := range am.animations {
		anim.Stop()
	}
	am.stop()
}

// Update updates all running animations
func (am *AnimationManager) Update() tea.Cmd {
	if !am.running {
		return nil
	}

	hasRunning := false
	for _, anim := range am.animations {
		if anim.Update() {
			hasRunning = true
		}
	}

	if !hasRunning {
		am.stop()
		return nil
	}

	// Return a command to trigger the next update
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return AnimationTickMsg{Time: t}
	})
}

// ensureRunning starts the animation loop if not already running
func (am *AnimationManager) ensureRunning() {
	if !am.running {
		am.running = true
	}
}

// stop stops the animation loop
func (am *AnimationManager) stop() {
	am.running = false
}

// HasRunningAnimations returns true if any animations are currently running
func (am *AnimationManager) HasRunningAnimations() bool {
	for _, anim := range am.animations {
		if anim.IsRunning() {
			return true
		}
	}
	return false
}

// AnimationTickMsg is sent to trigger animation updates
type AnimationTickMsg struct {
	Time time.Time
}

// DoorOpenAnimation represents a door opening animation
type DoorOpenAnimation struct {
	*Animation
	DoorIndex int
	Frames    []string
	Colors    []lipgloss.Color
}

// NewDoorOpenAnimation creates a new door opening animation
func NewDoorOpenAnimation(doorIndex int) *DoorOpenAnimation {
	frames := []string{
		"🚪", // Closed
		"🔓", // Unlocking
		"📂", // Opening
		"✨", // Revealing
	}

	colors := []lipgloss.Color{
		DoorColor,
		WarningColor,
		PrimaryColor,
		SecondaryColor,
	}

	anim := NewAnimation(
		"door_open_"+string(rune(doorIndex+'0')),
		time.Millisecond*800,
		EaseInOut,
	)

	return &DoorOpenAnimation{
		Animation: anim,
		DoorIndex: doorIndex,
		Frames:    frames,
		Colors:    colors,
	}
}

// GetCurrentFrame returns the current animation frame
func (doa *DoorOpenAnimation) GetCurrentFrame() (string, lipgloss.Color) {
	if len(doa.Frames) == 0 {
		return "🚪", DoorColor
	}

	frameIndex := int(doa.Progress * float64(len(doa.Frames)-1))
	if frameIndex >= len(doa.Frames) {
		frameIndex = len(doa.Frames) - 1
	}

	frame := doa.Frames[frameIndex]
	color := DoorColor
	if frameIndex < len(doa.Colors) {
		color = doa.Colors[frameIndex]
	}

	return frame, color
}

// PulseAnimation creates a pulsing effect
type PulseAnimation struct {
	*Animation
	BaseStyle  lipgloss.Style
	PulseColor lipgloss.Color
	Intensity  float64
}

// NewPulseAnimation creates a new pulse animation
func NewPulseAnimation(id string, baseStyle lipgloss.Style, pulseColor lipgloss.Color) *PulseAnimation {
	anim := NewAnimation(id, time.Millisecond*1000, EaseInOut)
	anim.Loop = true

	return &PulseAnimation{
		Animation:  anim,
		BaseStyle:  baseStyle,
		PulseColor: pulseColor,
		Intensity:  0.3,
	}
}

// GetCurrentStyle returns the current pulsed style
func (pa *PulseAnimation) GetCurrentStyle() lipgloss.Style {
	// Create a pulsing effect by interpolating between base and pulse colors
	intensity := pa.Intensity * math.Sin(pa.Progress*math.Pi*2)
	if intensity < 0 {
		intensity = 0
	}

	// For simplicity, just toggle bold on/off for pulse effect
	if intensity > 0.5 {
		return pa.BaseStyle.Copy().Bold(true)
	}
	return pa.BaseStyle
}

// TypewriterAnimation creates a typewriter text effect
type TypewriterAnimation struct {
	*Animation
	Text        string
	CurrentText string
}

// NewTypewriterAnimation creates a new typewriter animation
func NewTypewriterAnimation(id, text string, duration time.Duration) *TypewriterAnimation {
	anim := NewAnimation(id, duration, EaseLinear)

	ta := &TypewriterAnimation{
		Animation: anim,
		Text:      text,
	}

	anim.OnUpdate = func(progress float64) {
		length := int(progress * float64(len(ta.Text)))
		if length > len(ta.Text) {
			length = len(ta.Text)
		}
		ta.CurrentText = ta.Text[:length]
	}

	return ta
}

// GetCurrentText returns the current typewriter text
func (ta *TypewriterAnimation) GetCurrentText() string {
	return ta.CurrentText
}

// Particle system for visual effects
type Particle struct {
	X, Y    float64
	VX, VY  float64
	Life    float64
	MaxLife float64
	Char    string
	Color   lipgloss.Color
}

type ParticleSystem struct {
	particles []Particle
	width     int
	height    int
}

// NewParticleSystem creates a new particle system
func NewParticleSystem(width, height int) *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0),
		width:     width,
		height:    height,
	}
}

// AddWinningParticles adds celebration particles for winning
func (ps *ParticleSystem) AddWinningParticles(centerX, centerY int) {
	sparkles := []string{"✨", "⭐", "💫", "🌟", "✦", "✧", "🎉", "🎊"}
	colors := []lipgloss.Color{CarColor, SparkleColor, GlowColor, SecondaryColor}

	for i := 0; i < 20; i++ {
		particle := Particle{
			X:       float64(centerX),
			Y:       float64(centerY),
			VX:      (game.SecureFloat64() - 0.5) * 4,
			VY:      (game.SecureFloat64() - 0.5) * 4,
			Life:    1.0,
			MaxLife: 1.0,
			Char:    sparkles[game.SecureIntn(len(sparkles))],
			Color:   colors[game.SecureIntn(len(colors))]}
		ps.particles = append(ps.particles, particle)
	}
}

// Update updates all particles
func (ps *ParticleSystem) Update(deltaTime float64) {
	for i := len(ps.particles) - 1; i >= 0; i-- {
		p := &ps.particles[i]

		// Update position
		p.X += p.VX * deltaTime
		p.Y += p.VY * deltaTime

		// Update life
		p.Life -= deltaTime * 0.5

		// Remove dead particles
		if p.Life <= 0 {
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}
}

// Render renders all particles to a string grid
func (ps *ParticleSystem) Render() [][]string {
	grid := make([][]string, ps.height)
	for i := range grid {
		grid[i] = make([]string, ps.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for _, p := range ps.particles {
		x, y := int(p.X), int(p.Y)
		if x >= 0 && x < ps.width && y >= 0 && y < ps.height {
			style := lipgloss.NewStyle().Foreground(p.Color)
			grid[y][x] = style.Render(p.Char)
		}
	}

	return grid
}

// HasParticles returns true if there are active particles
func (ps *ParticleSystem) HasParticles() bool {
	return len(ps.particles) > 0
}

// GetParticleCount returns the number of active particles
func (ps *ParticleSystem) GetParticleCount() int {
	return len(ps.particles)
}
