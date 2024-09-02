package main

import (
	"errors"
	"fmt"
)

//func main() {
//	boot.Start()
//}

/*
Todos:
 - internal error handling (ctx)
 - external error definitions
 - log stack trace
 - authentication
 - data plane functions (exchange process)
*/

// todo: endpoints (transfers, start at provider), validations (data format, agreement, ack)

type ExternalError struct {
	err  error
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e ExternalError) Error() string {
	return e.err.Error()
}

func IncompatibleState(protocol, expected, received string) error {
	return ExternalError{
		err: fmt.Errorf("incompatible state (protocol: %s, expected: %s, received: %s)",
			protocol, expected, received),
		Code: 2001,
		Msg:  "testtt",
	}
}

func main() {
	//err := test()
	err := test2(`sample name`)
	//err := test1(`sampleeee`)
	e := &ExternalError{}
	if errors.As(err, e) {
		fmt.Println("YESSS: ", err)
		fmt.Println("okkk: ", e.Msg, e.Code)
	} else {
		fmt.Println("NOO: ", err)
	}

	ew := errors.Unwrap(err)
	fmt.Println("unwrapped: ", ew)
}

func test() error {
	return IncompatibleState(`test`, `test`, `test`)
}

func test2(name string) error {
	err := test1(name)
	return fmt.Errorf("err 2 (%s) - %w", name, err)
}

func test1(n string) error {
	err := test()
	return fmt.Errorf("err 2 %s- %w", n, err)
}
