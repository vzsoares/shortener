package handlers

import (
	"apps/engine/domain"
	"apps/engine/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"libs/utils"
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
		json.NewEncoder(w).Encode(utils.NewBody(nil,
			"Missing id", utils.CODE_BAD_REQUEST),
		)
		return
	}

	p, err := h.domain.GetUrl(ctx, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, utils.ItemNotFoundError) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(utils.NewBody(nil,
				"Not found", utils.CODE_DB_ITEM_NOT_FOUND),
			)
		} else {
			fmt.Printf("error: %+v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(utils.NewBody(nil,
				"Something went wrong", utils.CODE_INTERNAL_SERVER_ERROR),
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewBody(p, "Ok", utils.CODE_OK))
}

func (h *UrlHttpHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id := r.PathValue("id")

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.NewBody(nil,
			"Missing id", utils.CODE_BAD_REQUEST),
		)
		return
	}

	err := h.domain.DeleteUrl(ctx, id)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.NewBody(nil, "Error", utils.CODE_INTERNAL_SERVER_ERROR))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewBody(nil, "Ok", utils.CODE_OK))
}

func (h *UrlHttpHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	if r.Header.Get("Content-Type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.NewBody(nil,
			"content type not supported", utils.CODE_BAD_REQUEST),
		)
		return
	}

	url := &types.UrlBase{}
	err := json.NewDecoder(r.Body).Decode(url)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.NewBody(nil,
			"invalid json", utils.CODE_BAD_REQUEST),
		)
		return
	}

	err = h.domain.PutUrl(ctx, url)
	if err != nil {
		fmt.Printf("error: %+v", err.Error())
		if errors.Is(err, utils.InputValidationError) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(utils.NewBody(nil,
				"Validation error", utils.InputValidationError.Code),
			)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(utils.NewBody(nil, err.Error(), utils.CODE_INTERNAL_SERVER_ERROR))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewBody(nil, "Ok", utils.CODE_OK))
}
