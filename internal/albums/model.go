// Album struct + validation
package albums

import "errors"

type Album struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Artist string `json:"artist"`
	Price float64 `json:"price"`
}


func (a Album) Validate() error {
	if a.Title == "" {
		return errors.New("title is required")
	}
	if a.Artist == "" {
		return errors.New("artist is required")
	}
	if a.Price < 0 {
		return errors.New("price must be non-negative")
	}
	return nil
}
