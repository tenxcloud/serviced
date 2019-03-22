/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import (
	"github.com/json-iterator/go"
)

type serializerFactory struct {
}

func (serializerFactory) New() (s interface{}, err error) {
	s = &json{impl: jsoniter.ConfigFastest}
	return
}

type Serializer interface {
	Serialize(value interface{}) (content []byte, err error)
	Deserialize(content []byte, value interface{}) (err error)
}

type json struct {
	impl jsoniter.API
}

func (s json) Serialize(value interface{}) (content []byte, err error) {
	content, err = s.impl.Marshal(value)
	return
}

func (s json) Deserialize(content []byte, value interface{}) (err error) {
	err = s.impl.Unmarshal(content, value)
	return
}
