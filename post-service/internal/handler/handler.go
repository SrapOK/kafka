package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	Client *kafka.Writer
}

func (h *Handler) PostPerson(c *gin.Context) {
	var dto PersonDto

	if err := c.BindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, _ := json.Marshal(dto)

	msg := kafka.Message{
		Value: val,
	}

	if err := h.Client.WriteMessages(c.Request.Context(), msg); err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	c.String(http.StatusCreated, "Success")
}
