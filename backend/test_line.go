package main

import (
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
)

func main() {
	l := line.New()
	fmt.Printf("Line type: %T\n", l)
}
