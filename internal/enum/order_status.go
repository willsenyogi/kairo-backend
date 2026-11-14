package enum

type Status int

const (
	StatusPending   Status = 1
	StatusConfirmed Status = 2
	StatusShipped   Status = 3
	StatusDelivered Status = 4
)

func (s Status) IsValid() bool {
	return s == StatusPending || s == StatusConfirmed || s == StatusShipped || s == StatusDelivered
}
