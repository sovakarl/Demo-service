package rest

import (
	"demo-service/internal/transport/rest/handler/order"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewOrderRouter(handler *order.OrderHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/order", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("потом доделаю")) })
	r.Get("/order/{order_uid}", handler.GetOrder)
	return r
}
