package oauth

import (
	"errors"
)

type ServiceData[T any] struct {
	data             T
	errorCode        string
	errorDescription string
	errorURI         string
}

func Call[T any](f func() (T, error)) (ServiceData[T], error) {
	data, err := f()
	if err != nil {
		return newErrorServiceData(data, err), err
	}

	return newSuccessServiceData(data), nil
}

func (d *ServiceData[T]) Data() T {
	return d.data
}

func (d *ServiceData[T]) ErrorCode() string {
	return d.errorCode
}

func (d *ServiceData[T]) ErrorDescription() string {
	return d.errorDescription
}

func (d *ServiceData[T]) ErrorURI() string {
	return d.errorURI
}

func newSuccessServiceData[T any](data T) ServiceData[T] {
	return ServiceData[T]{
		data: data,
	}
}

func newErrorServiceData[T any](data T, err error) ServiceData[T] {
	code := ErrorCodeServerError
	description := ErrorDescriptionInternalServerError
	uri := ""

	var validationErr *Error
	if errors.As(err, &validationErr) {
		code = validationErr.Code
		description = validationErr.Description
		uri = validationErr.URI

	}

	return ServiceData[T]{
		data:             data,
		errorCode:        code,
		errorDescription: description,
		errorURI:         uri,
	}
}
