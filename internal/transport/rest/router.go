package rest

import (
	"github.com/go-chi/chi/v5"
)

func NewOrderRouter(handler *OrderHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/order", handler.Page)
	r.Get("/order/{order_uid}", handler.GetOrder)
	return r
}
