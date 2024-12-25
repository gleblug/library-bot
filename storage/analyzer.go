package storage

import (
	"slices"
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
)

type Tokens []string

func (t Tokens) OverlapCoefficient(query Tokens) int {
	if len(t) < 1 {
		return -1
	}
	sum := 0
	for _, queryToken := range query {
		closestToken := slices.MinFunc(t, func(t1 string, t2 string) int {
			return levenshtein.ComputeDistance(queryToken, t1) - levenshtein.ComputeDistance(queryToken, t2)
		})
		sum += levenshtein.ComputeDistance(queryToken, closestToken)
	}
	return sum
}

func tokenize(text string) Tokens {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func lowercaseFilter(tokens Tokens) Tokens {
	res := make([]string, len(tokens))
	for i, token := range tokens {
		res[i] = strings.ToLower(token)
	}
	return res
}

func Analyze(text string) Tokens {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	return tokens
}
