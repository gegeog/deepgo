package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserService struct {
	NotEmptyStruct bool
}
type MessageService struct {
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]func() interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]func() interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	fn, ok := constructor.(func() interface{})
	if ok {
		c.constructors[name] = fn
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	if c.constructors[name] == nil {
		return nil, errors.New("no constructor found")
	}

	return c.constructors[name](), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
