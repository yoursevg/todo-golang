package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Listening on port 8080")
}

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsComplete  bool   `json:"isComplete"`
}

var todos []Todo

func handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			if r.URL.Query().Get("id") != "" {
				handleGetTodo(w, r)
			} else {
				handleTodosGet(w, r)
			}
		}
	case http.MethodPost:
		handleTodosPost(w, r)
	case http.MethodPut:
		handleTodosUpdate(w, r)
	case http.MethodDelete:
		handleTodosDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Create New Todo
func handleTodosPost(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//Устанавливаем новый Id
	newTodo.Id = getNextID()

	//Добавляем элемент в наш слайс
	todos = append(todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&newTodo)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Get all todos
func handleTodosGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&todos)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Get Task by Id
func handleGetTodo(w http.ResponseWriter, r *http.Request) {

	// Получаем значение параметра "id" из URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	//Преобразуем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	// Используем полученный id для получения соответствующего Todo
	todo, _, err := findTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Отправляем полученный Todo в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Update existing Todo
func handleTodosUpdate(w http.ResponseWriter, r *http.Request) {

	// Получаем значение параметра "id" из URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	//Преобразуем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	//Находим нужный нам todo по айди
	todo, index, err := findTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Достаем из тела запроса новый Todo
	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	//Заменяем поля
	todo.Title = updatedTodo.Title
	todo.Description = updatedTodo.Description
	todo.IsComplete = updatedTodo.IsComplete

	todos[index] = *todo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&todo)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Delete existing Todo
func handleTodosDelete(w http.ResponseWriter, r *http.Request) {

	// Получаем значение параметра "id" из URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	//Преобразуем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	//Находим нужный нам todo по айди
	_, index, err := findTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Удаляем элемент из слайса
	todos = append(todos[:index], todos[index+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

func findTodoByID(id int) (*Todo, int, error) {
	for i, todo := range todos {
		if todo.Id == id {
			return &todo, i, nil
		}
	}
	return nil, -1, fmt.Errorf("todo with id %d not found", id)
}

func getNextID() int {
	if len(todos) == 0 {
		return 1
	}
	return todos[len(todos)-1].Id + 1
}

func main() {
	http.HandleFunc("/todos", handleTodos)
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
