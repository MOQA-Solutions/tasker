package server

import "github.com/MOQA-Solutions/tasker/types"
import . "github.com/MOQA-Solutions/tasker/registry"
import(
  "database/sql"
  "time"
)

func Loop(high, medium, low <-chan types.WorkerEvent) {
    defer recover()
    defer Loop(high, medium, low)

    for {
        select {
          case workerEvent := <-high:
            handle(workerEvent)
            continue
          default:
        }

        select {
        case workerEvent := <-high:
            handle(workerEvent)
            continue
        default:
        }

        select {
        case workerEvent := <-low:
            handle(workerEvent)
            continue
        default:
        }

        Pool.IncrementPool()

        select {
        case workerEvent := <-high:
            handle(workerEvent)
        case workerEvent := <-medium:
            handle(workerEvent)
        case workerEvent := <-low:
            handle(workerEvent)
        }
    }
}

func handle(workerEvent types.WorkerEvent) {
  action := workerEvent.Action 
  switch action {
  case "start_task":
    startTask(workerEvent)
  }
}


func startTask(workerEvent types.WorkerEvent) {
    Pool.DecrementPool()
    state := "processing" 
    metadata := workerEvent.Metadata
    done := workerEvent.Done
    id := metadata.ID
    task := metadata.Task
    db := metadata.DB  
    payload := metadata.Payload

    _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
    if err != nil {
        return 
    }
  
    switch task {
      case "email":
        handleEmailTask(id, db, payload, done) 
      case "image": 
        handleImageTask(id, db, payload, done) 
    }
}

func handleEmailTask(id int, db *sql.DB, _payload any, done <-chan struct{}) {
  select {
    case <- done: 
      state := "canceled" 
      _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
      if err != nil {
        return 
      }
      CancelContext.DeleteCancelFunction(id)

    case <-time.After(15 * time.Second): 
      select {
      case <- done: 
        state := "canceled" 
        _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
        if err != nil {
          return 
        }

      default:
        state := "terminated"
        _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
          if err != nil {
              return 
          }
      
        cancel, exist := CancelContext.GetCancelFunction(id)
        if !exist {
          return
        }
        cancel()
        CancelContext.DeleteCancelFunction(id)
      }
    }
}


func handleImageTask(id int, db *sql.DB, _payload any, done <-chan struct{}) {
  select {
    case <- done: 
      state := "canceled" 
      _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
      if err != nil {
        return 
      }
      CancelContext.DeleteCancelFunction(id)

    case <-time.After(15 * time.Second): 
      select {
      case <- done: 
        state := "canceled" 
        _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
        if err != nil {
          return 
        }

      default:
        state := "terminated"
        _, err := db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", state, id)
          if err != nil {
              return 
          }
      
        cancel, exist := CancelContext.GetCancelFunction(id)
        if !exist {
          return
        }
        cancel()
        CancelContext.DeleteCancelFunction(id)
      }
    }
}
