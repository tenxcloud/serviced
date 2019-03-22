/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package main

import (
	"context"
	"fmt"
	"github.com/tedli/serviced"
)

type One struct {
	Value string
}

type Two struct {
	Value string
}

type RemoteService interface {
	Do(one One) (two Two)
}

type proxy struct {
	serviced.TransparentProxy
}

func (p proxy) Do(one One) (two Two) {
	rvs, e := p.TransparentProxy.Call("Do", one)
	if e != nil {
		panic(e)
	}
	two = rvs[0].(Two)
	return
}

type LocalService interface {
	Do(two Two) (one One)
	GetDependence() RemoteService
}

type module struct {
	// this field can be auto injected
	Dependence RemoteService
}

func (module) Do(two Two) (one One) {
	fmt.Println("module")
	fmt.Printf("Two: %s\n", two.Value)
	one = One{Value: "one"}
	return
}

func (m module) GetDependence() RemoteService {
	return m.Dependence
}

// This is a debug entry. Server side should implement the service
// interface to process real business logic. Then return value can
// be returned through client proxy. On client side, by resolving
// a service object from kernel, then the service can be consumed
// through method calls, not matter the underneath object is a
// transparent proxy, or a module implementation struct instance.
func (proxy) Call(method string, arguments []interface{}) (returnValues []interface{}) {
	fmt.Printf("Method: %s\n", method)
	fmt.Printf("One: %s\n", arguments[0].(One).Value)
	returnValues = []interface{}{Two{Value: "two"}}
	return
}

// the io struct now is for debug and demo purpose.
// in real world, it should wrap a websocket or
// something like a connection object.
type io struct {
	channel chan *serviced.Message
}

type ioFactory struct {
}

func (ioFactory) New() (s interface{}, err error) {
	s = &io{channel: make(chan *serviced.Message, 100)}
	return
}

func (i io) In() (rawMessage <-chan *serviced.Message) {
	return i.channel
}

func (i io) Out() (rawMessage chan<- *serviced.Message) {
	return i.channel
}

func main() {
	builder := serviced.NewKernelBuilder()
	if err := serviced.RegisterDefault(builder); err != nil {
		panic(err)
	}
	if err := builder.RegisterService(new(RemoteService), new(proxy)); err != nil {
		panic(err)
	}
	if err := builder.RegisterFactory(new(serviced.IO), new(ioFactory)); err != nil {
		panic(err)
	}
	if err := builder.RegisterType(new(One)); err != nil {
		panic(err)
	}
	if err := builder.RegisterType(new(Two)); err != nil {
		panic(err)
	}
	if err := builder.RegisterService(new(LocalService), new(module)); err != nil {
		panic(err)
	}
	kernel := builder.Build()
	server, err := kernel.ResolveService(new(serviced.Server))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := server.(serviced.Server).Start(context.Background()); err != nil {
			panic(err)
		}
	}()
	service, err := kernel.ResolveService(new(LocalService))
	if err != nil {
		panic(err)
	}
	ls := service.(LocalService)
	one := ls.Do(Two{Value: "two"})
	fmt.Printf("Local service invocation: %s\n", one.Value)
	two := ls.GetDependence().Do(One{Value: "one"})
	fmt.Printf("Remote service: %s", two.Value)
}
