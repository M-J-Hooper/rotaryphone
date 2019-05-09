package functional

import (
	"fmt"
	"testing"

	"github.com/M-J-Hooper/rotaryphone"
)

func TestLatchFunctional(t *testing.T) {
	latch := rotaryphone.NewLatch()
	for {
		active := <-latch.Active
		if active {
			fmt.Println("Latch active")
		} else {
			fmt.Println("Latch inactive")
		}
	}
}
