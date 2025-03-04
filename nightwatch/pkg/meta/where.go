package meta

const (
	ListAll      = ""
	defaultLimit = 1000
	defaultOrder = "id desc"
)

type ListOption func(*ListOptions)

type ListOptions struct {
	Filters map[string]any
	Not     map[string]any
	Offset  int
	Limit   int
	Order   string
}

func NewListOptions(opts ...ListOption) ListOptions {
	los := ListOptions{
		Filters: map[string]any{},
		Offset:  0,
		Limit:   defaultLimit,
		Order:   defaultOrder,
	}

	for _, opt := range opts {
		opt(&los)
	}

	return los
}

func WithFilter(filter map[string]any) ListOption {
	return func(o *ListOptions) {
		o.Filters = filter
	}
}

func WithFilterNot(not map[string]any) ListOption {
	return func(o *ListOptions) {
		o.Not = not
	}
}

func WithOffset(offset int64) ListOption {
	return func(o *ListOptions) {
		if offset < 0 {
			offset = 0
		}
		o.Offset = int(offset)
	}
}

func WithLimit(limit int64) ListOption {
	return func(o *ListOptions) {
		if limit <= 0 {
			limit = defaultLimit
		}
		o.Limit = int(limit)
	}
}

func WithOrder(order string) ListOption {
	return func(o *ListOptions) {
		o.Order = order
	}
}
