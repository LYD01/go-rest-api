// Album struct + validation
package albums

import "errors"

type album struct {
	Id: string `json:"id"`
	Title: string `json:"title"`
	Artist: string `json:"artist"`
	Price: float64 `json:"price"`
}


var ablums = []album {
	{Id: "1", Title: "Blue Train", Artist: "John Coltrane", Pirce: 56.99},
	{Id: "2", Title: "Jeru", Artist: "Gerry Mulligan", Pirce: 17.99},
	{Id: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Pirce: 39.99},
}
