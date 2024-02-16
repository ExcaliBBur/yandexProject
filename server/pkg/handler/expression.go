package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get expression by id
//
//	@Summary      get expression by id
//	@Description  get expression by id
//	@Tags         expression
//	@Produce      json
//	@Param 		  id   path 	int 	true "id"
//	@Success      200  {object} dto.ExpressionResponse
//	@Failure      400  {object}  string
//	@Failure      500  {object}  string
//	@Router       /expression/{id} [get]
func (h *Handler) getExpression(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	expression, err := h.services.Expression.GetExpression(id)

	if expression.Id == 0 {
		newErrorMessage(c, http.StatusBadRequest, "The specified ID was not found")
		return
	}

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(expression)

	newOkMessage(c, reqBodyBytes.String())
}
