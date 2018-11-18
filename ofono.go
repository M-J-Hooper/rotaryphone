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
    modem dbus.ObjectPath
}

func NewOfonoAdapter() *OfonoAdapter {
    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to connect to SystemBus bus:", err)
        return nil //what to do about errors...
    }

    var modems [][]interface {}
    err = conn.Object(Object, "/").Call(GetModemInterface, 0).Store(&modems)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to get modems:", err)
        return nil
    }

    return &OfonoAdapter{*conn, modems[0][0].(dbus.ObjectPath)}
}



func (o *OfonoAdapter) Call(number string) {
    var path string
    err := o.conn.Object(Object, o.modem).Call(DialInterface, 0, number, "default").Store(&path)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to dial:", err)
        return
    }
    return
}

func (o *OfonoAdapter) Hangup() {
    var s string
    err := o.conn.Object(Object, o.modem).Call(HangupAllInterface, 0).Store(&s)
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

    var s2 string
    err = o.conn.Object(Object, o.modem).Call(DebugInterface, 0).Store(&s2)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to introspect ofono modem", err)
    }
    fmt.Println("Modem introspection:", s2)
    println()

    fmt.Println("Modem:", o.modem)
}
