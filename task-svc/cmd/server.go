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
	"github.com/Akshay9597/Task-Management/task-svc/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func createPostgresDB() (*sqlx.DB, error){
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("host=db port=5432 user=postgres dbname=tasks sslmode=disable password=qwerty"),
	)

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

	db,error := createPostgresDB()

	if(error != nil && db == nil){ // second one just to remove error at the moment
		// TODO: Log DB error
		return
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