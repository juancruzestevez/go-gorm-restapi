package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juancruzestevez/go-gorm-restapi/db"
	"github.com/juancruzestevez/go-gorm-restapi/middleware"
	"github.com/juancruzestevez/go-gorm-restapi/models"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.DB.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
}


func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "User not found in context", http.StatusUnauthorized)
        return
    }

    var task models.Task
    json.NewDecoder(r.Body).Decode(&task)
    
    task.UserID = userID

    if err := db.DB.Create(&task).Error; err != nil {
        w.WriteHeader(http.StatusBadGateway)
        w.Write([]byte(err.Error()))
        return
    }

    json.NewEncoder(w).Encode(&task)
}


func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task models.Task
	db.DB.First(&task, params["id"])
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	json.NewEncoder(w).Encode(&task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task models.Task
	db.DB.First(&task, params["id"])
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}
	db.DB.Unscoped().Delete(&task)
	w.WriteHeader(http.StatusNoContent)
}