package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorGeneric ErrorType = iota
	ErrorFilePermission
	ErrorDiskSpace
	ErrorInvalidInput
	ErrorNetwork
	ErrorSystem
	ErrorConfig
	ErrorStats
)

// EnhancedError represents an error with recovery suggestions
type EnhancedError struct {
	Type        ErrorType
	Message     string
	Cause       error
	Suggestions []string
	Context     map[string]string
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	return e.Message
}

// GetDisplayMessage returns a formatted error message with recovery suggestions
func (e *EnhancedError) GetDisplayMessage() string {
	var builder strings.Builder

	// Main error message
	builder.WriteString("❌ ")
	builder.WriteString(e.Message)

	// Add suggestions if available
	if len(e.Suggestions) > 0 {
		builder.WriteString("\n\n💡 Suggestions:")
		for _, suggestion := range e.Suggestions {
			builder.WriteString("\n  • ")
			builder.WriteString(suggestion)
		}
	}

	// Add context if available
	if len(e.Context) > 0 {
		builder.WriteString("\n\n📋 Details:")
		for key, value := range e.Context {
			builder.WriteString(fmt.Sprintf("\n  %s: %s", key, value))
		}
	}

	return builder.String()
}

// EnhanceError converts a regular error into an enhanced error with recovery suggestions
func EnhanceError(err error, context ...string) *EnhancedError {
	if err == nil {
		return nil
	}

	// Check if it's already an enhanced error
	if enhanced, ok := err.(*EnhancedError); ok {
		return enhanced
	}

	errorMsg := err.Error()
	errorType := classifyError(err)

	enhanced := &EnhancedError{
		Type:        errorType,
		Message:     errorMsg,
		Cause:       err,
		Suggestions: generateSuggestions(errorType, err),
		Context:     make(map[string]string),
	}

	// Add context information
	if len(context) > 0 {
		enhanced.Context["Operation"] = context[0]
	}
	if len(context) > 1 {
		enhanced.Context["File"] = context[1]
	}

	return enhanced
}

// classifyError determines the type of error based on its characteristics
func classifyError(err error) ErrorType {
	errorMsg := strings.ToLower(err.Error())

	// File permission errors
	if strings.Contains(errorMsg, "permission denied") ||
		strings.Contains(errorMsg, "access denied") ||
		strings.Contains(errorMsg, "operation not permitted") {
		return ErrorFilePermission
	}

	// Disk space errors
	if strings.Contains(errorMsg, "no space left") ||
		strings.Contains(errorMsg, "disk full") ||
		strings.Contains(errorMsg, "not enough space") {
		return ErrorDiskSpace
	}

	// Network errors
	if strings.Contains(errorMsg, "connection refused") ||
		strings.Contains(errorMsg, "network unreachable") ||
		strings.Contains(errorMsg, "timeout") ||
		strings.Contains(errorMsg, "dns") {
		return ErrorNetwork
	}

	// Invalid input errors
	if strings.Contains(errorMsg, "invalid") ||
		strings.Contains(errorMsg, "malformed") ||
		strings.Contains(errorMsg, "parse") ||
		strings.Contains(errorMsg, "unmarshal") {
		return ErrorInvalidInput
	}

	// Configuration errors
	if strings.Contains(errorMsg, "config") ||
		strings.Contains(errorMsg, "configuration") {
		return ErrorConfig
	}

	// Statistics errors
	if strings.Contains(errorMsg, "stats") ||
		strings.Contains(errorMsg, "statistics") {
		return ErrorStats
	}

	return ErrorGeneric
}

// generateSuggestions creates recovery suggestions based on error type
func generateSuggestions(errorType ErrorType, err error) []string {
	switch errorType {
	case ErrorFilePermission:
		return generateFilePermissionSuggestions(err)
	case ErrorDiskSpace:
		return generateDiskSpaceSuggestions()
	case ErrorInvalidInput:
		return generateInvalidInputSuggestions(err)
	case ErrorNetwork:
		return generateNetworkSuggestions()
	case ErrorSystem:
		return generateSystemSuggestions()
	case ErrorConfig:
		return generateConfigSuggestions()
	case ErrorStats:
		return generateStatsSuggestions()
	default:
		return generateGenericSuggestions()
	}
}

// generateFilePermissionSuggestions creates suggestions for file permission errors
func generateFilePermissionSuggestions(err error) []string {
	suggestions := []string{
		"Check if you have read/write permissions for the file or directory",
	}

	if runtime.GOOS != "windows" {
		suggestions = append(suggestions,
			"Try running: chmod 644 <filename> (for files)",
			"Try running: chmod 755 <directory> (for directories)",
			"Check if the file is owned by another user: ls -la <filename>",
		)
	} else {
		suggestions = append(suggestions,
			"Right-click the file/folder → Properties → Security → Edit permissions",
			"Try running the application as Administrator",
		)
	}

	suggestions = append(suggestions,
		"Ensure the parent directory exists and is writable",
		"Check if the file is currently open in another application",
	)

	return suggestions
}

// generateDiskSpaceSuggestions creates suggestions for disk space errors
func generateDiskSpaceSuggestions() []string {
	suggestions := []string{
		"Free up disk space by deleting unnecessary files",
		"Empty your trash/recycle bin",
		"Clear temporary files and caches",
	}

	if runtime.GOOS != "windows" {
		suggestions = append(suggestions,
			"Check disk usage: df -h",
			"Find large files: du -sh * | sort -hr",
			"Clear system logs: sudo journalctl --vacuum-time=7d",
		)
	} else {
		suggestions = append(suggestions,
			"Run Disk Cleanup utility",
			"Check disk space in File Explorer",
			"Consider moving files to external storage",
		)
	}

	suggestions = append(suggestions,
		"Try saving to a different location with more space",
		"Consider compressing old files",
	)

	return suggestions
}

// generateInvalidInputSuggestions creates suggestions for invalid input errors
func generateInvalidInputSuggestions(err error) []string {
	suggestions := []string{
		"Check the input format and try again",
		"Ensure all required fields are filled",
	}

	errorMsg := strings.ToLower(err.Error())

	if strings.Contains(errorMsg, "json") {
		suggestions = append(suggestions,
			"Verify JSON syntax is correct",
			"Check for missing commas or brackets",
			"Use a JSON validator to check the format",
		)
	}

	if strings.Contains(errorMsg, "number") || strings.Contains(errorMsg, "integer") {
		suggestions = append(suggestions,
			"Enter numbers only (1-9 for confirmation)",
			"Avoid letters, spaces, or special characters",
		)
	}

	if strings.Contains(errorMsg, "strategy") {
		suggestions = append(suggestions,
			"Valid strategies: 'switch', 'stay', or 'ask'",
			"Check spelling and use lowercase",
		)
	}

	if strings.Contains(errorMsg, "color") {
		suggestions = append(suggestions,
			"Valid color schemes: 'default', 'high-contrast', 'colorblind-safe'",
		)
	}

	return suggestions
}

// generateNetworkSuggestions creates suggestions for network errors
func generateNetworkSuggestions() []string {
	return []string{
		"Check your internet connection",
		"Verify the server is running and accessible",
		"Try again in a few moments",
		"Check if a firewall is blocking the connection",
		"Verify the URL or address is correct",
		"Contact your network administrator if the problem persists",
	}
}

// generateSystemSuggestions creates suggestions for system errors
func generateSystemSuggestions() []string {
	return []string{
		"Try restarting the application",
		"Check system resources (CPU, memory)",
		"Ensure your system meets the minimum requirements",
		"Update your operating system",
		"Check system logs for more details",
		"Contact support if the problem persists",
	}
}

// generateConfigSuggestions creates suggestions for configuration errors
func generateConfigSuggestions() []string {
	return []string{
		"Check the configuration file syntax",
		"Restore default configuration if corrupted",
		"Verify all configuration values are valid",
		"Check file permissions for the config directory",
		"Try deleting the config file to regenerate defaults",
		"Refer to the documentation for valid configuration options",
	}
}

// generateStatsSuggestions creates suggestions for statistics errors
func generateStatsSuggestions() []string {
	return []string{
		"Check if the statistics file is corrupted",
		"Try resetting statistics if the file is damaged",
		"Ensure sufficient disk space for statistics storage",
		"Verify write permissions for the statistics directory",
		"Consider backing up statistics before making changes",
		"Check if another instance of the application is running",
	}
}

// generateGenericSuggestions creates general suggestions for unclassified errors
func generateGenericSuggestions() []string {
	return []string{
		"Try the operation again",
		"Restart the application if the problem persists",
		"Check if you have sufficient permissions",
		"Ensure all required files are present",
		"Contact support if you continue to experience issues",
	}
}

// CreateFilePermissionError creates a specific file permission error
func CreateFilePermissionError(operation, filename string, cause error) *EnhancedError {
	return &EnhancedError{
		Type:        ErrorFilePermission,
		Message:     fmt.Sprintf("Permission denied: %s", operation),
		Cause:       cause,
		Suggestions: generateFilePermissionSuggestions(cause),
		Context: map[string]string{
			"Operation": operation,
			"File":      filename,
			"OS":        runtime.GOOS,
		},
	}
}

// CreateDiskSpaceError creates a specific disk space error
func CreateDiskSpaceError(operation string, cause error) *EnhancedError {
	return &EnhancedError{
		Type:        ErrorDiskSpace,
		Message:     "Insufficient disk space",
		Cause:       cause,
		Suggestions: generateDiskSpaceSuggestions(),
		Context: map[string]string{
			"Operation": operation,
		},
	}
}

// CreateInvalidInputError creates a specific invalid input error
func CreateInvalidInputError(input, expected string) *EnhancedError {
	return &EnhancedError{
		Type:    ErrorInvalidInput,
		Message: fmt.Sprintf("Invalid input: expected %s", expected),
		Suggestions: []string{
			fmt.Sprintf("Valid options: %s", expected),
			"Check spelling and format",
			"Refer to help documentation for examples",
		},
		Context: map[string]string{
			"Input":    input,
			"Expected": expected,
		},
	}
}

// CreateConfigError creates a specific configuration error
func CreateConfigError(setting string, cause error) *EnhancedError {
	return &EnhancedError{
		Type:        ErrorConfig,
		Message:     fmt.Sprintf("Configuration error: %s", setting),
		Cause:       cause,
		Suggestions: generateConfigSuggestions(),
		Context: map[string]string{
			"Setting": setting,
		},
	}
}

// WrapError wraps an existing error with enhanced error handling
func WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for specific error types
	if os.IsPermission(err) {
		return CreateFilePermissionError(operation, "", err)
	}

	if isNoSpaceError(err) {
		return CreateDiskSpaceError(operation, err)
	}

	// Default to enhanced error
	return EnhanceError(err, operation)
}

// isNoSpaceError checks if an error is related to disk space
func isNoSpaceError(err error) bool {
	if err == nil {
		return false
	}

	errorMsg := strings.ToLower(err.Error())
	return strings.Contains(errorMsg, "no space left") ||
		strings.Contains(errorMsg, "disk full") ||
		strings.Contains(errorMsg, "not enough space")
}

// FormatErrorForDisplay formats an error for display in the UI
func FormatErrorForDisplay(err error) string {
	if err == nil {
		return ""
	}

	if enhanced, ok := err.(*EnhancedError); ok {
		return enhanced.GetDisplayMessage()
	}

	// For regular errors, enhance them first
	enhanced := EnhanceError(err)
	return enhanced.GetDisplayMessage()
}
