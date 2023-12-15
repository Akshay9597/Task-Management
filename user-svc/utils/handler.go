package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"time"
	"context"
	"strconv"
	"strings"
	"runtime/debug"
	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
)

var b_ctx = context.Background()

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

func (h * Handler) handleAuth(ctx *gin.Context){
	auth_header := ctx.GetHeader("Authorization")
	fmt.Print(auth_header)

	header_parts := strings.Split(auth_header, " ");
	fmt.Print(header_parts)

	if len(header_parts) != 2 || header_parts[0] != "Bearer" {
		fmt.Print("This 1\n")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"Incorrect Authorization header"})
		return
	}

	access_token, err := parseToken(header_parts[1])

	if err != nil {
		fmt.Print("This 2\n")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	}

	ctx.Set("id", access_token.UserId)


}

func (h *Handler) signup( ctx *gin.Context){

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

	client := redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Your Redis server address
        Password: "",                // No password
        DB:       0,                 // Default DB
    })

	message := "Pusing to redis !"

	r_err := client.RPush(b_ctx, "logs", message).Err()

	if r_err != nil {
		fmt.Print("Error pushing log to Redis:", r_err.Error())
    }

	var request User

	fmt.Print(request.FirstName)

	if err:= ctx.BindJSON(&request); err!=nil {
		// TODO: Log Bad request
		client.RPush(b_ctx, "logs", "Signup: BadRequest")
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

	client.RPush(b_ctx, "logs", "Signup: Successful")

	ctx.JSON(http.StatusOK, signupResponse{id})

	client.Close()

}

type tokenInput struct {
	Username string `json:"username" binding:"required,min=5,max=20"`
	Password string `json:"password" binding:"required,min=8,max=30"`
}

type tokenResponse struct {
	AccessToken AccessToken `json:"access_token"`
}

func (h *Handler) generateNewTokenForUser(ctx *gin.Context){
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

	client := redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Your Redis server address
        Password: "",                // No password
        DB:       0,                 // Default DB
    })

	var input tokenInput

	if err:= ctx.BindJSON(&input); err!=nil {
		// TODO: Log Bad request
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		fmt.Print(err.Error())
		return
	}

	var user User

	err := h.db.Get(&user, "SELECT id, first_name, last_name, username  FROM users WHERE username=$1 AND password=$2", input.Username, input.Password)
	if err != nil {
		client.RPush(b_ctx, "logs", "Token: Invalid Login. Can't generate token")
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	client.RPush(b_ctx, "logs", "Token: Generated")

	ttl, err := time.ParseDuration("12h")

	token := generateNewToken(TokenInput{
		UserId: user.Id,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	})

	ctx.JSON(http.StatusOK, tokenResponse{token})

	client.Close()


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

	client := redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Your Redis server address
        Password: "",                // No password
        DB:       0,                 // Default DB
    })

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		client.RPush(b_ctx, "logs", "Profile: Invalid ID")
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"Invalid id"})
		return
	}

	var user User
	err = h.db.Get(&user, "SELECT id, first_name, last_name, username FROM users WHERE id=$1", id)
	if err != nil {
		client.RPush(b_ctx, "logs", "Profile: Incorrect Credentials")
		fmt.Print(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	client.RPush(b_ctx, "logs", "Profile: Fetch successful")
	ctx.JSON(http.StatusOK, user)
	client.Close()
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
	routes.GET("/token",h.generateNewTokenForUser)
	routes.GET("/users/:id",h.getRecord)
}