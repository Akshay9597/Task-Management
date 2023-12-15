package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strconv"
	"strings"
	"runtime/debug"
	"github.com/jmoiron/sqlx"
)

func (h *Handler) createPostgresDB(cfg DBConfig) (*sqlx.DB, error){
	fmt.Print()
	command := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.SSLMode, cfg.Password,
	)
	fmt.Print(command)

	db, err := sqlx.Connect("postgres", command)

	if(err != nil){
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if(err != nil){ // second one just to remove error at the moment
		// TODO: Log DB error
		fmt.Print(err.Error())
	}

	h.db = db

	return db, err
}

type createTaskResponse struct {
	Id int `json:"id"`
}

type fetchAllRecordsResponse struct {
	Records []Task `json:"tasks"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{ db: db}
}

func (h * Handler) handleAuth(ctx *gin.Context){
	auth_header := ctx.GetHeader("Authorization")

	header_parts := strings.Split(auth_header, " ");

	if len(header_parts) != 2 || header_parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"Incorrect Authorization header"})
		return
	}

	access_token, err := parseToken(header_parts[1])

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	}

	ctx.Set("user_id", access_token.UserId)


}

func (h *Handler) newTask( ctx *gin.Context){

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
			if(h.db == nil){
				fmt.Printf("DB is null")
				cfg, err := Init()
				if(err != nil){ 
					// TODO: Log host not specified
				}
				h.createPostgresDB(cfg)
				
			}
		}
	}()

	var request Task

	fmt.Print(request.Title)

	if err:= ctx.BindJSON(&request); err!=nil {
		// TODO: Log Bad request
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		fmt.Print(err.Error())
		return
	}

	var id int

	row:= h.db.QueryRow("INSERT into tasks (title,user_id) VALUES ($1, $2) RETURNING id", request.Title, request.UserId)

	if err:= row.Scan(&id); err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createTaskResponse{id})

}

func (h *Handler) fetchAllRecords(ctx *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
			if(h.db == nil){
				fmt.Printf("DB is null")
				cfg, err := Init()
				if(err != nil){ 
					// TODO: Log host not specified
				}
				h.createPostgresDB(cfg)
				
			}
		}
	}()
	var records [] Task
	err:= h.db.Select(&records, "SELECT * FROM tasks")

	if err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, fetchAllRecordsResponse{records})
}

func (h *Handler) getRecord(ctx *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
			if(h.db == nil){
				fmt.Printf("DB is null")
				cfg, err := Init()
				if(err != nil){ 
					// TODO: Log host not specified
				}
				h.createPostgresDB(cfg)
				
			}
		}
	}()
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"Invalid id"})
		return
	}

	var task Task
	err = h.db.Get(&task, "SELECT * FROM tasks WHERE id=$1", id)
	if err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
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
	routes.POST("/tasks",h.newTask)
	routes.GET("/tasks",h.fetchAllRecords)
	routes.GET("/tasks/:id",h.getRecord)
}