package fizzbuzz

import "strconv"

// Generate returns strings from 1..limit where multiples
// of int1 are replaced by str1, multiples of int2 by str2,
// and multiples of both by str1+str2.
func Generate(int1, int2, limit int, str1, str2 string) []string {
	out := make([]string, 0, limit)

	for i := 1; i <= limit; i++ {
		s := ""

		if i%int1 == 0 {
			s += str1
		}
		if i%int2 == 0 {
			s += str2
		}

		if s == "" {
			s = strconv.Itoa(i)
		}

		out = append(out, s)
	}

	return out
}
