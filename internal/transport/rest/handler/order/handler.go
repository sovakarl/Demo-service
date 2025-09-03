package order

import (
	"demo-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type OrderHandler struct {
	service service.Service
}

func NewOrderHandler(s service.Service) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	if orderUID == "" {
		http.Error(w, `{"error": "order_uid is missing in path"}`, http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(orderUID)
	if err != nil {
		//TODO ДОДЕЛАТЬ ЭТУ ХУЙНЮ
		return
	}

	// Устанавливаем заголовок и отдаём JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, `{"error": "failed to serialize order"}`, http.StatusInternalServerError)
		return
	}
}
