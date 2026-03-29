package main

import "github.com/MOQA-Solutions/tasker/types"
import (
    "bufio"
    "fmt"
    "net/http"
    "os"
    "strings"
	"log"
	"encoding/json"
	"bytes"
)

var (
	  api_key string
)

func init () {
	api_key = os.Getenv("api_key_1")
}

func main () {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), " ")
        switch parts[0] {
        case "new_task":
            newTask(parts)

        case "cancel_task":
            cancelTask(parts)

	    case "task_status":
            taskStatus(parts)

		case "pool_status":
            poolStatus()
        }
	}
}

func newTask(parts []string) {
	var req *http.Request
	switch parts[1] {
	case "email":
		taskType := "email"
		payload := types.EmailPayload {
                    Sender: "abdelghani",
                    Receiver: "hachemaoui",
	                Subject: "update from cloudflare",
	                Files: "https://github.com/MOQA-Solutions/Tasks.intro.pdf",
                }

		var jsonPayload []byte 
		var jsonRequest []byte
		var err error

	    jsonPayload, err = json.Marshal(payload)
		if err != nil {
			return
		}
		
		request := types.Request{
			            Type: taskType,
						Payload: jsonPayload,
		           }

		jsonRequest, err = json.Marshal(request)
	    if err != nil {
		  log.Fatal(err)
	    }

		req, err = http.NewRequest("POST", "http://localhost:9090/tasks/start?priority=high", bytes.NewBuffer(jsonRequest))
		if err != nil {
			log.Fatal(err)
		}

		sendRequest(req, "task")		

	case "image":
		taskType := "image"
		payload := types.ImagePayload {
                    Url: "https://images/image.png", 
					Processing: "a futuristic robot that simle to the front of screen",
                }

		var jsonPayload []byte 
		var jsonRequest []byte
		var err error

	    jsonPayload, err = json.Marshal(payload)
		if err != nil {
			return
		}
		
		request := types.Request{
			            Type: taskType,
						Payload: jsonPayload,
		           }

		jsonRequest, err = json.Marshal(request)
	    if err != nil {
		  log.Fatal(err)
	    }

		req, err = http.NewRequest("POST", "http://localhost:9090/tasks/start?priority=high", bytes.NewBuffer(jsonRequest))
		if err != nil {
			log.Fatal(err)
		}

		sendRequest(req, "task")	

	default: 
	  panic("Action not supported")
	}	
}

func cancelTask(parts []string) {
	id:= parts[1]
	
	req, err := http.NewRequest("Get", fmt.Sprintf("http://localhost:9090/tasks/cancel/%s", id), nil)
		if err != nil {
			log.Fatal(err)
		}
	sendRequest(req, "task")
}

func taskStatus(parts []string) {
	id:= parts[1]
	
	req, err := http.NewRequest("Get", fmt.Sprintf("http://localhost:9090/tasks/status/%s", id), nil)
		if err != nil {
			log.Fatal(err)
		}
	sendRequest(req, "task")
}

func poolStatus() {
	req, err := http.NewRequest("Get", "http://localhost:9090/pool/status", nil)
		if err != nil {
			log.Fatal(err)
		}
	sendRequest(req, "pool")
}

func sendRequest(req *http.Request, requestType string) {
	var taskResponse types.TaskResponse 
	var poolResponse types.PoolResponse

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", api_key) 
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if (requestType == "task") {
		err = json.NewDecoder(resp.Body).Decode(&taskResponse)
		if err != nil {
		  fmt.Printf("invalid json")
		  return
		  }
		fmt.Println(taskResponse) 
	} else {
		err = json.NewDecoder(resp.Body).Decode(&poolResponse)
		if err != nil {
		  fmt.Printf("invalid json")
		  return
		}
		fmt.Println(poolResponse)  
	}
}