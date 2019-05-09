package functional

import (
	"fmt"
	"testing"

	"github.com/M-J-Hooper/rotaryphone"
)

func TestDialFunctional(t *testing.T) {
	dial := rotaryphone.NewDial()
	for {
		digit := <-dial.Digit
		fmt.Println("Dialed digit", digit)
	}
}
