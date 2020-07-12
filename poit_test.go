package poit

import (
	"fmt"
	"testing"
)

func TestSearchTygelsjo(t *testing.T) {

	for a := range Search("Tygelsjö") {
		fmt.Println(a)
	}

}

func TestSearchMalmo(t *testing.T) {

	for a := range Search("Malmö") {
		fmt.Println(a)
	}

}
