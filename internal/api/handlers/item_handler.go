package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"training/proj/internal/api/models"
	"training/proj/internal/customerrors"
	"training/proj/internal/db/repositories"

	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	ItemRepository *repositories.ItemRepository
}

func NewItemHandler(ir *repositories.ItemRepository) *ItemHandler {
	return &ItemHandler{
		ItemRepository: ir,
	}
}

func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, crudErr := h.ItemRepository.GetAll()

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "item_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	item, crudErr := h.ItemRepository.GetById(id)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	var itemReq models.Item

	decodeErr := json.NewDecoder(r.Body).Decode(&itemReq)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	itemResp, err := h.ItemRepository.Create(&itemReq)

	if err != nil {
		customerrors.ServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(itemResp)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "item_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	rowsAffecetd, crudErr := h.ItemRepository.Delete(id)

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	if rowsAffecetd == 0 {
		customerrors.NotFoundResponse(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) PutItem(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "item_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	var itemReq models.Item

	decodeErr := json.NewDecoder(r.Body).Decode(&itemReq)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	itemResp, crudErr := h.ItemRepository.Update(id, &itemReq)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(itemResp)
}

func (h *ItemHandler) GetItemCategories(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "item_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	categories, crudErr := h.ItemRepository.GetItemCategories(id)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
