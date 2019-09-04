package types

// RuleUsernameLength - light weight username check.
func RuleUsernameLength(username AccountKey) bool {
	if len(username) < MinimumUsernameLength ||
		len(username) > MaximumUsernameLength {
		return false
	}
	return true
}
