package rotaryphone

import (
    "time"
    "strconv"
)

type OutputAdapter interface {
    Call(string)
    Hangup()
}

type InputAdapter interface {
    Run()
    GetDigitChannel() chan int
    GetHangupChannel() chan struct{}
}

type Rotary struct {
    digitTimeout time.Duration
    input InputAdapter
    output OutputAdapter
}

func NewRotary(digitTimeout time.Duration) *Rotary {
    input := GpioAdapter{make(chan int), make(chan struct{})}
    output := NewOfonoAdapter()

    return &Rotary{digitTimeout, input, output}
}

func (r *Rotary) Run() {
    var number string
    lastDigit := time.Now()

    go r.input.Run()
    for {
        select {
        case digit := <-r.input.GetDigitChannel():
            number += strconv.Itoa(digit)
            println("New number is " + number)
            lastDigit = time.Now()
        case <-r.input.GetHangupChannel():
            println("Hanging up")
            //r.output.Hangup()
        default:
            if time.Since(lastDigit) > r.digitTimeout && len(number) > 0 {
                r.output.Call(number)
                println("Calling " + number)
                lastDigit = time.Now()
                number = ""
            }
            time.Sleep(100 * time.Millisecond)
        }
    }
}
