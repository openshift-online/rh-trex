package util

import (
	"context"
	"fmt"
)

// ToPtr returns a pointer copy of value.
func ToPtr[T any](v T) *T {
	return &v
}

// FromPtr returns the pointer value or empty.
func FromPtr[T any](v *T) T {
	if v == nil {
		return Empty[T]()
	}
	return *v
}

// FromEmptyPtr emulates ToPtr(FromPtr(x)) sequence
func FromEmptyPtr[T any](v *T) *T {
	if v == nil {
		x := Empty[T]()
		return &x
	}
	return v
}

// Empty returns an empty value of type T.
func Empty[T any]() T {
	var zero T
	return zero
}

func EmptyStringToNil(a string) *string {
	if a == "" {
		return nil
	}
	return &a
}

func NilToEmptyString(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

func GetAccountIDFromContext(ctx context.Context) string {
	accountID := ctx.Value("accountID")
	if accountID == nil {
		return ""
	}
	return fmt.Sprintf("%v", accountID)
}
