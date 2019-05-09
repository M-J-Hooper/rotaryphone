package main

import (
	"time"

	"github.com/M-J-Hooper/rotaryphone"
)

func main() {
	phone := rotaryphone.NewRotary(4 * time.Second)
	phone.Run()
	//phone.Debug()
}
