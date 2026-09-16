package httpdelivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"automation-test/orders/internal/delivery/rpc"
	"automation-test/orders/internal/domain"
	"automation-test/orders/internal/usecase"
	jsonhttp "automation-test/orders/pkg/httputil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	orders   *usecase.OrderUsecase
	payments *rpc.PaymentClient
}

func NewHandler(orders *usecase.OrderUsecase, payments *rpc.PaymentClient) *Handler {
	return &Handler{orders: orders, payments: payments}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/v1/orders", h.list)
	mux.HandleFunc("POST /api/v1/orders", h.create)
	mux.HandleFunc("GET /api/v1/orders/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/orders/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/orders/{id}", h.delete)
	return cors(mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "orders"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orders.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]any{"data": orders, "count": len(orders)})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateOrderInput
	if err := decodeJSON(w, r, &input); err != nil {
		jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	order, err := h.orders.Create(r.Context(), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusCreated, map[string]any{"data": order})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	order, err := h.orders.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	response := map[string]any{"data": order}
	if r.URL.Query().Get("include_payment") == "true" {
		payment, rpcErr := h.payments.GetByOrderID(r.Context(), order.ID)
		switch {
		case rpcErr == nil:
			response["payment"] = payment
		case status.Code(rpcErr) == codes.NotFound:
			response["payment"] = nil
		default:
			response["payment_unavailable"] = true
		}
	}
	jsonhttp.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateOrderInput
	if err := decodeJSON(w, r, &input); err != nil {
		jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	order, err := h.orders.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	jsonhttp.WriteJSON(w, http.StatusOK, map[string]any{"data": order})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.orders.Delete(r.Context(), r.PathValue("id")); err != nil {
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
