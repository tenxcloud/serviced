/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import "reflect"

type Factory interface {
	New() (value interface{}, err error)
}

type Kernel interface {
	ResolveService(contract interface{}) (service interface{}, err error)
	ResolveType(typeName string) (theType reflect.Type, err error)
}

type KernelBuilder interface {
	RegisterType(value interface{}) (err error)
	RegisterCompatibleType(typeName string, value interface{}) (err error)
	RegisterService(contract interface{}, service interface{}) (err error)
	RegisterFactory(contract interface{}, service Factory) (err error)
	Build() (kernel Kernel)
}

func NewKernelBuilder() KernelBuilder {
	return &kernel{
		types:     make(map[string]reflect.Type),
		factories: make(map[string]reflect.Value),
	}
}

type kernel struct {
	types     map[string]reflect.Type
	factories map[string]reflect.Value
}

func (k *kernel) RegisterType(value interface{}) (err error) {
	err = k.registerInternal(nil, value)
	return
}

func (k *kernel) RegisterCompatibleType(typeName string, value interface{}) (err error) {
	err = k.registerInternal(&typeName, value)
	return
}

func (k *kernel) RegisterFactory(contract interface{}, service Factory) (err error) {
	v := reflect.Indirect(reflect.ValueOf(contract))
	if v.Kind() != reflect.Interface {
		err = ErrContractNotInterface
		return
	}
	key := v.Type().String()
	if err = k.registerInternal(&key, service); err != nil {
		return
	}
	err = k.RegisterType(service)
	return
}

func (k *kernel) registerInternal(typeName *string, value interface{}) (err error) {
	if value == nil {
		err = ErrRegisterNilValue
		return
	}
	vv := reflect.Indirect(reflect.ValueOf(value))
	var key string
	if typeName == nil {
		key = vv.Type().String()
	} else if key = *typeName; key == "" {
		err = ErrRegisterWithEmptyTypeName
		return
	}
	if _, known := k.types[key]; known {
		err = ErrTypeAlreadyRegistered
		return
	}
	k.types[key] = vv.Type()
	return
}

func (k *kernel) RegisterService(contract interface{}, service interface{}) (err error) {
	v := reflect.Indirect(reflect.ValueOf(contract))
	if v.Kind() != reflect.Interface {
		err = ErrContractNotInterface
		return
	}
	if !reflect.TypeOf(service).Implements(v.Type()) {
		err = ErrTypeNotImplementService
		return
	}
	key := v.Type().String()
	if err = k.registerInternal(&key, service); err != nil {
		return
	}
	err = k.RegisterType(service)
	return
}

func (k kernel) Build() (kernel Kernel) {
	return k
}

var (
	kernelKey = reflect.Indirect(reflect.ValueOf(new(Kernel))).Type().String()
)

func (k kernel) ResolveService(contract interface{}) (service interface{}, err error) {
	var key string
	if k, ok := contract.(string); ok {
		key = k
	} else {
		key = reflect.Indirect(reflect.ValueOf(contract)).Type().String()
	}
	if key == kernelKey {
		service = k
		return
	}
	service, err = k.resolveServiceInternal(key)
	return
}

func (k kernel) ResolveType(typeName string) (theType reflect.Type, err error) {
	var known bool
	if theType, known = k.types[typeName]; !known {
		err = ErrUnresolvableType
	}
	return
}

func (k kernel) resolveServiceInternal(key string) (service interface{}, err error) {
	if t, known := k.types[key]; known {
		service = reflect.New(t).Interface()
	} else {
		err = ErrUnresolvableType
		return
	}
	if err = k.injectDependencies(service); err != nil {
		return
	}
	if factory, ok := service.(Factory); ok {
		if service, err = factory.New(); err != nil {
			return
		}
		if err = k.injectDependencies(service); err != nil {
			return
		}
	}
	return
}

func (k kernel) injectDependencies(service interface{}) (err error) {
	v := reflect.Indirect(reflect.ValueOf(service))
	for i, fieldCount := 0, v.NumField(); i < fieldCount; i++ {
		fv := reflect.Indirect(v.Field(i))
		if fvt := fv.Type().String(); fvt == kernelKey {
			fv.Set(reflect.ValueOf(k))
		} else if fs, e := k.resolveServiceInternal(fvt); e == nil {
			dv := reflect.Indirect(reflect.ValueOf(fs))
			if _, ok := fs.(TransparentProxy); ok {
				handleTransparentProxy(dv, v)
			}
			fv.Set(dv)
		} else if e != ErrUnresolvableType {
			err = e
			return
		}
	}
	return
}

func handleTransparentProxy(value, service reflect.Value) {
	value.FieldByName("Service").Set(reflect.ValueOf(service.Type()))
}
