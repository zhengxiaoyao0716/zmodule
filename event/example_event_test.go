package event_test

import (
	"fmt"

	"github.com/zhengxiaoyao0716/zmodule/event"
)

var ks = [][2]string{
	{"CUS", "start"},
	{"CUS", "stop"},
}

const (
	KeyCustom event.KeyIndex = event.KeyCustom + iota // your.KeyCustom == event.KeyCustom
	KeyStart
	KeyStop
)

func Example() {
	event.OnInit(func(e event.Event) error { fmt.Println(e); return nil })
	event.Init(nil, ks...)

	event.On(KeyStart, func(e event.Event) error { fmt.Println(e); return nil })
	event.On(KeyStop, func(e event.Event) error { fmt.Println(e); return nil })
	event.Emit(KeyStart, nil)
	event.Emit(KeyStop, nil)
}
