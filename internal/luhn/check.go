package luhn

import (
	"fmt"
	"github.com/dlc/go-market/internal/model/apperrors"
	"strconv"
)

func ValidIDErr(numberS string) error {
	number, err := strconv.Atoi(numberS)
	if err != nil {
		return fmt.Errorf("error while decoding order id : %s", err)
	}
	if !ValidID(number) {
		return apperrors.NewUnprocessableContent("wrong id")
	}
	return nil
}

func ValidID(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
