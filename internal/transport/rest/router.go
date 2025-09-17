package rest

import (
	"demo-service/internal/transport/rest/handler/order"

	"github.com/go-chi/chi/v5"
)

func NewOrderRouter(handler *order.OrderHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/order", handler.Page)
	r.Get("/order/{order_uid}", handler.GetOrder)
	return r
}
