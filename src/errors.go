package src

import (
	"fmt"
	"os"
	"time"
)

var globalErrors = make(chan error, 100)

func init() {
	go func() {
		f, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		for e := range globalErrors {
			f.WriteString(fmt.Sprintf("%s: %v\n", time.Now(), e))
			f.Sync()
		}
	}()
}

func LogErrorf(format string, args ...any) {
	globalErrors <- fmt.Errorf(format, args...)
}
