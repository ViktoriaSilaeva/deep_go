package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	container map[string]func() any
	singleton map[string]any
}

func NewContainer() *Container {
	return &Container{
		container: make(map[string]func() any),
		singleton: make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	switch constr := constructor.(type) {
	case func() any:
		c.container[name] = constr
	default:
		panic("unexpected constructor type")
	}
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	switch constr := constructor.(type) {
	case func() any:
		c.singleton[name] = constr()
	default:
		panic("unexpected constructor type")
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	if function, ok := c.container[name]; ok {
		return function(), nil
	}
	return nil, errors.New("unexpected service name")
}

func (c *Container) ResolveSingletonType(name string) (interface{}, error) {
	if value, ok := c.singleton[name]; ok {
		return value, nil
	}
	return nil, errors.New("unexpected service name")
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

	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{}
	})

	userService3, err := container.ResolveSingletonType("UserService")
	assert.NoError(t, err)
	userService4, err := container.ResolveSingletonType("UserService")
	assert.NoError(t, err)

	u3 := userService3.(*UserService)
	u4 := userService4.(*UserService)
	assert.True(t, u3 == u4)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
