package anxiety

import (
	"errors"
	"testing"
)

func TestPanicWithBetaBlockers(t *testing.T) {
	BetaBlockers = true

	defer func() {
		if rvr := recover(); rvr != nil {
			t.Fatal("panicked, wanted no panic")
		}
	}()

	Panic(errors.New("test"))
}

func TestPanicWithoutBetaBlockers(t *testing.T) {
	BetaBlockers = false

	defer func() {
		if rvr := recover(); rvr == nil {
			t.Fatal("no panic, wanted panic")
		}
	}()

	Panic(errors.New("test"))
}
