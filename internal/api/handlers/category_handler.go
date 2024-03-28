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
	"github.com/jackc/pgerrcode"
)

type CategoryHandler struct {
	CategoryRepository     *repositories.CategoryRepository
	CategoryItemRepository *repositories.CategoryItemRepository
}

func NewCategoryHandler(cr *repositories.CategoryRepository, cir *repositories.CategoryItemRepository) *CategoryHandler {
	return &CategoryHandler{
		CategoryRepository:     cr,
		CategoryItemRepository: cir,
	}
}

func (h *CategoryHandler) PutCategoryItem(w http.ResponseWriter, r *http.Request) {
	categoryId, catConvErr := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)

	if catConvErr != nil {
		customerrors.BadRequestResponse(w, r, catConvErr)
		return
	}

	itemId, itemConvErr := strconv.ParseInt(chi.URLParam(r, "item_id"), 10, 64)

	if itemConvErr != nil {
		customerrors.BadRequestResponse(w, r, itemConvErr)
		return
	}

	crudErr := h.CategoryItemRepository.Create(categoryId, itemId)

	if crudErr != nil {
		switch crudErr.Code {
		case pgerrcode.UniqueViolation:
			customerrors.EditConflictResponse(w, r)
		case pgerrcode.ForeignKeyViolation:
			customerrors.NotFoundResponse(w, r)
		default:
			customerrors.ServerErrorResponse(w, r, crudErr)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, crudErr := h.CategoryRepository.GetAll()

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	category, crudErr := h.CategoryRepository.GetById(id)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) PostCategory(w http.ResponseWriter, r *http.Request) {
	var categoryReq models.Category

	decodeErr := json.NewDecoder(r.Body).Decode(&categoryReq)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	categoryResp, crudErr := h.CategoryRepository.Create(&categoryReq)

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoryResp)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	rowsAffecetd, crudErr := h.CategoryRepository.Delete(int64(id))

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

func (h *CategoryHandler) PutCategory(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	var categoryReq models.Category

	decodeErr := json.NewDecoder(r.Body).Decode(&categoryReq)

	if decodeErr != nil {
		customerrors.BadRequestResponse(w, r, decodeErr)
		return
	}

	categoryResp, crudErr := h.CategoryRepository.Update(id, &categoryReq)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categoryResp)
}

func (h *CategoryHandler) GetCategoryItems(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)

	if convErr != nil {
		customerrors.BadRequestResponse(w, r, convErr)
		return
	}

	items, crudErr := h.CategoryRepository.GetCategoryItems(id)

	if crudErr == sql.ErrNoRows {
		customerrors.NotFoundResponse(w, r)
		return
	}

	if crudErr != nil {
		customerrors.ServerErrorResponse(w, r, crudErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
