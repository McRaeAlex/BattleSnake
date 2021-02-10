package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcraealex/BattleSnake2019/logic"
	"github.com/mcraealex/BattleSnake2019/structs"
)

// Handlers controls the routes of the project
type Handlers struct {
	logger *log.Logger
	//db     *database.DatabaseHandler
}

// NewHandlers returns a new handlers struct
func NewHandlers(logger *log.Logger /*,db *database.DatabaseHandler*/) *Handlers {
	return &Handlers{
		logger: logger,
		//db:     db,
	}
}

// SetupRoutes adds routes to a mux
func (h *Handlers) SetupRoutes(mux *mux.Router) {
	mux.HandleFunc("/ping", h.pingHandler).Methods("POST")
	mux.HandleFunc("/start", h.startHandler).Methods("POST")
	mux.HandleFunc("/move", h.moveHandler).Methods("POST")
	mux.HandleFunc("/end", h.endHandler).Methods("POST")
	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
}

func (h *Handlers) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Please hire me\n")
}

func (h *Handlers) startHandler(w http.ResponseWriter, r *http.Request) {
	// throw a error if no data is found
	// might make this middleware
	if r.Body == nil {
		http.Error(w, "Please send request body", 400)
		h.logger.Println("Error: No body on request")
		return
	}

	// this is the json they sent us. if i have enough time to add the database
	// I will keep track of wins vs losses
	var requestJSON structs.GeneralRequest
	json.NewDecoder(r.Body).Decode(&requestJSON)

	// create response
	response := structs.StartResponse{
		Color: "#9effb3",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TODO: Logic
func (h *Handlers) moveHandler(w http.ResponseWriter, r *http.Request) {
	// throw a error if no data is found
	// might make this middleware
	if r.Body == nil {
		http.Error(w, "Please send request body", 400)
		h.logger.Println("Error: No body on request")
		return
	}

	// this is the json they sent us. if i have enough time to add the database
	// I will keep track of wins vs losses
	var requestJSON structs.GeneralRequest
	err := json.NewDecoder(r.Body).Decode(&requestJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	// create response

	responseI := logic.Logic(requestJSON)
	var str string
	if responseI == 0 {
		str = "right"
	}
	if responseI == 1 {
		str = "down"
	}
	if responseI == 2 {
		str = "left"
	}
	if responseI == 3 {
		str = "up"
	}
	fmt.Printf("Answer: %s val: %v\n", str, responseI)
	h.logger.Println(str)
	response := structs.MoveResponse{
		Move: str,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) endHandler(w http.ResponseWriter, r *http.Request) {
	// throw a error if no data is found
	// might make this middleware
	if r.Body == nil {
		http.Error(w, "Please send request body", 400)
		h.logger.Println("Error: No body on request")
		return
	}

	// this is the json they sent us. if i have enough time to add the database
	// I will keep track of wins vs losses
	var requestJSON structs.GeneralRequest
	json.NewDecoder(r.Body).Decode(&requestJSON)
	// do something with the request

	w.WriteHeader(http.StatusOK)
}
