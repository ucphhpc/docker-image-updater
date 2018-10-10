package main

import "fmt"

type arrayFlags map[string]struct{}

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(Value string) error {

	*i["sdfdsf"]struct {}{}
	return nil
}