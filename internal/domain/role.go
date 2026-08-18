package domain

const (
	RolePriest          = "priest"
	RoleAssistantPriest = "assistant_priest"
	RoleDeacon          = "deacon"
	RoleCongregation    = "congregation"
	RoleChanter         = "chanter"
	RoleReader          = "reader"
	RoleRubric          = "rubric"
)

func IsValidRole(role string) bool {
	switch role {
	case RolePriest,
		RoleAssistantPriest,
		RoleDeacon,
		RoleCongregation,
		RoleChanter,
		RoleReader,
		RoleRubric:
		return true
	default:
		return false
	}
}
