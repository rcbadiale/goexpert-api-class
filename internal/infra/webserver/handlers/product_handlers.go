package handlers

import (
	"encoding/json"
	"goexpert-api/internal/dto"
	"goexpert-api/internal/entity"
	"goexpert-api/internal/infra/database"
	entityPkg "goexpert-api/pkg/entity"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductService database.ProductInterface
}

func NewProductHandler(service database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductService: service,
	}
}

// Create product godoc
// @Summary      Create a new product
// @Description  Create a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateProductInput true "product data"
// @Success      201
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      403
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /products [post]
// @Security     ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.ProductService.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: "error creating product"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateProductOutput{ID: p.ID.String()})
}

// Get product godoc
// @Summary      Get a product data
// @Description  Get a product data
// @Tags         products
// @Produce      json
// @Param        id       path      string true "product id"
// @Success      200      {object}  entity.Product
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      403
// @Failure      404      {object}  dto.ErrorOutput
// @Router       /products/{id} [get]
// @Security     ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}
	product, err := h.ProductService.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.ErrorOutput{Message: "product not found"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update product godoc
// @Summary      Update a product data
// @Description  Update a product data
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string true "product id"
// @Param        request  body      dto.CreateProductInput true "product data"
// @Success      200
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      403
// @Failure      404      {object}  dto.ErrorOutput
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /products/{id} [put]
// @Security     ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductService.Update(&product)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			error := dto.ErrorOutput{Message: "product not found"}
			json.NewEncoder(w).Encode(error)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: "server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete product godoc
// @Summary      Delete a product data
// @Description  Delete a product data
// @Tags         products
// @Produce      json
// @Param        id       path      string true "product id"
// @Success      200
// @Failure      400      {object}  dto.ErrorOutput
// @Failure      403
// @Failure      404      {object}  dto.ErrorOutput
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /products/{id} [delete]
// @Security     ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err := entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "invalid format"}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductService.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			error := dto.ErrorOutput{Message: "product not found"}
			json.NewEncoder(w).Encode(error)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: "server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Get all products godoc
// @Summary      Get all products data
// @Description  Get all products data
// @Tags         products
// @Produce      json
// @Param        page     query     string false "page number"
// @Param        limit    query     string false "limit"
// @Success      200      {array}   entity.Product
// @Failure      403
// @Failure      500      {object}  dto.ErrorOutput
// @Router       /products [get]
// @Security     ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}

	products, err := h.ProductService.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: "server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
