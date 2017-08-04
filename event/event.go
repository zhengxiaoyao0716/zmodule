package event

import (
	"github.com/zhengxiaoyao0716/util/event"
)

// Event .
type Event event.Event

var (
	pool *event.Pool
	ks   = []event.Key{
		event.Key{"ERR", "unused"},

		event.Key{"SYS", "init"},
		event.Key{"SYS", "start"},
		event.Key{"SYS", "stop"},

		event.Key{"ERR", "custom"},
	}
)

// KeyIndex is the index of prebuild event keys
type KeyIndex int

// Index emulate of prebuild event keys
const (
	// KeyUnused is used to mark the start of key index, should never be used.
	KeyUnused KeyIndex = iota

	keyInit // Notice the `k` is in lower case so that you should not use it directly, either the value (1) of it.
	KeyStart
	KeyStop

	// KeyCustom is used to mark the start of custom index (also end of sys), should never be used directly.
	// To custom your keys, used as:
	// ``` Go
	// var ks = [][2]string{
	// 	{"CUS", "start"},
	// 	{"CUS", "stop"},
	// }
	// const (
	// 	KeyCustom event.KeyIndex = event.KeyCustom + iota // your.KeyCustom == event.KeyCustom
	// 	KeyStart
	// 	KeyStop
	// )
	// ...n
	// event.Init(ks...)
	// ```
	KeyCustom
)

// Emit .
func Emit(i KeyIndex, d interface{}) { pool.Emit(ks[i], d) }

// On .
func On(i KeyIndex, h func(Event) error) {
	pool.On(ks[i], func(e event.Event) error { return h(Event(e)) })
}

// Pool provide a single instance of event pool
func Pool() *event.Pool { return pool }

var inits []func(Event) error

// Init .
func Init(payload interface{}, cusKs ...[2]string) {
	for _, k := range cusKs {
		ks = append(ks, k)
	}
	pool = event.NewRestrictPool(ks...)

	for _, h := range inits {
		On(keyInit, h)
	}
	inits = nil
	Emit(keyInit, payload)
}

// OnInit used to listener the `keyInit` event, witch send inside `Init` function.
func OnInit(h func(Event) error) {
	inits = append(inits, h)
}
