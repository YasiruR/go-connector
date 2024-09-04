package main

import "github.com/YasiruR/connector/boot"

func main() {
	boot.Start()
}

/*
Todos:
 - internal error handling (ctx)
 - external error definitions
 - log stack trace
 - authentication
 - data plane functions (exchange process)
*/

// todo: endpoints (transfers, start at provider), validations (data format, agreement, ack)
