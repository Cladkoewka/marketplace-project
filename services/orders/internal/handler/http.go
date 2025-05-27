package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Cladkoewka/marketplace-project/services/orders/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/metrics"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func SetupRouter(h *OrderHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(metrics.MetricsMiddleware)
	router.Handle("/metrics", promhttp.Handler())

	router.Post("/orders", h.create)
	router.Get("/orders/{id}", h.getByID)
	router.Get("/orders", h.getByCustomerID)
	router.Get("/orders/health", h.healthCheck)

	return router
}

func (h *OrderHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /healthz (orders) called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("orders ok"))
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST /orders called")
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		slog.Error("failed to decode order", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &order); err != nil {
		slog.Error("failed to create order", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("order created", slog.Any("order", order))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	slog.Info("GET /orders/{id} called", slog.String("id", idStr))

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("invalid order id", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetById(r.Context(), id)
	if err != nil {
		slog.Error("order not found", slog.Int64("id", id))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("order retrieved", slog.Any("order", order))
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) getByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	slog.Info("GET /orders with customer_id", slog.String("customer_id", customerIDStr))

	if customerIDStr == "" {
		http.Error(w, "missing customer_id", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		slog.Error("invalid customer_id", slog.String("customer_id", customerIDStr))
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetByCustomerId(r.Context(), customerID)
	if err != nil {
		slog.Error("failed to get orders by customer_id", slog.Int64("customer_id", customerID))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("orders retrieved", slog.Int("count", len(orders)))
	json.NewEncoder(w).Encode(orders)
}
