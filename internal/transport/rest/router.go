package rest

import (
	"demo-service/internal/transport/rest/handler/order"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewOrderRouter(handler *order.OrderHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/order/{order_uid}", handler.GetOrder)
	return r
}
