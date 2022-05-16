package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"strconv"
)

type task struct {
	ID   int   `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type allTasks []task

var tasks = allTasks {
	{
		ID: 1,
		Name: "Task 1",
		Content: "This is the first task",
	},
}


func getTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getTasks")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createTask")
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getTask")
	vars := mux.Vars(r)

	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Kindly enter a valid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == key {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteTask")
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Kindly enter a valid ID")
		return
	}

	for index, task := range tasks {
		if task.ID == key {
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been deleted successfully", key)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateTask")
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Kindly enter a valid ID")
		return
	}

	var updatedTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &updatedTask)

	for index, task := range tasks {
		if task.ID == key {
			tasks = append(tasks[:index], tasks[index+1:]...)
			updatedTask.ID = key
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "The task with ID %v has been updated successfully", key)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, this is my first Go web app")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))

}