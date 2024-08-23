package network

import (
	. "Project/dataenums"
)
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// Helper function to check if HRAInput is empty
func isEmptyHRAInput(input HRAInput) bool {
	return len(input.HallRequests) == 0 && len(input.States) == 0
}
