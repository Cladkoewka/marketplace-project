package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/service"
	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func SetupRouter(h *CustomerHandler) http.Handler {
	router := chi.NewRouter()

	router.Post("/customers", h.create)
	router.Get("/customers/{id}", h.getByID)
	router.Get("/customers", h.getByEmailQueryParam)

	return router
}

func (h *CustomerHandler) create(w http.ResponseWriter, r *http.Request) {
	var c domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CustomerHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(c)
}

func (h *CustomerHandler) getByEmailQueryParam(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(c)
}
