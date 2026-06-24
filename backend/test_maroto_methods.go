package main

import (
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"reflect"
)

func main() {
	r := row.New(10)
	fmt.Printf("Row methods:\n")
	t := reflect.TypeOf(r)
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Printf("- %s\n", t.Method(i).Name)
	}

	c := col.New(12)
	fmt.Printf("\nCol methods:\n")
	tc := reflect.TypeOf(c)
	for i := 0; i < tc.NumMethod(); i++ {
		fmt.Printf("- %s\n", tc.Method(i).Name)
	}
}
