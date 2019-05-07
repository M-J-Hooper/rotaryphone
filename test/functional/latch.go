package main

import (
    "fmt"
    "github.com/M-J-Hooper/rotaryphone"
)

func main() {
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