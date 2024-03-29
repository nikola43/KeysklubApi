package controllers

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/nikola43/KeysklubApi/utils"
	"github.com/panjf2000/ants"
	"log"
	"net/http"
)

var con bolt.Conn
var myPool *ants.Pool

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.Router.Use(JwtAuthentication) //attach JWT auth middleware
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	// configure allowed requests
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	// open db connection
	con = utils.CreateConnection()
	defer func() {
		_ = con.Close()
	}()

	// create pool
	myPool, _ = ants.NewPool(100000)
	defer myPool.Release()

	// run server
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(a.Router)))
}

func (a *App) initializeRoutes() {
	// USER
	a.Router.HandleFunc("/api/user/login", Login).Methods("POST")
	a.Router.HandleFunc("/api/user/new", SingUp).Methods("POST")
	a.Router.HandleFunc("/api/user", GetAll).Methods("GET")
	a.Router.HandleFunc("/api/user/{id:[0-9]+}", GetUserByID).Methods("GET")
	a.Router.HandleFunc("/api/user/{invite_id:[0-9]+}/{invited_id:[0-9]+}", InviteUser).Methods("GET")
	a.Router.HandleFunc("/api/user/{id:[0-9]+}/invited", GetNumberOfInvitedUsers).Methods("GET")
	a.Router.HandleFunc("/api/user/total", GetNumberOfUsers).Methods("GET")
}
