/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import (
	"reflect"
	"sync"
)

type TransparentProxy interface {
	Call(method string, arguments ...interface{}) (returnValues []interface{}, err error)
}

type proxy struct {
	Service         reflect.Type
	Kernel          Kernel
	Serializer      Serializer
	Invocations     Invocations
	InvocationQueue InvocationQueue
}

func (p proxy) ensureMethod(method string) (err error) {
	if _, exist := p.Service.MethodByName(method); !exist {
		err = ErrNosuchMethod
	}
	return
}

func (p proxy) prepareArguments(arguments []interface{}) (objects []object, err error) {
	objects, err = interfaceToObject(p.Serializer, arguments)
	return
}

func buildRequestMessage(service, method string) *Message {
	return &Message{
		RequestID:  NewUUID(service, method),
		IsResponse: false,
		Service:    service,
		Method:     method,
	}
}

func (p proxy) prepareReturnValues(response *Message) (returnValues []interface{}, err error) {
	returnValues, err = objectToInterface(p.Kernel, p.Serializer, response.ReturnValues)
	return
}

func (p proxy) Call(method string, arguments ...interface{}) (returnValues []interface{}, err error) {
	if err = p.ensureMethod(method); err != nil {
		return
	}
	request := buildRequestMessage(p.Service.String(), method)
	var objects []object
	if objects, err = p.prepareArguments(arguments); err != nil {
		return
	}
	request.Arguments = objects
	responseChannel := p.Invocations.New(request.RequestID)
	defer p.Invocations.Remove(request.RequestID)
	p.InvocationQueue.In() <- request
	response := <-responseChannel
	returnValues, err = p.prepareReturnValues(response)
	return
}

var (
	i = &invocations{holder: new(sync.Map)}
	q = &queue{holder: make(chan *Message, 100)}
)

type invocationsFactory struct {
}

func (invocationsFactory) New() (s interface{}, err error) {
	s = i
	return
}

type invocationQueueFactory struct {
}

func (invocationQueueFactory) New() (s interface{}, err error) {
	s = q
	return
}

type Invocations interface {
	New(requestID string) (response <-chan *Message)
	Get(requestID string) (response chan<- *Message, err error)
	Remove(requestID string)
}

type InvocationQueue interface {
	Out() (out <-chan *Message)
	In() (in chan<- *Message)
	Close()
}

type invocations struct {
	holder *sync.Map
}

func (i invocations) New(requestID string) (response <-chan *Message) {
	r := make(chan *Message)
	i.holder.Store(requestID, r)
	response = r
	return
}

func (i invocations) Get(requestID string) (response chan<- *Message, err error) {
	invocation, exist := i.holder.Load(requestID)
	if !exist {
		err = ErrNoSuchInvocation
		return
	}
	response = invocation.(chan *Message)
	return
}

func (i invocations) Remove(requestID string) {
	i.holder.Delete(requestID)
}

type queue struct {
	holder chan *Message
}

func (q queue) Out() (out <-chan *Message) {
	return q.holder
}

func (q queue) In() (in chan<- *Message) {
	return q.holder
}

func (q queue) Close() {
	close(q.holder)
}
