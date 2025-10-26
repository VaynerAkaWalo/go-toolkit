package xevent

import (
	"context"
	"reflect"
	"sync"
)

type (
	Broker struct {
		channels map[reflect.Type]interface{}
		mutex    *sync.RWMutex
	}
)

func NewBroker(types ...reflect.Type) *Broker {
	channels := make(map[reflect.Type]interface{})

	for _, t := range types {
		channels[t] = reflect.SliceOf(t)
	}

	return &Broker{
		channels: channels,
		mutex:    &sync.RWMutex{},
	}
}

func RegisterListener[T any](broker *Broker, ctx context.Context) chan T {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	channel := make(chan T)

	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	rawChannels, ok := broker.channels[eventType]
	if !ok {
		panic("cannot process unregistered events")
	}

	channels, ok := rawChannels.([]chan T)

	channels = append(channels, channel)

	broker.channels[eventType] = channels

	return channel
}

func RemoveListener[T any](broker *Broker, ctx context.Context, channel chan T) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	rawChannels, ok := broker.channels[eventType]
	if !ok {
		panic("cannot process unregistered events")
	}

	channels, ok := rawChannels.([]chan T)

	close(channel)

	for index, cha := range channels {
		if cha == channel {
			broker.channels[eventType] = append(channels[:index], channels[index+1:]...)
			return
		}
	}
}

func PublishEvent[T any](broker *Broker, ctx context.Context, event T) {
	eventType := reflect.TypeOf(event)

	broker.mutex.RLock()
	defer broker.mutex.RUnlock()

	rawChannels, ok := broker.channels[eventType]
	if !ok {
		panic("cannot process unregistered events")
	}

	channels, ok := rawChannels.([]chan T)

	for _, cha := range channels {
		cha <- event
	}
}
