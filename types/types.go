package types 

import (
  "sync"
  "encoding/json" 
  _ "github.com/jackc/pgx/v5/stdlib" 
  "database/sql"
  "context"
)

type Request struct {
    Type string          `json:"type"`
    Payload json.RawMessage `json:"data"`
}

type TaskResponse struct {
    ID     int    `json:"id"`
    State string `json:"state"`
}

type PoolResponse struct {
    Active     int    `json:"active"`
    Available int `json:"available"` 
    Capacity int `json:"capacity"`
}

type WorkerEvent struct {
  Action string 
  Metadata Metadata
  Done <-chan struct{}
}

type Metadata struct {
  ID int 
  Task string 
  DB *sql.DB 
  Payload any 
}

type EmailPayload struct {
    Sender  string `json:"sender"`
    Receiver string `json:"receiver"`
	  Subject string `json:"subject"`
	  Files string `json:"files"`
}

type ImagePayload struct {
    Url  string `json:"url"`
    Processing string `json:"processing"`
}

type Pool struct {
    active int
    mu sync.Mutex
}

type ProtectedChannel struct {
  ch chan WorkerEvent 
  mu sync.Mutex
}

type CancelContext struct {
  cancelFunctions map[int]context.CancelFunc
  mu sync.Mutex
}

/////////////////////////////////////////////////////////////////////////////////

func NewPool() *Pool {
  return &Pool{active: 0}
}

func (pool *Pool) IncrementPool() {
  pool.mu.Lock() 
  pool.active ++ 
  pool.mu.Unlock()
}

func (pool *Pool) DecrementPool() {
  pool.mu.Lock() 
  if pool.active != 0 {
    pool.active --
  } 
  pool.mu.Unlock()
}

func (pool *Pool) PoolActive() int {
  return pool.active
}

func NewProtectedChannel() *ProtectedChannel {
  return &ProtectedChannel{ch: make(chan WorkerEvent, 10000)}
}

func (protectedChannel *ProtectedChannel) GetChannel() chan WorkerEvent {
  return protectedChannel.ch 
}

func (protectedChannel *ProtectedChannel) GetSyncChannel() chan WorkerEvent {
  protectedChannel.mu.Lock() 
  return protectedChannel.ch 
}

func (protectedChannel *ProtectedChannel) ReturnSyncChannel() {
  protectedChannel.mu.Unlock()
} 

func NewCancelContext() *CancelContext {
  return &CancelContext{cancelFunctions: make(map[int]context.CancelFunc)}
}

func (cancelContext *CancelContext)GetCancelFunction(id int) (context.CancelFunc, bool) {
  val, exist := cancelContext.cancelFunctions[id] 
  return val, exist
}

func (cancelContext *CancelContext) AddCancelFunction(id int, cancelFunction context.CancelFunc) {
  cancelContext.mu.Lock()
  cancelContext.cancelFunctions[id] = cancelFunction 
  cancelContext.mu.Unlock()
}

func (cancelContext *CancelContext) DeleteCancelFunction(id int) {
  cancelContext.mu.Lock()
  delete(cancelContext.cancelFunctions, id)
  cancelContext.mu.Unlock()
}




