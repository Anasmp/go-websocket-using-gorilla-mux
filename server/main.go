package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	Todo struct {
		ID          int    `json:"id,omitempty"`
		Description string `json:"description,omitempty"`
		Done        bool   `json:"done"`
	}
	Todos         []Todo
	clientRequest struct {
		Username string `json:"username,omitempty"`
		Type     string `json:"type,omitempty"`
		Todo     `json:"todo,omitempty"`
		ID       int `json:"id,omitempty"`
	}
	clientResponse struct {
		Todos `json:"todos,omitempty"`
	}
	Connections []*websocket.Conn
	Client      struct {
		Todos
		Connections
	}
)

var upgrader websocket.Upgrader
var db map[string]*Client
var todoID int

func main() {
	db = make(map[string]*Client)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		creq := &clientRequest{}
		err := conn.ReadJSON(creq)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("message from client", creq)
		if 0 == len(creq.Username) {
			return
		}
		cresp := &clientResponse{}
		var todos Todos
		var username = creq.Username
		switch creq.Type {
		case "hello":
			doLogin(username, conn)
			todos = getTodos(creq.Username)
		case "add":
			// fmt.Println(db["anas"])
			todos = addTodo(username, creq.Todo)
		case "delete":
			todos = removeTodo(username, creq.ID)
		case "toggle.done":
			todos = toggleDone(username, creq.ID)
		}

		cresp.Todos = todos
		connections := getConnections(username)
		for _, c := range connections {
			fmt.Println("Updating %v  clients for user %v", len(connections), username)
			// c.WriteJSON(cresp)
			if err := c.WriteJSON(cresp); err != nil {
				doLogout(username, c)
			}
		}

	}

}

func doLogout(username string, c *websocket.Conn) {
	conns := db[username].Connections
	var tmp Connections
	for _, v := range conns {
		if v != c {
			tmp = append(tmp, v)
		}
	}
	db[username].Connections = tmp
}

func getConnections(username string) Connections {
	return db[username].Connections
}

func doLogin(username string, c *websocket.Conn) {
	if db[username] == nil {
		db[username] = &Client{}
	}
	db[username].Connections = append(db[username].Connections, c)
}

func getTodos(username string) Todos {
	return db[username].Todos
}
func addTodo(username string, todo Todo) Todos {
	todoID++
	todo.ID = todoID
	db[username].Todos = append(db[username].Todos, todo)
	return db[username].Todos
}

func removeTodo(username string, id int) Todos {
	todos := db[username].Todos
	var tmp Todos
	for _, v := range todos {
		if id != v.ID {
			tmp = append(tmp, v)
		}
	}
	db[username].Todos = tmp
	return tmp
}

func toggleDone(username string, id int) Todos {
	// todos := db[username]

	for i, v := range db[username].Todos {
		if id == v.ID {
			db[username].Todos[i].Done = !v.Done
		}
	}
	return db[username].Todos
}
