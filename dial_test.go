package rotaryphone

import (
	"fmt"
	"testing"
)

func TestDialFunctional(t *testing.T) {
	dial := NewDial()
	for {
		digit := <-dial.Digit
		fmt.Println("Dialed digit", digit)
	}
}
