/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import "reflect"

type object struct {
	Type  string
	Value []byte
}

type Message struct {
	RequestID    string
	IsResponse   bool
	Service      string
	Method       string
	Arguments    []object
	ReturnValues []object
}

func objectToInterface(kernel Kernel, serializer Serializer, objects []object) (interfaces []interface{}, err error) {
	var values []reflect.Value
	if values, err = objectToValue(kernel, serializer, objects); err != nil {
		return
	}
	count := len(values)
	interfaces = make([]interface{}, count)
	for i := 0; i < count; i++ {
		interfaces[i] = values[i].Interface()
	}
	return
}

func interfaceToObject(serializer Serializer, interfaces []interface{}) (objects []object, err error) {
	count := len(interfaces)
	objects = make([]object, count)
	for i := 0; i < count; i++ {
		argument := interfaces[i]
		var value []byte
		if value, err = serializer.Serialize(argument); err != nil {
			return
		}
		objects[i] = object{
			Type:  reflect.Indirect(reflect.ValueOf(argument)).Type().String(),
			Value: value,
		}
	}
	return
}

func objectToValue(kernel Kernel, serializer Serializer, objects []object) (values []reflect.Value, err error) {
	count := len(objects)
	values = make([]reflect.Value, count)
	for i := 0; i < count; i++ {
		object := &(objects[i])
		var theType reflect.Type
		if theType, err = kernel.ResolveType(object.Type); err != nil {
			return
		}
		rv := reflect.New(theType)
		if err = serializer.Deserialize(object.Value, rv.Interface()); err != nil {
			return
		}
		values[i] = reflect.Indirect(rv)
	}
	return
}

func valueToObject(serializer Serializer, values []reflect.Value) (objects []object, err error) {
	count := len(values)
	interfaces := make([]interface{}, count)
	for i := 0; i < count; i++ {
		interfaces[i] = reflect.Indirect(values[i]).Interface()
	}
	objects, err = interfaceToObject(serializer, interfaces)
	return
}
