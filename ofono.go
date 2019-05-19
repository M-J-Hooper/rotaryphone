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

type OfonoPhone struct {
	conn *dbus.Conn
}

func NewOfonoPhone() *OfonoPhone {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Println("Failed to connect to SystemBus bus:", err)
		return nil //what to do about errors...
	}
	return &OfonoPhone{conn}
}

func (o *OfonoPhone) GetModem() dbus.ObjectPath {
	var modems [][]interface{}
	err := o.conn.Object(Object, "/").Call(GetModemInterface, 0).Store(&modems)
	if err != nil {
		fmt.Println("Failed to get modems:", err)
	}
	fmt.Println("Got modems", modems)
	for _, modem := range modems {
		path, details := modem[0].(dbus.ObjectPath), modem[1].(map[string]dbus.Variant)
		if details["Online"].Value().(bool) {
			return path
		}
	}
	return ""
}

func (o *OfonoPhone) Call(number string) {
	var path string
	err := o.conn.Object(Object, o.GetModem()).Call(DialInterface, 0, number, "default").Store(&path)
	if err != nil {
		fmt.Println("Failed to dial:", err)
		return
	}
	return
}

func (o *OfonoPhone) Hangup() {
	var s string
	err := o.conn.Object(Object, o.GetModem()).Call(HangupAllInterface, 0).Store(&s)
	if err != nil {
		fmt.Println("Failed to dial:", err)
		return
	}
	return
}

func (o OfonoPhone) Debug() {
	var s string
	err := o.conn.Object(Object, "/").Call(DebugInterface, 0).Store(&s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to introspect ofono root", err)
	}
	fmt.Println("Root introspection:", s)

	modem := o.GetModem()
	fmt.Println("Modem:", modem)

	var s2 string
	err = o.conn.Object(Object, modem).Call(DebugInterface, 0).Store(&s2)
	if err != nil {
		fmt.Println("Failed to introspect ofono modem", err)
	}
	fmt.Println("Modem introspection:", s2)
	println()
}
