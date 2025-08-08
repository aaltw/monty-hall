package ui

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestEnhanceError(t *testing.T) {
	// Test nil error
	enhanced := EnhanceError(nil)
	if enhanced != nil {
		t.Error("EnhanceError should return nil for nil input")
	}

	// Test regular error
	originalErr := errors.New("test error")
	enhanced = EnhanceError(originalErr)

	if enhanced == nil {
		t.Fatal("EnhanceError should not return nil for valid error")
	}

	if enhanced.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", enhanced.Message)
	}

	if enhanced.Cause != originalErr {
		t.Error("Enhanced error should preserve original cause")
	}

	// Test already enhanced error
	enhanced2 := EnhanceError(enhanced)
	if enhanced2 != enhanced {
		t.Error("EnhanceError should return the same enhanced error if already enhanced")
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		errorMsg string
		expected ErrorType
	}{
		{"permission denied", ErrorFilePermission},
		{"access denied", ErrorFilePermission},
		{"operation not permitted", ErrorFilePermission},
		{"no space left on device", ErrorDiskSpace},
		{"disk full", ErrorDiskSpace},
		{"invalid input", ErrorInvalidInput},
		{"malformed json", ErrorInvalidInput},
		{"connection refused", ErrorNetwork},
		{"network unreachable", ErrorNetwork},
		{"config error", ErrorConfig},
		{"statistics failed", ErrorStats},
		{"unknown error", ErrorGeneric},
	}

	for _, tt := range tests {
		t.Run(tt.errorMsg, func(t *testing.T) {
			err := errors.New(tt.errorMsg)
			errorType := classifyError(err)
			if errorType != tt.expected {
				t.Errorf("Expected error type %v, got %v", tt.expected, errorType)
			}
		})
	}
}

func TestGenerateFilePermissionSuggestions(t *testing.T) {
	err := errors.New("permission denied")
	suggestions := generateFilePermissionSuggestions(err)

	if len(suggestions) == 0 {
		t.Error("Should generate suggestions for file permission errors")
	}

	// Check for platform-specific suggestions
	hasChmod := false
	hasWindowsSpecific := false

	for _, suggestion := range suggestions {
		if strings.Contains(suggestion, "chmod") {
			hasChmod = true
		}
		if strings.Contains(suggestion, "Administrator") || strings.Contains(suggestion, "Properties") {
			hasWindowsSpecific = true
		}
	}

	if runtime.GOOS != "windows" && !hasChmod {
		t.Error("Should include chmod suggestions on Unix-like systems")
	}

	if runtime.GOOS == "windows" && !hasWindowsSpecific {
		t.Error("Should include Windows-specific suggestions on Windows")
	}
}

func TestGenerateDiskSpaceSuggestions(t *testing.T) {
	suggestions := generateDiskSpaceSuggestions()

	if len(suggestions) == 0 {
		t.Error("Should generate suggestions for disk space errors")
	}

	// Check for platform-specific suggestions
	hasUnixCommands := false
	hasWindowsSpecific := false

	for _, suggestion := range suggestions {
		if strings.Contains(suggestion, "df -h") || strings.Contains(suggestion, "du -sh") {
			hasUnixCommands = true
		}
		if strings.Contains(suggestion, "Disk Cleanup") {
			hasWindowsSpecific = true
		}
	}

	if runtime.GOOS != "windows" && !hasUnixCommands {
		t.Error("Should include Unix commands on Unix-like systems")
	}

	if runtime.GOOS == "windows" && !hasWindowsSpecific {
		t.Error("Should include Windows-specific suggestions on Windows")
	}
}

func TestGenerateInvalidInputSuggestions(t *testing.T) {
	tests := []struct {
		errorMsg     string
		expectedText string
	}{
		{"invalid json format", "JSON"},
		{"invalid number", "numbers only"},
		{"invalid strategy", "switch"},
		{"invalid color scheme", "color schemes"},
	}

	for _, tt := range tests {
		t.Run(tt.errorMsg, func(t *testing.T) {
			err := errors.New(tt.errorMsg)
			suggestions := generateInvalidInputSuggestions(err)

			if len(suggestions) == 0 {
				t.Error("Should generate suggestions for invalid input errors")
			}

			found := false
			for _, suggestion := range suggestions {
				if strings.Contains(strings.ToLower(suggestion), strings.ToLower(tt.expectedText)) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Should include suggestion about %s", tt.expectedText)
			}
		})
	}
}

func TestCreateSpecificErrors(t *testing.T) {
	// Test file permission error
	permErr := CreateFilePermissionError("write file", "test.txt", errors.New("permission denied"))
	if permErr.Type != ErrorFilePermission {
		t.Error("Should create file permission error type")
	}
	if permErr.Context["File"] != "test.txt" {
		t.Error("Should include filename in context")
	}

	// Test disk space error
	diskErr := CreateDiskSpaceError("save data", errors.New("no space left"))
	if diskErr.Type != ErrorDiskSpace {
		t.Error("Should create disk space error type")
	}

	// Test invalid input error
	inputErr := CreateInvalidInputError("abc", "numbers 1-9")
	if inputErr.Type != ErrorInvalidInput {
		t.Error("Should create invalid input error type")
	}
	if inputErr.Context["Expected"] != "numbers 1-9" {
		t.Error("Should include expected input in context")
	}

	// Test config error
	configErr := CreateConfigError("color_scheme", errors.New("invalid value"))
	if configErr.Type != ErrorConfig {
		t.Error("Should create config error type")
	}
}

func TestWrapError(t *testing.T) {
	// Test nil error
	wrapped := WrapError(nil, "test operation")
	if wrapped != nil {
		t.Error("WrapError should return nil for nil input")
	}

	// Test permission error
	permErr := &os.PathError{Op: "open", Path: "test.txt", Err: os.ErrPermission}
	wrapped = WrapError(permErr, "read file")

	enhanced, ok := wrapped.(*EnhancedError)
	if !ok {
		t.Fatal("WrapError should return EnhancedError")
	}

	if enhanced.Type != ErrorFilePermission {
		t.Error("Should detect permission error type")
	}

	// Test generic error
	genericErr := errors.New("generic error")
	wrapped = WrapError(genericErr, "test operation")

	enhanced, ok = wrapped.(*EnhancedError)
	if !ok {
		t.Fatal("WrapError should return EnhancedError for generic errors")
	}
}

func TestEnhancedErrorMethods(t *testing.T) {
	err := &EnhancedError{
		Type:    ErrorFilePermission,
		Message: "Test error",
		Suggestions: []string{
			"Try this",
			"Or try that",
		},
		Context: map[string]string{
			"File": "test.txt",
			"Op":   "read",
		},
	}

	// Test Error() method
	if err.Error() != "Test error" {
		t.Errorf("Expected 'Test error', got '%s'", err.Error())
	}

	// Test GetDisplayMessage() method
	display := err.GetDisplayMessage()

	if !strings.Contains(display, "❌ Test error") {
		t.Error("Display message should contain error icon and message")
	}

	if !strings.Contains(display, "💡 Suggestions:") {
		t.Error("Display message should contain suggestions section")
	}

	if !strings.Contains(display, "Try this") {
		t.Error("Display message should contain suggestions")
	}

	if !strings.Contains(display, "📋 Details:") {
		t.Error("Display message should contain details section")
	}

	if !strings.Contains(display, "File: test.txt") {
		t.Error("Display message should contain context information")
	}
}

func TestFormatErrorForDisplay(t *testing.T) {
	// Test nil error
	formatted := FormatErrorForDisplay(nil)
	if formatted != "" {
		t.Error("Should return empty string for nil error")
	}

	// Test enhanced error
	enhanced := &EnhancedError{
		Message:     "Test enhanced error",
		Suggestions: []string{"Test suggestion"},
	}

	formatted = FormatErrorForDisplay(enhanced)
	if !strings.Contains(formatted, "Test enhanced error") {
		t.Error("Should format enhanced error correctly")
	}

	// Test regular error
	regular := errors.New("regular error")
	formatted = FormatErrorForDisplay(regular)

	if !strings.Contains(formatted, "regular error") {
		t.Error("Should enhance and format regular error")
	}

	if !strings.Contains(formatted, "💡 Suggestions:") {
		t.Error("Should include suggestions for regular error")
	}
}

func TestIsNoSpaceError(t *testing.T) {
	tests := []struct {
		errorMsg string
		expected bool
	}{
		{"no space left on device", true},
		{"disk full", true},
		{"not enough space", true},
		{"permission denied", false},
		{"file not found", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.errorMsg, func(t *testing.T) {
			var err error
			if tt.errorMsg != "" {
				err = errors.New(tt.errorMsg)
			}

			result := isNoSpaceError(err)
			if result != tt.expected {
				t.Errorf("Expected %v for '%s', got %v", tt.expected, tt.errorMsg, result)
			}
		})
	}

	// Test nil error
	if isNoSpaceError(nil) {
		t.Error("Should return false for nil error")
	}
}

// Benchmark tests
func BenchmarkEnhanceError(b *testing.B) {
	err := errors.New("test error")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		EnhanceError(err)
	}
}

func BenchmarkClassifyError(b *testing.B) {
	err := errors.New("permission denied")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		classifyError(err)
	}
}

func BenchmarkFormatErrorForDisplay(b *testing.B) {
	err := errors.New("test error")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FormatErrorForDisplay(err)
	}
}
