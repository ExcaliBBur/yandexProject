package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	entity "server/model/entity"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Update duration
//
//		@Summary      update duration
//		@Description  update duration
//		@Tags         duration
//		@Produce      json
//		@Success      200  {object} string
//	 	@Failure      400  {object}  string
//	 	@Failure      500  {object}  string
//		@Router       /duration [put]
func (h *Handler) putDuration(c *gin.Context) {
	var input entity.Duration

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	h.services.Duration.UpdateDuration(input)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(input)

	err := h.services.Duration.SendMessage(reqBodyBytes.Bytes())

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	newOkMessage(c, "updated")
}

// Get current duration
//
//	@Summary      get current duration
//	@Description  get current duration
//	@Tags         duration
//	@Produce      json
//	@Param request body entity.Duration true "query params"
//	@Success      200  {object} entity.Duration
//	@Failure      500  {object}  string
//	@Router       /duration [get]
func (h *Handler) getDuration(c *gin.Context) {
	duration, err := h.services.Duration.GetDuration()

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(duration)

	err = h.services.Duration.SendMessage(reqBodyBytes.Bytes())

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	newOkMessage(c, reqBodyBytes.String())
}
