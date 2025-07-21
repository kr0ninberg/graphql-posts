package helpers

func Paginate[T any](items []*T, limit *int32, offset *int32) []*T {
	start := 0
	end := len(items)

	if offset != nil && int(*offset) < len(items) {
		start = int(*offset)
	}
	if limit != nil && start+int(*limit) < len(items) {
		end = start + int(*limit)
	}

	if start > end {
		start = end
	}

	return items[start:end]
}
