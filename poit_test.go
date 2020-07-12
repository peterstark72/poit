package poit

import (
	"fmt"
	"testing"
)

func ExamplePoitClient_Search() {
	poc := NewClient()

	for a := range poc.Search("Tygelsjö") {
		fmt.Println(a)
	}
}

func TestSearchTygelsjo(t *testing.T) {

	poc := NewClient()

	var announcements []Announcement
	for a := range poc.Search("Tygelsjö") {
		announcements = append(announcements, a)
	}
	if len(announcements) == 0 {
		t.Error("No announcements found")
	}
}

func TestSearchMalmo(t *testing.T) {

	poc := NewClient()
	for a := range poc.Search("Malmö") {
		fmt.Println(a)
		if a.Type == "Bygglov" {
			for n, s := range poc.GetAnnouncementText(a) {
				fmt.Println(n, s)
			}
		}
	}

}
