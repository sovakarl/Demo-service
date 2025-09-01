package rest

import "demo-service/internal/service"

type OrderHandler struct {
	service service.Service
}

func NewOrderHandler(s service.Service) *OrderHandler{
	return &OrderHandler{service: s}
}
