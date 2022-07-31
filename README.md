# go-map-functions
Useful functions when working with hashmaps or configuring a field to field mapping

## Examples
The following structs and variables will be used in the examples below
```go
type Price struct {
    Id    int     `json:"id""`
    Value float32 `json:"value"`
}

type RelatedProduct struct {
    Name string `json:"name"`
}

type Product struct {
    Name           string           `json:"name"`
    Description    string           `json:"description"`
    Prices         []Price          `json:"prices"`	
    RelatedProduct RelatedProduct   `json:"relatedProduct"`
}

p := Product{
	Name: "Sample", 
	Description: "Sample Product"
	Prices: []Price {
		{
			Id: 1, 
			Value: 10.0
		},
		{
			Id: 2, 
			Value: 15.0
		}
	}, 
	RelatedProduct: {
		Name: "Related Sample Product"
	}
}

// building a map out of the sample struct
var productMap map[string]any
jsonProduct, _ := json.Marshal(p)
json.Unmarshal(jsonProduct, &productMap)
```

### Accessing Map Fields by using Dot-Notation
```go
// building a map out of struct for demonstrate what could be possible
import mf "github.com/flo-hame/go-map-functions"

...

// accessing productMap name
productName, err := mf.GetValueByFieldPathDotNotation("name", productMap)

// accessing nested related product name
relatedProductName, err := mf.GetValueByFieldPathDotNotation("relatedProduct.name", productMap)

// accessing nested value within a slice
secondPrice, err := mf.GetValueByFieldPathDotNotation("prices.[1].value", productMap)

...

```

### Using Field Mapping
Doc will be coming soon...

