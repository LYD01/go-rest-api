// HTTP Handlers e.g., (getAll, getById, create)
package albums

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LYD01/my-go-app/internal/respond"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.JSON(w, http.StatusOK, h.store.GetAll())
	}
}

func (h *Handler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		album, ok := h.store.GetById(id)
		if !ok {
			respond.Error(w, http.StatusNotFound, "album not found")
			return
		}
		respond.JSON(w, http.StatusOK, album)
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a Album
		if err := respond.Decode(r, &a); err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if a.Id == "" {
			a.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		if err := a.Validate(); err != nil {
			respond.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		h.store.Create(a)
		respond.JSON(w, http.StatusCreated, a)
	}
}
