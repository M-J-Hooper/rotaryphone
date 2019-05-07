package main

import (
	"fmt"

	"github.com/M-J-Hooper/rotaryphone"
)

func main() {
	dial := rotaryphone.NewDial()
	for {
		digit := <-dial.Digit
		fmt.Println("Dialed digit", digit)
	}
}
