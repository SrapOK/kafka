package handler

import (
	"get-service/internal/model/person"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *person.PersonRepo
}

func (h *Handler) GetPerson(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	ps, err := h.Repo.GetPerson(&ctx, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *ps)
}

func (h *Handler) GetRandomPerson(c *gin.Context) {
	ctx := c.Request.Context()

	ps, err := h.Repo.GetRandomPerson(&ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *ps)
}

func (h *Handler) GetAllPersons(c *gin.Context) {
	ctx := c.Request.Context()

	ps, err := h.Repo.GetAllPersons(&ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ps)
}
