package handler

import (
	"apps/engine/domain"
	"apps/engine/tools"
	"apps/engine/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UrlHttpHandler struct {
	domain *domain.UrlDomain
	ctx    context.Context
}

func NewHttpHandler(ctx context.Context, s *domain.UrlDomain) *UrlHttpHandler {
	return &UrlHttpHandler{
		domain: s,
		ctx:    ctx,
	}
}

func (h *UrlHttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id := r.PathValue("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(tools.NewBody(nil,
			"Missing id", tools.CODE_BAD_REQUEST),
		)
		return
	}

	p, err := h.domain.GetUrl(ctx, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, tools.ItemNotFoundError) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(tools.NewBody(nil,
				"Not found", tools.CODE_DB_ITEM_NOT_FOUND),
			)
		} else {
			fmt.Printf("error: %+v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(tools.NewBody(nil,
				"Something went wrong", tools.CODE_INTERNAL_SERVER_ERROR),
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tools.NewBody(p, "Ok", tools.CODE_OK))
}

func (h *UrlHttpHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id := r.PathValue("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(tools.NewBody(nil,
			"Missing id", tools.CODE_BAD_REQUEST),
		)
		return
	}

	err := h.domain.DeleteUrl(ctx, id)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(tools.NewBody(nil, "Error", tools.CODE_INTERNAL_SERVER_ERROR))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tools.NewBody(nil, "Ok", tools.CODE_OK))
}

func (h *UrlHttpHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	if r.Header.Get("Content-Type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(tools.NewBody(nil,
			"content type not supported", tools.CODE_BAD_REQUEST),
		)
		return
	}

	url := &types.UrlBase{}
	err := json.NewDecoder(r.Body).Decode(url)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(tools.NewBody(nil,
			"invalid json", tools.CODE_BAD_REQUEST),
		)
		return
	}

	err = h.domain.PutUrl(ctx, url)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
		if errors.Is(err, tools.InputValidationError) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(tools.NewBody(nil,
				"Validation error", tools.InputValidationError.Code),
			)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(tools.NewBody(nil, err.Error(), tools.CODE_INTERNAL_SERVER_ERROR))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tools.NewBody(nil, "Ok", tools.CODE_OK))
}
