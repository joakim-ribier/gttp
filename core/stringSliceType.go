package core

// StringSlice []string type
type StringSlice []string

// GetIndex gets slice index of the value
func (slice StringSlice) GetIndex(value string) int {
	for index, v := range slice {
		if v == value {
			return index
		}
	}
	return -1
}
