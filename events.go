package events

import (
	"log"
	"reflect"
	"sync"
)

const (
	Version             = "0.0.1"
	DefaultMaxListeners = 0
	EnableWarning       = false
)

type (
	EventName string
	Listener  func(...interface{})
	Events    map[EventName][]Listener

	EventEmitter interface {
		AddListener(EventName, ...Listener)
		Emit(EventName, ...interface{})
		EventNames() []EventName
		GetMaxListeners() int
		ListenerCount(EventName) int
		Listeners(EventName) []Listener
		On(EventName, ...Listener)
		Once(EventName, ...Listener)
		RemoveAllListeners(EventName) bool
		RemoveListener(EventName, Listener) bool
		Clear()
		SetMaxListeners(int)
		Len() int
	}

	emitter struct {
		maxListeners int
		evtListeners sync.Map
		mu           sync.Mutex
	}
)

func New() EventEmitter {
	return &emitter{maxListeners: DefaultMaxListeners}
}

var (
	_              EventEmitter = &emitter{}
	defaultEmitter EventEmitter = New()
)

func AddListener(evt EventName, listener ...Listener) {
	defaultEmitter.AddListener(evt, listener...)
}

func (e *emitter) AddListener(evt EventName, listeners ...Listener) {
	e.mu.Lock()
	defer e.mu.Unlock()

	currentListenersRaw, _ := e.evtListeners.LoadOrStore(evt, []Listener{})
	currentListeners := currentListenersRaw.([]Listener)

	if e.maxListeners > 0 && len(currentListeners) >= e.maxListeners {
		if EnableWarning {
			log.Printf(`(events) warning: possible EventEmitter memory leak detected. %d listeners added. Use emitter.SetMaxListeners(n int) to increase limit.`, len(currentListeners))
		}
		return
	}

	e.evtListeners.Store(evt, append(currentListeners, listeners...))
}

func Emit(evt EventName, data ...interface{}) {
	defaultEmitter.Emit(evt, data...)
}

func (e *emitter) Emit(evt EventName, data ...interface{}) {
	if listenersRaw, ok := e.evtListeners.Load(evt); ok {
		listeners := listenersRaw.([]Listener)
		for _, listener := range listeners {
			listener(data...)
		}
	}
}

func EventNames() []EventName {
	return defaultEmitter.EventNames()
}

func (e *emitter) EventNames() []EventName {
	var names []EventName
	e.evtListeners.Range(func(key, value interface{}) bool {
		names = append(names, key.(EventName))
		return true
	})
	return names
}

func GetMaxListeners() int {
	return defaultEmitter.GetMaxListeners()
}

func (e *emitter) GetMaxListeners() int {
	return e.maxListeners
}

func ListenerCount(evt EventName) int {
	return defaultEmitter.ListenerCount(evt)
}

func (e *emitter) ListenerCount(evt EventName) int {
	if listenersRaw, ok := e.evtListeners.Load(evt); ok {
		listeners := listenersRaw.([]Listener)
		return len(listeners)
	}
	return 0
}

func Listeners(evt EventName) []Listener {
	return defaultEmitter.Listeners(evt)
}

func (e *emitter) Listeners(evt EventName) []Listener {
	if listenersRaw, ok := e.evtListeners.Load(evt); ok {
		listeners := listenersRaw.([]Listener)
		return listeners
	}
	return nil
}

func On(evt EventName, listener ...Listener) {
	defaultEmitter.On(evt, listener...)
}

func (e *emitter) On(evt EventName, listeners ...Listener) {
	e.AddListener(evt, listeners...)
}

func Once(evt EventName, listener ...Listener) {
	defaultEmitter.Once(evt, listener...)
}

func (e *emitter) Once(evt EventName, listeners ...Listener) {
	e.AddListener(evt, listeners...)
}

func RemoveAllListeners(evt EventName) bool {
	return defaultEmitter.RemoveAllListeners(evt)
}

func (e *emitter) RemoveAllListeners(evt EventName) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.evtListeners.Delete(evt)
	_, found := e.evtListeners.Load(evt)
	return !found
}

func RemoveListener(evt EventName, listener Listener) bool {
	return defaultEmitter.RemoveListener(evt, listener)
}

func (e *emitter) RemoveListener(evt EventName, listener Listener) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if listenersRaw, ok := e.evtListeners.Load(evt); ok {
		listeners := listenersRaw.([]Listener)
		for i, l := range listeners {
			if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
				e.evtListeners.Store(evt, append(listeners[:i], listeners[i+1:]...))
				return true
			}
		}
	}
	return false
}

func Clear() {
	defaultEmitter.Clear()
}

func (e *emitter) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.evtListeners = sync.Map{}
}

func SetMaxListeners(n int) {
	defaultEmitter.SetMaxListeners(n)
}

func (e *emitter) SetMaxListeners(n int) {
	if n < 0 {
		if EnableWarning {
			log.Printf("(events) warning: MaxListeners must be a positive number, tried to set: %d", n)
		}
		return
	}
	e.maxListeners = n
}

func Len() int {
	return defaultEmitter.Len()
}

func (e *emitter) Len() int {
	length := 0
	e.evtListeners.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}
