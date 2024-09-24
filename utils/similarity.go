package utils

import "strings"

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func LevenshteinDistance(str1, str2 string) int {
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)
	s1, s2 := []rune(str1), []rune(str2)
	lenS1, lenS2 := len(s1), len(s2)
	if lenS1 == 0 {
		return lenS2
	}
	if lenS2 == 0 {
		return lenS1
	}

	d := make([][]int, lenS1+1)
	for i := range d {
		d[i] = make([]int, lenS2+1)
	}

	for i := 0; i <= lenS1; i++ {
		d[i][0] = i
	}
	for j := 0; j <= lenS2; j++ {
		d[0][j] = j
	}

	for i := 1; i <= lenS1; i++ {
		for j := 1; j <= lenS2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	return d[lenS1][lenS2]
}

func Similarity(str1, str2 string) float64 {
	str1 = strings.ReplaceAll(str1, " ", "")
	str2 = strings.ReplaceAll(str2, " ", "")
	distance := LevenshteinDistance(str1, str2)
	maxLength := len(str1)
	if len(str2) > maxLength {
		maxLength = len(str2)
	}

	djustedMaxLength := maxLength + (len(str1) + len(str2))

	if maxLength == 0 {
		return 1.0
	}

	similarity := 1.0 - float64(distance)/float64(djustedMaxLength)
	return similarity
}
