package converters

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var humanReadableRegex = regexp.MustCompile(`(\d+)([BbMmKk]?)`)

func HumanReadable(value string) (uint64, error) {
	result := humanReadableRegex.FindAllStringSubmatch(value, -1)
	if result == nil || len(result) != 1 || len(result[0]) != 3 {
		return 0, errors.New(fmt.Sprintf("Failed to parse %s. result is %q", value, result))
	}

	quantity, err := strconv.Atoi(result[0][1])
	if err != nil {
		return 0, err
	}

	switch label := strings.ToLower(result[0][2]); label {
	case "k":
		return uint64(quantity * 1000), nil
	case "m":
		return uint64(quantity * 1000000), nil
	case "b":
		return uint64(quantity * 1000000000), nil
	default:
		return uint64(quantity), nil
	}
}
