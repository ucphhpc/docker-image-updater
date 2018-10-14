package main

import "strings"

type arrayFlags map[string]struct{}

func (i *arrayFlags) String() string {
	arr := make([]string, 0)
	for k := range *i {
		arr = append(arr, k)
	}
	return strings.Join(arr, ",")
}

func (i *arrayFlags) Set(Value string) error {
	if *i == nil {
		*i = make(arrayFlags)
	}
	(*i)[Value] = struct{}{}
	return nil
}