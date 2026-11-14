package enum

type Role int

const (
	RoleBuyer  Role = 1
	RoleSeller Role = 2
	RoleAdmin  Role = 3
)

func (r Role) IsValid() bool {
	return r == RoleBuyer || r == RoleSeller || r == RoleAdmin
}
