package mistake

import (
	"errors"
	"runtime"
)

func Unwrap(err error) {
	if err != nil {
		panic(errors.Join(runtime.StartTrace(), err))
	}
}

func UnwrapNotTrace(err error) {
	if err != nil {
		panic(err)
	}
}
