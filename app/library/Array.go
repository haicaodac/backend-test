package library

// ArrayStringContains ...
func ArrayStringContains(array []string, str string) bool {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return false
}

// ArrayIntContains ...
func ArrayIntContains(array []int, str int) bool {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return false
}

// ArrayRemoveInt ...
func ArrayRemoveInt(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}
