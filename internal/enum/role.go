package enum

type Role string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
)

func (r Role) IsValid() bool {
	return r == RoleBuyer || r == RoleSeller
}