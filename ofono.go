package rotaryphone

import (
    "fmt"
    "os"

    "github.com/godbus/dbus"
)

const Object = "org.ofono"
const GetModemInterface = "org.ofono.Manager.GetModems"
const DialInterface = "org.ofono.VoiceCallManager.Dial"
const HangupAllInterface = "org.ofono.VoiceCallManager.HangupAll"
const DebugInterface = "org.freedesktop.DBus.Introspectable.Introspect"

type OfonoAdapter struct {
    conn dbus.Conn
}

func NewOfonoAdapter() *OfonoAdapter {
    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to connect to SystemBus bus:", err)
        return nil //what to do about errors...
    }
    return &OfonoAdapter{*conn}
}

func (o *OfonoAdapter) GetModem() dbus.ObjectPath {
    var modems [][]interface {}
    err := o.conn.Object(Object, "/").Call(GetModemInterface, 0).Store(&modems)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to get modems:", err)
    }
    return modems[0][0].(dbus.ObjectPath)
}

func (o *OfonoAdapter) Call(number string) {
    var path string
    err := o.conn.Object(Object, o.GetModem()).Call(DialInterface, 0, number, "default").Store(&path)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to dial:", err)
        return
    }
    return
}

func (o *OfonoAdapter) Hangup() {
    var s string
    err := o.conn.Object(Object, o.GetModem()).Call(HangupAllInterface, 0).Store(&s)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to dial:", err)
        return
    }
    return
}

func (o OfonoAdapter) Debug() {
    var s string
    err := o.conn.Object(Object, "/").Call(DebugInterface, 0).Store(&s)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to introspect ofono root", err)
    }
    fmt.Println("Root introspection:", s)
    println()

    modem := o.GetModem()
    fmt.Println("Modem:", modem)

    var s2 string
    err = o.conn.Object(Object, modem).Call(DebugInterface, 0).Store(&s2)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to introspect ofono modem", err)
    }
    fmt.Println("Modem introspection:", s2)
    println()
}
