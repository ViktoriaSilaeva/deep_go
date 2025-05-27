package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	count     uint64
	errString string
}

func (e *MultiError) Error() string {
	return fmt.Sprint(e.count) + " errors occured:\n" + e.errString + "\n"

}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	var me *MultiError
	if err == nil {
		me = &MultiError{}
	} else {
		me = err.(*MultiError)
	}

	for _, val := range errs {
		me.errString = me.errString + "\t* " + val.Error()
		me.count++
	}

	return me
}

func (e *MultiError) Is(target error) bool {
	if e == nil {
		return e == target
	}

	return strings.Contains(e.errString, target.Error())
}

func (e *MultiError) As(target error) bool {
	if e == nil {
		return e == target
	}

	return strings.Contains(e.errString, target.Error())
}

func (e *MultiError) Unwrap() error {
	split := strings.SplitN(e.errString, "\t", -1)
	return &MultiError{count: e.count - 1, errString: strings.Join(split[0:len(split)-1], "\t")}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	assert.True(t, err.(*MultiError).Is(errors.New("error 1")))
	assert.False(t, err.(*MultiError).As(errors.New("error 3")))

	expectedMessage = "1 errors occured:\n\t* error 1\n"
	assert.EqualError(t, err.(*MultiError).Unwrap(), expectedMessage)
}
