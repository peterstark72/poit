# POIT

Package poit implements a search query to the POIT website, poit.bolagsverket.se.

## Example usage

```
func ExampleClient_Search() {
	poc := NewClient()
	for a := range poc.Search("Tygelsjö") {
		fmt.Println(a)
	}
}
```
