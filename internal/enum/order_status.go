package enum

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
)

func (s Status) IsValid() bool {
	return s == StatusPending || s == StatusConfirmed || s == StatusShipped || s == StatusDelivered
}
