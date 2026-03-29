package main 

import "github.com/MOQA-Solutions/tasker/server"
import "github.com/MOQA-Solutions/tasker/types"
import . "github.com/MOQA-Solutions/tasker/registry"
import (
	"net/http"
	"database/sql"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"context"
  )

var (
     db *sql.DB
    )

func init() {
	var err error
	db, err = sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/tasks?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}	

    Workers = make(map[int][]*types.ProtectedChannel)
	for i:=0; i<1000; i++ {
	  high := types.NewProtectedChannel()
	  medium := types.NewProtectedChannel()
	  low := types.NewProtectedChannel()
	  Workers[i] = []*types.ProtectedChannel{high, medium, low}
	  go server.Loop(high.GetChannel(), medium.GetChannel(), low.GetChannel())
    }
}

func main() {
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/start", taskStartHandler)
	mux.HandleFunc("/tasks/status/{id}", taskStatusHandler)
	mux.HandleFunc("/tasks/cancel/{id}", taskCancelHandler)
	mux.HandleFunc("/pool/status", poolStatusHandler)

	handler := apiKeyMiddleware(mux)

	err := http.ListenAndServe(":9090", handler)
	fmt.Println("ListenAndServe error:", err)
}

func apiKeyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    key := r.Header.Get("X-API-Key")
        if key == "" {
          http.Error(w, "missing api key", http.StatusBadRequest)
          return
        } 

		var exists bool
		hash := server.HashAPIKey(key)
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE key = $1)", hash).Scan(&exists)
		
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) 
	}) 
}

func taskStartHandler(w http.ResponseWriter, r *http.Request) {
	priority := r.URL.Query().Get("priority")

    if priority != "high" && 
	   priority != "medium" && 
	   priority != "low" {
        http.Error(w, "priority error", http.StatusUnauthorized)
        return
      }

	key := r.Header.Get("X-API-Key")
	hash := server.HashAPIKey(key)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE key = $1", hash).Scan(&count)
	
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if count >= 2 {
		http.Error(w, "too much tasks for user", http.StatusUnauthorized)
		return
	}

	var request types.Request

	err = json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
      http.Error(w, "invalid json", http.StatusBadRequest)
      return
    } 

	switch request.Type {
    case "email":
		var email types.EmailPayload 
		err = json.Unmarshal(request.Payload, &email)
		handleTask(w, "email", email, hash, priority)

    case "image":
        var image types.ImagePayload 
		err = json.Unmarshal(request.Payload, &image)
		handleTask(w, "image", image, hash, priority)
	
	default:
		http.Error(w, "task not supported", http.StatusUnauthorized)
		return
	}

}

func taskStatusHandler(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
    if strID == "" {
        http.Error(w, "ID error", http.StatusUnauthorized)
        return
      }

	id, err := strconv.Atoi(strID)
	if err != nil {
        http.Error(w, "ID wrong format", http.StatusUnauthorized)
        return
      }

	var state string
	err = db.QueryRow("SELECT state FROM tasks WHERE id = $1", id).Scan(&state) 
	
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid Task", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error Database", http.StatusInternalServerError) 
		return
	} 

	writeTaskResponse(w, id, state)
}

func taskCancelHandler(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
    if strID == "" {
        http.Error(w, "ID error", http.StatusUnauthorized)
        return
      }

	id, err := strconv.Atoi(strID)
	if err != nil {
        http.Error(w, "ID wrong format", http.StatusUnauthorized)
        return
      }

	var state string
	err = db.QueryRow("SELECT state FROM tasks WHERE id = $1", id).Scan(&state) 
	
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid Task", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error Database", http.StatusInternalServerError) 
		return
	} 

	if state != "waiting" && state != "processing" {
		http.Error(w, "Task not processed in realtime", http.StatusUnauthorized)
		return
	}

	cancel, exist := CancelContext.GetCancelFunction(id) 
	if !exist {
		http.Error(w, "Task not processed in realtime", http.StatusUnauthorized) 
		return
	} 

	cancel()
	writeTaskResponse(w, id, state)
}

func poolStatusHandler(w http.ResponseWriter, _r *http.Request) {
	writePoolResponse(w)
}


func handleTask(w http.ResponseWriter, task string, payload any, hash string, priority string) {
    var id int
    state := "waiting"

    err := db.QueryRow(
      "INSERT INTO tasks (key, task, state) VALUES ($1, $2, $3) RETURNING id",
      hash, task, state,
      ).Scan(&id)
    if err != nil {
        log.Fatal(err)
    }

	channel := getChannel(id, priority)

	action := "start_task"
	metadata := types.Metadata{
		            ID: id,
					Task: task,
					DB: db, 
					Payload: payload,
	            } 
	ctx, cancel := context.WithCancel(context.Background()) 
	done := ctx.Done() 
	workerEvent := types.WorkerEvent {
		               Action: action, 
					   Metadata: metadata, 
					   Done: done,
	                } 

	channel <- workerEvent 

	CancelContext.AddCancelFunction(id, cancel)   	                   

    writeTaskResponse(w, id, state)
} 

func writeTaskResponse(w http.ResponseWriter, id int, state string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)

    json.NewEncoder(w).Encode(types.TaskResponse{
        ID: id,
        State: state,
    })
}

func writePoolResponse(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)

    json.NewEncoder(w).Encode(types.PoolResponse{
        Active: Pool.PoolActive(),
        Available: 1000 - Pool.PoolActive(), 
		Capacity: (Pool.PoolActive() * 100) / 1000,
    })
}

func getChannel(id int, priority string) chan types.WorkerEvent {
	index := server.Phash(strconv.Itoa(id), 1000)
	channels := Workers[index]
	switch priority {
	case "high":
      return channels[0].GetChannel()
	case "medium": 
	  return channels[1].GetChannel()
	case "low": 
	  return channels[2].GetChannel()
	default: 
	  panic("priority not defined")
	}
}



