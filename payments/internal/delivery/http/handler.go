package httpdelivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"automation-test/payments/internal/delivery/rpc"
	"automation-test/payments/internal/domain"
	"automation-test/payments/internal/usecase"
	jsonhttp "automation-test/payments/pkg/httputil"
)

type Handler struct {
	payments *usecase.PaymentUsecase
	orders   *rpc.OrderClient
}

func NewHandler(payments *usecase.PaymentUsecase, orders *rpc.OrderClient) *Handler {
	return &Handler{payments: payments, orders: orders}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/v1/payments", h.list)
	mux.HandleFunc("POST /api/v1/payments", h.create)
	mux.HandleFunc("GET /api/v1/payments/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/payments/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/payments/{id}", h.delete)
	return cors(mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "payments"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	payments, err := h.payments.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]any{"data": payments, "count": len(payments)})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePaymentInput
	if err := decodeJSON(w, r, &input); err != nil {
		jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	payment, err := h.payments.Create(r.Context(), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusCreated, map[string]any{"data": payment})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	payment, err := h.payments.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response := map[string]any{"data": payment}
	if r.URL.Query().Get("include_order") == "true" {
		order, rpcErr := h.orders.GetByID(r.Context(), payment.OrderID)
		if rpcErr == nil {
			response["order"] = order
		} else {
			response["order_unavailable"] = true
		}
	}
	jsonhttp.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdatePaymentInput
	if err := decodeJSON(w, r, &input); err != nil {
		jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	payment, err := h.payments.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]any{"data": payment})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.payments.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		jsonhttp.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidID), errors.Is(err, domain.ErrInvalidInput):
		jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
	default:
		jsonhttp.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, ", "))
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
