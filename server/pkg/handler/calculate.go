package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"server/model/entity"

	"github.com/gin-gonic/gin"
)

// calculate expression
//
//	@Summary      Calculate expression
//	@Description  calculate expression
//	@Tags         calculate
//	@Produce      json
//	@Param request body dto.ExpressionRequest true "query params"
//	@Success      200  {object} string
//	@Failure      400  {object}  string
//	@Failure      500  {object}  string
//	@Router       /calculate [post]
func (h *Handler) calculate(c *gin.Context) {
	var input entity.Expression

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	if res, _ := regexp.MatchString("[/*+-]", input.Expression); !res {
		newErrorMessage(c, http.StatusBadRequest, "Expression must contain operators")
		return
	}
	if input.Expression == "" {
		newErrorMessage(c, http.StatusBadRequest, "Expression is empty")
		return
	}

	idempotency_key := c.Request.Header.Get("X-Request-ID")
	if idempotency_key == "" {
		newErrorMessage(c, http.StatusBadRequest, "X-Request-ID required")
		return
	}

	isExists, err := h.services.Idempotency.IsIdempotencyKeyExists(idempotency_key)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	var id int

	id, err = h.services.Idempotency.GetExpressionId(idempotency_key)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isExists {
		newOkMessage(c, fmt.Sprintf("Expression already exists, id: %d", id))
		return
	}

	_, err = h.services.Expression.ParseExpression(input)

	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Expression %s accepted for processing", input.Expression)

	if id, err = h.services.Expression.CreateExpression(input); err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.Id = id

	err = h.services.Idempotency.CreateIdempotencyKey(idempotency_key, id)

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
	}

	go func() {
		err = h.services.Expression.EvaluateAndUpdateExpression(input)

		if err != nil {
			newErrorMessage(c, http.StatusInternalServerError, err.Error())
		}
	}()

	newOkMessage(c, fmt.Sprint(id))
}
