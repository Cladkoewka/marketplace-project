package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/metrics"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func SetupRouter(h *CustomerHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(metrics.MetricsMiddleware)
	router.Handle("/metrics", promhttp.Handler())

	router.Post("/customers", h.create)
	router.Get("/customers/{id}", h.getByID)
	router.Get("/customers", h.getByEmailQueryParam)
	router.Get("/customers/health", h.healthCheck)

	return router
}

func (h *CustomerHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /healthz (orders) called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("customers ok"))
}

func (h *CustomerHandler) create(w http.ResponseWriter, r *http.Request) {
	slog.Info("POST /customers called")
	var c domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		slog.Error("failed to decode customer", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &c); err != nil {
		slog.Error("failed to create customer", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("customer created", slog.Any("customer", c))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CustomerHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	slog.Info("GET /customers/{id} called", slog.String("id", idStr))

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("invalid customer id", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("customer not found", slog.Int64("id", id))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("customer retrieved", slog.Any("customer", c))
	json.NewEncoder(w).Encode(c)
}

func (h *CustomerHandler) getByEmailQueryParam(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	slog.Info("GET /customers with email", slog.String("email", email))

	if email == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		slog.Error("customer not found", slog.String("email", email))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("customer retrieved", slog.Any("customer", c))
	json.NewEncoder(w).Encode(c)
}
