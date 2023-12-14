package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strconv"
	"runtime/debug"
	"github.com/jmoiron/sqlx"
	"github.com/Akshay9597/Task-Management/user-svc/internal/users"
	"github.com/Akshay9597/Task-Management/user-svc/internal/config"
)

func (h *Handler) createPostgresDB(cfg config.DBConfig) (*sqlx.DB, error){
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

type signupResponse struct {
	Id int `json:"id"`
}

// type fetchAllRecordsResponse struct {
// 	Records []tasks.Task `json:"tasks"`
// }

type errorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{ db: db}
}

func (h *Handler) signup( ctx *gin.Context){

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
			if(h.db == nil){
				fmt.Printf("DB is null")
				cfg, err := config.Init()
				if(err != nil){ 
					// TODO: Log host not specified
				}
				h.createPostgresDB(cfg)
				
			}
		}
	}()

	var request users.User

	fmt.Print(request.FirstName)

	if err:= ctx.BindJSON(&request); err!=nil {
		// TODO: Log Bad request
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		fmt.Print(err.Error())
		return
	}

	var id int

	row:= h.db.QueryRow("INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4) RETURNING id", request.FirstName, request.LastName, request.Username, request.Password)

	if err:= row.Scan(&id); err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, signupResponse{id})

}

func (h *Handler) getRecord(ctx *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("%v, %s", panicInfo, string(debug.Stack()))
			if(h.db == nil){
				fmt.Printf("DB is null")
				cfg, err := config.Init()
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

	var user users.User
	err = h.db.Get(&user, "SELECT id, first_name, last_name, username FROM users WHERE id=$1", id)
	if err != nil {
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
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
	routes.POST("/signup",h.signup)
	routes.GET("/users/:id",h.getRecord)
}