package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App struct
type App struct {
	Router *mux.Router
	DB     *sql.DB
} // App struct

// Initialize 
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
} // Initialize

// Run
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
} // Run

// InitializeRoutes - set up routes for REST
func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/items", a.getItems).Methods("GET")
	a.Router.HandleFunc("/items", a.createItem).Methods("POST")
	a.Router.HandleFunc("/item/{id:[0-9]+}", a.getItem).Methods("GET")
	a.Router.HandleFunc("/item/{id:[0-9]+}", a.updateItem).Methods("PUT")
	a.Router.HandleFunc("/item/{id:[0-9]+}", a.deleteItem).Methods("DELETE")
} // InitializeRoutes

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
} // respondWithError

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
} // respondWithJSON

func (a *App) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	i := item{ID: id}
	if err := i.getItem(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Item not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, i)
} // getItem

func (a *App) getItems(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	items, err := getItems(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, items)
} // getItems

func (a *App) createItem(w http.ResponseWriter, r *http.Request) {
	var i item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid CREATE data")
		return
	}
	defer r.Body.Close()

	if err := i.createItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, i)
} // createItem

func (a *App) updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var i item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UPDATE data")
		return
	}
	defer r.Body.Close()
	i.ID = id

	if err := i.updateItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, i)
} // updateItem

func (a *App) deleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Item ID")
		return
	}

	i := item{ID: id}
	if err := i.deleteItem(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
} // deleteItem
