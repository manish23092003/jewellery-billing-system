package main

import (
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"reflect"
)

func main() {
	// Let's print out fields in props.Color and props.Text
	t := props.Text{}
	fmt.Printf("Text props: %+v\n", reflect.TypeOf(t))
	for i := 0; i < reflect.TypeOf(t).NumField(); i++ {
		fmt.Printf("- %s\n", reflect.TypeOf(t).Field(i).Name)
	}

	r := row.New(10)
	fmt.Printf("Row type: %T\n", r)
	
	c := props.Cell{}
	fmt.Printf("Cell props: %+v\n", reflect.TypeOf(c))
	for i := 0; i < reflect.TypeOf(c).NumField(); i++ {
		fmt.Printf("- %s\n", reflect.TypeOf(c).Field(i).Name)
	}
}
