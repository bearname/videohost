package util

import "strconv"

func StrToInt(str string) (int, bool) {
	number, err := strconv.ParseInt(str, 0, 16)
	if err != nil {
		return 0, false
	}

	return int(number), true
}
