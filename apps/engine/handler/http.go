package handler

import (
	"apps/engine/domain"
	"apps/engine/tools"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type UrlHttpHandler struct {
	store *domain.UrlDomain
	ctx   context.Context
}

func NewHttpHandler(ctx context.Context, s *domain.UrlDomain) *UrlHttpHandler {
	return &UrlHttpHandler{
		store: s,
		ctx:   ctx,
	}
}

func (h *UrlHttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id := r.PathValue("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(tools.NewBody(nil, "Missing id", tools.CODE_BAD_REQUEST))
		return
	}

	p, err := h.store.GetUrl(ctx, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, tools.ItemNotFoundError) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(tools.NewBody(nil, "Not found", tools.CODE_DB_ITEM_NOT_FOUND))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(tools.NewBody(nil, "Not found", tools.CODE_INTERNAL_SERVER_ERROR))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tools.NewBody(p, "Ok", tools.CODE_OK))
}
