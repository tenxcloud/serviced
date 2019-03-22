/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-20  @author lizhen
 */

package serviced

import "errors"

var (
	ErrRegisterNilValue          = errors.New("register nil value")
	ErrTypeAlreadyRegistered     = errors.New("type already registered")
	ErrRegisterWithEmptyTypeName = errors.New("register with empty type name")
	ErrContractNotInterface      = errors.New("contract not interface")
	ErrUnresolvableType          = errors.New("unresolvable type")
	ErrTypeNotImplementService   = errors.New("type not implement service")
	ErrInvalidMessage            = errors.New("invalid message")
	ErrNoSuchInvocation          = errors.New("no such invocation")
	ErrNoSuchMethod              = errors.New("no such method")
)
