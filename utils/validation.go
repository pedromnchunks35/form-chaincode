package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func RemoveStringSpaces(value string) string {
	return strings.ReplaceAll(value, " ", "")
}

func IsValidString(value string) bool {
	return len(value) != 0
}

func ValueExists(err error, value []byte) bool {
	return err != nil || len(value) != 0
}

func ValidatePageAndSize(pageString string, sizeString string) (int, int, error) {
	page, err := convertStringToInt(pageString)
	if err != nil {
		return 0, 0, err
	}

	size, err := convertStringToInt(sizeString)
	if err != nil {
		return 0, 0, err
	}

	bothLegit := arePageAndSizeLegit(page, size)
	if !bothLegit {
		return 0, 0, fmt.Errorf("page and size are not consistent")
	}

	return page, size, nil
}

func arePageAndSizeLegit(page int, size int) bool {
	return !isNumberNegative(page) && !isNumberNegative(size) && isNumberDifferentThatZero(size)
}

func convertStringToInt(toConvert string) (int, error) {
	return strconv.Atoi(toConvert)
}

func isNumberNegative(number int) bool {
	return number < 0
}

func isNumberDifferentThatZero(number int) bool {
	return number != 0
}
