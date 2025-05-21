package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Cladkoewka/marketplace-project/services/orders/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/service"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func SetupRouter(h *OrderHandler) http.Handler {
	router := chi.NewRouter()
	router.Post("/orders", h.create)
	router.Get("/orders/{id}", h.getByID)
	router.Get("/orders", h.getByCustomerID)

	return router
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getByCustomerID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	customerIDStr := query.Get("customer_id")
	if customerIDStr == "" {
		http.Error(w, "missing customer_id", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetByCustomerId(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
