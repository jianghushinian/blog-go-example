package store

const (
	defaultLimitValue = 10
	defaultOrderValue = "id desc"
)

func defaultLimit(limit int) int {
	if limit == 0 {
		limit = defaultLimitValue
	}
	return limit
}

func defaultOrder(order string) string {
	if order == "" {
		order = defaultOrderValue
	}
	return order
}
