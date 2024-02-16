package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all exciting workers
//
//	@Summary      List workers
//	@Description  get workers
//	@Tags         workers
//	@Produce      json
//	@Success      200  {array} []entity.Worker
//	@Failure      500  {object}  string
//	@Router       /workers [get]
func (h *Handler) getWorkers(c *gin.Context) {

	workers, err := h.services.Worker.UpdateWorkers(h.services.GetDelay())

	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(workers)

	newOkMessage(c, reqBodyBytes.String())
}
