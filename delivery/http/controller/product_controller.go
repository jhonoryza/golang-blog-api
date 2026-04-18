package controller

import (
	"api_blog/delivery/http/response"
	"api_blog/infrastructure/exception"
	"api_blog/usecase/product"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProductController struct {
	Usecase *product.ProductUsecase
}

func NewProductController(uc *product.ProductUsecase) *ProductController {
	return &ProductController{Usecase: uc}
}

func (c *ProductController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	products := c.Usecase.FindAll()
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    products,
	}

	resp.ToJson(w)
}

func (c *ProductController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productId, err := strconv.Atoi(p.ByName("productId"))
	if err != nil {
		exception.BadRequestError(w, r, "BAD REQUEST")
		return
	}
	product := c.Usecase.FindOneById(productId)
	if product == nil {
		exception.BadRequestError(w, r, "BAD REQUEST")
		return
	}

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    product,
	}

	resp.ToJson(w)
}