package utils

import (
	"fmt"
	"github.com/kr/pretty"
)

func AP(data interface{}) {
	fmt.Printf("%# v", pretty.Formatter(data))
}