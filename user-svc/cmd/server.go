package main

import(
	"context"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
	"github.com/jmoiron/sqlx"
	"fmt"
	"github.com/Akshay9597/Task-Management/user-svc/internal/handlers"
	"github.com/Akshay9597/Task-Management/user-svc/internal/config"
	// postgres driver import
	_ "github.com/lib/pq"
)

type Server struct {
	httpServer *http.Server
}

func createPostgresDB(cfg config.DBConfig) (*sqlx.DB, error){
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

	return db, err
}

func createServer(handler http.Handler) *Server{
	return &Server{
		httpServer: &http.Server{
			Addr: ":8000",
			Handler: handler,
			ReadTimeout: 10*time.Second,
			WriteTimeout: 10*time.Second,
			MaxHeaderBytes: 1<<20,
		},
	}
}

func (s *Server) runServer() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) stopServer(ctx context.Context)  error {
	return s.httpServer.Shutdown(ctx)
}

func main(){

	cfg, err := config.Init()

	if(err != nil){
		// TODO: Log host not specified
	}


	db,error := createPostgresDB(cfg)
	// createPostgresDB()

	if(error != nil){ // second one just to remove error at the moment
		// TODO: Log DB error
		fmt.Print(error.Error())
	}

	handler := handlers.NewHandler(db)

	server := createServer(handler.Init())

	go func() {
		if err := server.runServer(); err != nil {
			// TODO: Log error
			
		}
	}()

	// TODO: Log Server started

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	// TODO: Log Server shutdown

	if err:= server.stopServer(context.Background()); err != nil {
		// TODO: Log Error shutting down server
	}

}