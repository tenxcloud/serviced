/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-22  @author lizhen
 */

package serviced

import (
	"github.com/satori/go.uuid"
)

const (
	Namespace = "c1e6f276-0d67-4a84-ad11-5cf62b3ce7fc"
)

var (
	namespace = uuid.Must(uuid.FromString(Namespace))
)

func NewUUID(service, method string) string {
	return uuid.NewV5(namespace, service+method).String()
}
