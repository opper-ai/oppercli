package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// truncateString shortens a string to maxLen characters, adding "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ConfirmDeletion asks the user to confirm deletion unless force is true
func ConfirmDeletion(resourceType, name string, force bool) (bool, error) {
	if force {
		return true, nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete %s '%s'? This cannot be undone. [y/N]: ", resourceType, name)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("error reading confirmation: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}
