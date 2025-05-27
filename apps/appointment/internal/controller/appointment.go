package controller

import (
	"appointment/internal/dto"
	"appointment/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	service service.AppointmentService
}

func (c *AppointmentController) RegisterRoutes(router *gin.RouterGroup) {
	appointments := router.Group("/appointments")
	{
		appointments.POST("", c.Create)
		appointments.GET("/:id", c.GetByID)
		appointments.DELETE("/:id", c.Delete)
		appointments.GET("/user/:user_id", c.GetByUserID)
		appointments.GET("/doctor/:doctor_id", c.GetByDoctorID)
		appointments.PATCH("/:id/status", c.ChangeStatus)
	}
}

func NewAppointmentController(service service.AppointmentService) *AppointmentController {
	return &AppointmentController{
		service: service,
	}
}

func (c *AppointmentController) Create(ctx *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *AppointmentController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	err = c.service.Delete(&dto.AppointmentIDRequest{ID: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *AppointmentController) GetByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	resp, err := c.service.GetByUserID(&dto.GetAppointmentsByUserIDRequest{UserID: userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AppointmentController) GetByDoctorID(ctx *gin.Context) {
	doctorID, err := uuid.Parse(ctx.Param("doctor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	resp, err := c.service.GetByDoctorID(&dto.GetAppointmentsByDoctorIDRequest{DoctorID: doctorID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AppointmentController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	resp, err := c.service.GetByID(&dto.AppointmentIDRequest{ID: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AppointmentController) ChangeStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment ID"})
		return
	}

	var req dto.ChangeAppointmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	resp, err := c.service.ChangeStatus(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
