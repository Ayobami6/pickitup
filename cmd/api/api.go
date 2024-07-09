package api

import (
	"log"
	"net/http"

	"github.com/Ayobami6/pickitup/services/rider"
	"github.com/Ayobami6/pickitup/services/root"
	"github.com/Ayobami6/pickitup/services/user"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// create a new APIServer struct

type APIServer struct {
    addr string
	db *gorm.DB
}

// create a Global function ti instatiate a new APIServer

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
    return &APIServer{addr: addr, db:db}
}

// implement the Run method to start the server

func (a *APIServer) Run() error {
	// TODO: Implement the server setup and start logic here
    // create new mux router
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	rootHandler := root.NewRootHandler()
	rootHandler.RegisterRoutes(subrouter)
	userRepo := user.NewUserRepoImpl(a.db)
	userHandler := user.NewUserHandler(*userRepo)
	userHandler.RegisterRoutes(subrouter)
	// rider stuffs
	riderRepo := rider.NewRiderRepositoryImpl(a.db)
	// instantiate the rider handler
	riderHandler := rider.NewRiderHandler(*riderRepo)
    // register the rider routes
    riderHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on :", a.addr)

	// return http listening on the router
	return http.ListenAndServe(a.addr, router)
}