package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/FollG/kafka-with-go/internal/domain/models"
	"github.com/FollG/kafka-with-go/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProductHandler struct {
	productUC *usecases.ProductUseCase
}

func NewProductHandler(productUC *usecases.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUC: productUC,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
		return
	}

	product := &models.Product{
		Name:   req.Name,
		Weight: req.Weight,
		Unit:   req.Unit,
		Color:  req.Color,
		Type:   models.ProductType(req.Type),
		Price:  req.Price,
		Attributes: models.Attributes{
			Size:               req.Attributes.Size,
			HeadCircumference:  req.Attributes.HeadCircumference,
			ChestCircumference: req.Attributes.ChestCircumference,
			WaistCircumference: req.Attributes.WaistCircumference,
			HipCircumference:   req.Attributes.HipCircumference,
			FootSize:           req.Attributes.FootSize,
			ExpiryDate:         req.Attributes.ExpiryDate,
			NutritionalInfo:    req.Attributes.NutritionalInfo,
			WarrantyMonths:     req.Attributes.WarrantyMonths,
			Voltage:            req.Attributes.Voltage,
			Dimensions:         req.Attributes.Dimensions,
			Material:           req.Attributes.Material,
		},
	}

	if err := h.productUC.CreateProduct(ctx, product); err != nil {
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create product",
			})
		}
		return
	}

	render.Status(r, http.StatusAccepted) // 202 Accepted - операция принята в обработку
	render.JSON(w, r, CreateProductResponse{
		ID:      product.ID,
		Message: "Product creation accepted",
		Status:  "processing",
	})
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid product ID",
		})
		return
	}

	product, err := h.productUC.GetProduct(ctx, id)
	if err != nil {
		if err == models.ErrProductNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get product",
		})
		return
	}

	if product == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	render.JSON(w, r, ProductResponse{
		ID:     product.ID,
		Name:   product.Name,
		Weight: product.Weight,
		Unit:   product.Unit,
		Color:  product.Color,
		Type:   string(product.Type),
		Price:  product.Price,
		Attributes: AttributesResponse{
			Size:               product.Attributes.Size,
			HeadCircumference:  product.Attributes.HeadCircumference,
			ChestCircumference: product.Attributes.ChestCircumference,
			WaistCircumference: product.Attributes.WaistCircumference,
			HipCircumference:   product.Attributes.HipCircumference,
			FootSize:           product.Attributes.FootSize,
			ExpiryDate:         product.Attributes.ExpiryDate,
			NutritionalInfo:    product.Attributes.NutritionalInfo,
			WarrantyMonths:     product.Attributes.WarrantyMonths,
			Voltage:            product.Attributes.Voltage,
			Dimensions:         product.Attributes.Dimensions,
			Material:           product.Attributes.Material,
		},
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := models.ProductFilter{
		Limit:  25, // дефолтный лимит
		Offset: 0,
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil && minPrice >= 0 {
			filter.MinPrice = &minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil && maxPrice >= 0 {
			filter.MaxPrice = &maxPrice
		}
	}

	if color := r.URL.Query().Get("color"); color != "" {
		filter.Color = color
	}

	if types := r.URL.Query()["type"]; len(types) > 0 {
		filter.Types = make([]models.ProductType, len(types))
		for i, t := range types {
			filter.Types[i] = models.ProductType(t)
		}
	}

	products, err := h.productUC.ListProducts(ctx, filter)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list products",
		})
		return
	}

	response := ListProductsResponse{
		Products: make([]ProductResponse, len(products)),
		Total:    len(products),
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}

	for i, product := range products {
		response.Products[i] = ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Weight: product.Weight,
			Unit:   product.Unit,
			Color:  product.Color,
			Type:   string(product.Type),
			Price:  product.Price,
			Attributes: AttributesResponse{
				Size:               product.Attributes.Size,
				HeadCircumference:  product.Attributes.HeadCircumference,
				ChestCircumference: product.Attributes.ChestCircumference,
				WaistCircumference: product.Attributes.WaistCircumference,
				HipCircumference:   product.Attributes.HipCircumference,
				FootSize:           product.Attributes.FootSize,
				ExpiryDate:         product.Attributes.ExpiryDate,
				NutritionalInfo:    product.Attributes.NutritionalInfo,
				WarrantyMonths:     product.Attributes.WarrantyMonths,
				Voltage:            product.Attributes.Voltage,
				Dimensions:         product.Attributes.Dimensions,
				Material:           product.Attributes.Material,
			},
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		}
	}

	render.JSON(w, r, response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid product ID",
		})
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
		return
	}

	product := &models.Product{
		ID:     id,
		Name:   req.Name,
		Weight: req.Weight,
		Unit:   req.Unit,
		Color:  req.Color,
		Type:   models.ProductType(req.Type),
		Price:  req.Price,
		Attributes: models.Attributes{
			Size:               req.Attributes.Size,
			HeadCircumference:  req.Attributes.HeadCircumference,
			ChestCircumference: req.Attributes.ChestCircumference,
			WaistCircumference: req.Attributes.WaistCircumference,
			HipCircumference:   req.Attributes.HipCircumference,
			FootSize:           req.Attributes.FootSize,
			ExpiryDate:         req.Attributes.ExpiryDate,
			NutritionalInfo:    req.Attributes.NutritionalInfo,
			WarrantyMonths:     req.Attributes.WarrantyMonths,
			Voltage:            req.Attributes.Voltage,
			Dimensions:         req.Attributes.Dimensions,
			Material:           req.Attributes.Material,
		},
	}

	if err := h.productUC.UpdateProduct(ctx, product); err != nil {
		switch {
		case err == models.ErrProductNotFound:
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
		case strings.Contains(err.Error(), "validation failed"):
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to update product",
			})
		}
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, UpdateProductResponse{
		Message: "Product update accepted",
		Status:  "processing",
	})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid product ID",
		})
		return
	}

	if err := h.productUC.DeleteProduct(ctx, id); err != nil {
		if err == models.ErrProductNotFound {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete product",
		})
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, DeleteProductResponse{
		Message: "Product deletion accepted",
		Status:  "processing",
	})
}
