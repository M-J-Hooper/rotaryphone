package rotaryphone

import (
	"fmt"
	"testing"
)

func TestLatchFunctional(t *testing.T) {
	latch := NewLatch()
	for {
		active := <-latch.Active
		if active {
			fmt.Println("Latch active")
		} else {
			fmt.Println("Latch inactive")
		}
	}
}
