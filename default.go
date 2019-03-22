/*
 * Licensed Materials - Property of tenxcloud.com
 * (C) Copyright 2019 TenxCloud. All Rights Reserved.
 * 2019-03-22  @author lizhen
 */

package serviced

func RegisterDefault(builder KernelBuilder) (err error) {
	if err = builder.RegisterService(new(Server), new(server)); err != nil {
		return
	}
	if err = builder.RegisterService(new(TransparentProxy), new(proxy)); err != nil {
		return
	}
	if err = builder.RegisterFactory(new(Invocations), new(invocationsFactory)); err != nil {
		return
	}
	if err = builder.RegisterFactory(new(InvocationQueue), new(invocationQueueFactory)); err != nil {
		return
	}
	if err = builder.RegisterFactory(new(Serializer), new(serializerFactory)); err != nil {
		return
	}
	return
}
