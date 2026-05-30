package api

import "github.com/jakej985-rgb/m3tal-core/core/services"

func StartStack(name string) error {
	return (&services.StackService{}).Up(name)
}

func StopStack(name string) error {
	return (&services.StackService{}).Down(name)
}
