package example

import (
	"context"
	"fmt"
)

type TestFn func(context.Context)

type TestFn2 func(string)

type TestFn3 func(string) error

type TestFn4 func(string, int64) error

type TestFn5 func(string, int64) ([]byte, string)

func Test1(ctx context.Context, id string) {
	fmt.Println("test")
}

func Test2(id, name string) {}

func Test3(id, name string, count int64) {}

func Test4(id string) error { return nil }

func Test5(id, name string) error { return nil }

func Test6(id, name string, count int64) error { return nil }

func Test7(id string) ([]byte, error) { return nil, nil }

func Test8(id, name string) ([]byte, error) { return nil, nil }

func Test9(id, name string, count int64) ([]byte, error) { return nil, nil }

type Interface interface {
	Test1(id string)

	Test2(id, name string)

	Test3(id, name string, count int64)

	Test4(id string) error

	Test5(id, name string) error

	Test6(id, name string, count int64) error

	Test7(id string) ([]byte, error)

	Test8(id, name string) ([]byte, error)

	Test9(id, name string, count int64) ([]byte, error)
}

type Struct struct{}

func (s *Struct) Test1(ctx context.Context, id string) {}

func (s *Struct) Test2(id, name string) {}

func (s *Struct) Test3(id, name string, count int64) {}

func (s *Struct) Test4(id string) error { return nil }

func (s *Struct) Test5(id, name string) error { return nil }

func (s *Struct) Test6(id, name string, count int64) error { return nil }

func (s *Struct) Test7(id string) ([]byte, error) { return nil, nil }

func (s *Struct) Test8(id, name string) ([]byte, error) { return nil, nil }

func (s *Struct) Test9(id, name string, count int64) ([]byte, error) { return nil, nil }
