package common

import (
	"errors"
	"strconv"
	"unicode"
)

func ParseMessageSize(maxMessageSize string) (int, error) {
	if maxMessageSize == "0" || len(maxMessageSize) < 2 || maxMessageSize[0] == '0' {
		return 0, errors.New("invalid max_message_size")
	}

	r := []rune(maxMessageSize)
	sizeR := []rune{}
	idx := 0

	for unicode.IsDigit(r[idx]) {
		sizeR = append(sizeR, r[idx])
		idx++
	}

	size, err := strconv.Atoi(string(sizeR))
	if err != nil {
		return 0, errors.New("invalid max_message_size")
	}
	switch maxMessageSize[idx:] {
	case "B", "b", "":
		break
	case "KB", "kb", "Kb":
		size = size * 1024
	case "MB", "mb", "Mb":
		size = size * 1024 * 1024
	case "GB", "gb", "Gb":
		size = size * 1024 * 1024 * 1024
	default:
		return 0, errors.New("invalid max_message_size")
	}
	return size, nil
}
