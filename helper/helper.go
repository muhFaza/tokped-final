package helper

import "strconv"

func FormatIDR(n int) string {
	in := strconv.Itoa(n)
	var out []rune
	l := len(in)
	for i, r := range in {
		out = append(out, r)
		if (l-i)%3 == 1 && i < l-1 {
			out = append(out, ',')
		}
	}
	return "Rp. " + string(out)
}
