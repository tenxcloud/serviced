/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import (
	"context"
	"reflect"
)

type IO interface {
	In() (rawMessage <-chan *Message)
	Out() (rawMessage chan<- *Message)
}

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	Queue       InvocationQueue
	Invocations Invocations
	Kernel      Kernel
	Serializer  Serializer
	IO          IO
}

func (s server) Start(ctx context.Context) error {
	defer s.Queue.Close()
	for {
		select {
		case message := <-s.IO.In():
			if err := s.dispatch(message); err != nil {
				return err
			}
		case message := <-s.Queue.Out():
			s.IO.Out() <- message
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s server) dispatch(message *Message) (err error) {
	if message.RequestID == "" {
		err = ErrInvalidMessage
		return
	}
	if message.IsResponse {
		err = s.handleResponse(message)
	} else {
		err = s.handleRequest(message)
	}
	return
}

func (s server) handleDebug(debugger Debugger, message *Message) (response *Message, err error) {
	var arguments []interface{}
	if arguments, err = objectToInterface(s.Kernel, s.Serializer, message.Arguments); err != nil {
		return
	}
	returnValues := debugger.Call(message.Method, arguments)
	response = new(Message)
	response.RequestID = message.RequestID
	response.IsResponse = true
	response.ReturnValues, err = interfaceToObject(s.Serializer, returnValues)
	return
}

func (s server) handleRequest(message *Message) (err error) {
	var service interface{}
	if service, err = s.Kernel.ResolveService(message.Service); err != nil {
		return
	}
	var response *Message
	if d, ok := service.(Debugger); ok {
		if response, err = s.handleDebug(d, message); err != nil {
			return
		}
	} else {
		if response, err = s.call(service, message); err != nil {
			return
		}
	}
	s.Queue.In() <- response
	return
}

func (s server) prepareArguments(service reflect.Value, message *Message) (arguments []reflect.Value, err error) {
	count := len(message.Arguments)
	var a []reflect.Value
	if a, err = objectToValue(s.Kernel, s.Serializer, message.Arguments); err != nil {
		return
	}
	arguments = make([]reflect.Value, 0, count+1)
	arguments[0] = service
	for i := 0; i < count; i++ {
		arguments[i+1] = a[i]
	}
	return
}

func (s server) prepareReturnValues(returnValues []reflect.Value) (objects []object, err error) {
	objects, err = valueToObject(s.Serializer, returnValues)
	return
}

func (s server) prepareResponse(request *Message, returnValues []reflect.Value) (response *Message, err error) {
	response = new(Message)
	response.RequestID = request.RequestID
	response.IsResponse = true
	response.ReturnValues, err = s.prepareReturnValues(returnValues)
	return
}

type Debugger interface {
	Call(method string, arguments []interface{}) (returnValues []interface{})
}

func (s server) call(service interface{}, message *Message) (response *Message, err error) {
	var arguments []reflect.Value
	sv := reflect.ValueOf(service)
	if arguments, err = s.prepareArguments(sv, message); err != nil {
		return
	}
	returnValues := sv.MethodByName(message.Method).Call(arguments)
	response, err = s.prepareResponse(message, returnValues)
	return
}

func (s server) handleResponse(message *Message) (err error) {
	var response chan<- *Message
	if response, err = s.Invocations.Get(message.RequestID); err != nil {
		return
	}
	defer close(response)
	response <- message
	return
}
