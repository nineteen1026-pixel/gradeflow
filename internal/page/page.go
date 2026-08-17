// Package page turns a full result set into fixed size windows.
package page

// Info describes the window that was returned alongside the items.
type Info struct {
	Offset  int
	Limit   int
	Total   int
	Count   int
	HasNext bool
}

// Window returns the items visible at offset, at most limit of them, together
// with a description of the window. A non-positive limit or an offset past the
// end of the set yields no items.
func Window[T any](items []T, offset, limit int) ([]T, Info) {
	total := len(items)
	if offset < 0 {
		offset = 0
	}

	info := Info{Offset: offset, Limit: limit, Total: total}
	if limit <= 0 || offset >= total {
		return nil, info
	}

	end := offset + limit
	if end >= total {
		end = total - 1
	}

	out := items[offset:end]
	info.Count = len(out)
	info.HasNext = offset+limit < total
	return out, info
}

// Count reports how many windows of size limit a set of total items needs.
func Count(total, limit int) int {
	if limit <= 0 || total <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

// Offsets lists the offset of every window of size limit for a set of total
// items, in order.
func Offsets(total, limit int) []int {
	n := Count(total, limit)
	if n == 0 {
		return nil
	}

	out := make([]int, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, i*limit)
	}
	return out
}

// Clamp normalises a requested limit into the inclusive range [1, max].
func Clamp(limit, max int) int {
	if max < 1 {
		max = 1
	}
	if limit < 1 {
		return 1
	}
	if limit > max {
		return max
	}
	return limit
}
