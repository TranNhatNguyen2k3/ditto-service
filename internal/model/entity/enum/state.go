package enum

const (
	ActionDelete = "DELETE"
	ActionUpdate = "UPDATE"
	ActionCreate = "CREATE"
)

const (
	StateActive   = "ACTIVE"
	StateInactive = "INACTIVE"
)

// Convert state from int to string
func EnumState(state int) string {
	switch state {
	case 0:
		return StateActive
	case 1:
		return StateInactive
	default:
		return ""
	}
}

// Convert state from string to int
func EnumStateDB(state string) int {
	switch state {
	case StateActive:
		return 0
	case StateInactive:
		return 1
	default:
		return -1
	}
}
