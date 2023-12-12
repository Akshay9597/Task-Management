package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/Akshay9597/Task-Management/task-svc/internal/tasks"
)

type createTaskResponse struct {
	Id int `json:"id"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) createTask( ctx *gin.Context){
	var request tasks.Task

	if err:= ctx.BindJSON(&request); err!=nil {
		// TODO: Log Bad request
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	var id int

	row:= h.db.QueryRow("INSERT into tasks (title,user_id) VALUES ($1, $2) RETURNING id", request.Title, request.UserId)

	if err:= row.Scan(&id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createTaskResponse{id})

}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(ctx *gin.Context) {ctx.String(http.StatusOK, "pong")})

	h.configureV1Routes(router)

	return router
}

func (h *Handler) configureV1Routes(router *gin.Engine){
	routes := router.Group("/api/v1")
	routes.POST("/tasks",h.createTask)
}