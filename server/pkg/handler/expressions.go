package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get expression lists all existing expressions
//
//	@Summary      List expression
//	@Description  get expressions
//	@Tags         expression
//	@Produce      json
//	@Success      200  {array}   []dto.ExpressionResponse
//	@Failure      400  {object}  string
//	@Failure      500  {object}  string
//	@Router       /expressions [get]
func (h *Handler) getExpressions(c *gin.Context) {
	params := c.Request.URL.Query()

	pageNumber, err := strconv.Atoi(params["page_number"][0])

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	pageSize, err := strconv.Atoi(params["page_size"][0])

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	expressions, err := h.services.Expression.GetExpressions(pageNumber, pageSize)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(expressions)

	newOkMessage(c, reqBodyBytes.String())
}
