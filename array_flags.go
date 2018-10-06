package main

import "fmt"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(Value string) error {
	*i = append(*i, Value)
	return nil
}