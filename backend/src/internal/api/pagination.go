package api

import (
	"errors"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
)

type PaginationResponse struct {
	PageIndex  int64 `json:"pageIndex" description:"Индекс запрошенной страницы начиная с нуля" example:"0"`
	PageSize   int64 `json:"pageSize" description:"Размер запрашиваемой страницы" example:"50"`
	TotalPages int64 `json:"totalPages" description:"Общее кол-во страниц в лидерборде с учетом pageSize" example:"10"`
} //@name PaginationResponse

type paginationRequest struct {
	PageIndex *int64 `form:"pageIndex" binding:"required"`
	PageSize  *int64 `form:"pageSize" binding:"required"`
}

func newPaginationRequest(c *gin.Context) (*paginationRequest, error) {
	var req paginationRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if pointer.Get(req.PageIndex) < 0 {
		return nil, errors.New("pageIndex can not be negative")
	}

	if pointer.Get(req.PageSize) <= 0 {
		return nil, errors.New("pageSize can not be less or equal zero")
	}

	if pointer.Get(req.PageSize) > 100 {
		return nil, errors.New("pageSize can not be more than 100")
	}

	return &req, nil
}

func GetPaginationQueryParams(c *gin.Context) (pageIndex, pageSize int64, err error) {
	req, err := newPaginationRequest(c)
	if err != nil {
		return 0, 0, err
	}

	return pointer.Get(req.PageIndex), pointer.Get(req.PageSize), nil
}
