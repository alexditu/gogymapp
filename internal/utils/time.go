package utils

import (
	"os"
	"strings"
)

func LocalTimeZone() (string, error) {
	tz, err := os.ReadFile("/etc/timezone")
	if err != nil {
		return "", err
	}

	return strings.Trim(string(tz), " \n"), nil
}
