package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var res strings.Builder
	var prev rune
	var isPrevDoubleSlash bool
	for _, cur := range str {
		num, err := strconv.Atoi(string(cur))

		var isCurNumber bool
		if err == nil {
			isCurNumber = true
		}

		if prev == '\\' && !isPrevDoubleSlash {
			// Slash + slash? Set the flag and skip
			if cur == '\\' {
				isPrevDoubleSlash = true
				continue
			}

			// Slash + number? Write as prev and skip
			if isCurNumber {
				isPrevDoubleSlash = false
				prev = rune(num) + '0'
				continue
			}

			// Slash + other? Return error
			return "", ErrInvalidString
		}

		// The number first? Return error
		if prev == 0 && isCurNumber {
			return "", ErrInvalidString
		}

		// Not shilding number? Repeat symbol num times
		if isCurNumber {
			res.WriteString(strings.Repeat(string(prev), num))
			prev = 0
			isPrevDoubleSlash = false
			continue
		}

		if prev != 0 {
			res.WriteRune(prev)
		}

		prev = cur
		isPrevDoubleSlash = false
	}

	if prev != 0 {
		res.WriteRune(prev)
	}
	return res.String(), nil
}
