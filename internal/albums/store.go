// In-memory store (repository pattern)

/*
Go doesn't have a built-in constructor keyword. Instead, the idiomatic way to create and initialize an object is to write a regular function that starts with New... and returns a pointer to the struct.
*/

package albums

import "sync"

type Store struct {
	mu     sync.RWMutex
	albums []Album
}

func NewStore() *Store {
	return &Store{
		albums: []Album{
			{Id: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
			{Id: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
			{Id: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
		},
	}
}

/*
Create a GetAll() method
no args,
returns array of album,
unlock the data
defer the unlock
make result from data
copy result
return result
*/

func (s *Store) GetAll() []Album {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Album, len(s.albums))
	copy(result, s.albums)
	return  result

}
/*
Write a GetById method that return the album and true
lock the data
deffer the lock
loop through albums return true on id match, return false on miss
*/
func (s *Store) GetById(id string) (Album, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.albums {
		if id == a.Id {
			return a, true
		}
	}
	return Album{}, false
}


func (s *Store) Create(a Album) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.albums = append(s.albums, a)
}



