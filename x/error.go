package x

import (
	"fmt"
)

func E(msg string, errors ...error) error {
	f := "msg: %s"
	for idx, _ := range errors {
		tmp := ",error_" + I2S(idx) + ": %w"
		f += tmp
	}
	return fmt.Errorf(f, msg, errors[:])
}
