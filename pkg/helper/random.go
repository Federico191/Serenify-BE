package helper

import (
	"github.com/thanhpk/randstr"
)

func GenerateCode() string {
	return randstr.Hex(20)
}