package order

import (
	"demo-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"database/sql"
)

type OrderHandler struct {
	service service.Service
}

func NewOrderHandler(s service.Service) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderUID := chi.URLParam(r, "order_uid")

	if orderUID == "" {
		writeErrorJSON(w, http.StatusBadRequest, "order_uid is missing in path")
		return
	}

	order, err := h.service.GetOrder(ctx, orderUID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeErrorJSON(w, http.StatusNotFound, "order not found")
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "internal/web/order.html")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	writeJSON(w, status, errorResponse{Error: message})
}