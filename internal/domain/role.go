package domain

const (
	RolePriest          = "priest"
	RoleAssistantPriest = "assistant_priest"
	RoleDeacon          = "deacon"
	RoleAssistantDeacon = "assistant_deacon"
	RoleCongregation    = "congregation"
	RoleChanter         = "chanter"
	RoleReader          = "reader"
	RoleRubric          = "rubric"
	RoleMixed           = "mixed"
)

func IsValidRole(role string) bool {
	switch role {
	case RolePriest,
		RoleAssistantPriest,
		RoleDeacon,
		RoleAssistantDeacon,
		RoleCongregation,
		RoleChanter,
		RoleReader,
		RoleRubric,
		RoleMixed:
		return true
	default:
		return false
	}
}
