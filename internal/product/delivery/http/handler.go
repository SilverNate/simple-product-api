package http

import (
	"github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"simple-product-api/internal/product"
	"simple-product-api/internal/product/usecase"
	"simple-product-api/pkg/common"
	validatorPkg "simple-product-api/pkg/validator"
)

type Handler struct {
	Usecase usecase.ProductUsecase
	Log     *logrus.Logger
}

func NewHandler(uc usecase.ProductUsecase, log *logrus.Logger) *Handler {
	return &Handler{Usecase: uc, Log: log}
}

func (h *Handler) Register(r fiber.Router) {
	r.Post("/", h.CreateProduct)
	r.Post("/list", h.ListProduct)
	r.Get("/:id", h.GetProductById)
}

// CreateProduct godoc
// @Summary Create products
// @Description Create product
// @Tags Products
// @Accept  json
// @Produce  json
// @Param   name query string true "Product Name"
// @Param   type query string true "Product Type"
// @Param   price query int false "Product Price"
// @Success 201 {object} common.Response
// @Failure 400 {object} common.Response
// @Router /api/v1/products [post]
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	h.Log.Info("received request to create products")

	var p product.Product
	if err := c.BodyParser(&p); err != nil {
		return common.BadRequest(c, err)
	}

	if err := validatorPkg.Validate.Struct(&p); err != nil {
		errs := map[string]string{}
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = e.ActualTag()
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(common.Response{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "validation failed",
			Data:    errs,
		})
	}

	if err := h.Usecase.CreateProduct(c.Context(), &p); err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err)
	}

	return common.Created(c, p, "product created successfully")
}

// ListProducts godoc
// @Summary Get list of products
// @Description Get paginated product list with filter & sort
// @Tags Products
// @Accept  json
// @Produce  json
// @Param   name query string false "Search query name"
// @Param   type query string false "Product type"
// @Param   sort_by query string false "Sort field"
// @Param   order query string false "Sort order"
// @Param   page query int false "Page number"
// @Param   limit query int false "Page size"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response
// @Router /api/v1/products/list [post]
func (h *Handler) ListProduct(c *fiber.Ctx) error {
	h.Log.Info("received request to list products")

	filter := product.ListFilter{
		Query:    c.Query("name"),
		Type:     c.Query("type"),
		SortBy:   c.Query("sort_by"),
		Order:    c.Query("order"),
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("limit", 10),
	}
	list, total, err := h.Usecase.ListProduct(c.Context(), filter)
	if err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err)
	}

	meta := product.MetaPage{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: (total + filter.PageSize - 1) / filter.PageSize,
	}

	return common.Success(c, list, "successfully fetched products", meta)
}

// GetProductById godoc
// @Summary Get products by id
// @Description Get product by using id
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} common.Response
// @Failure 404 {object} common.Response
// @Router /api/v1/products [get]
func (h *Handler) GetProductById(c *fiber.Ctx) error {
	h.Log.Info("received request get product by id")

	id := c.Params("id")
	result, err := h.Usecase.GetProductByID(c.Context(), id)
	if err != nil {
		return common.NotFound(c, err)
	}

	return common.Success(c, result, "successfully fetched products")
}
