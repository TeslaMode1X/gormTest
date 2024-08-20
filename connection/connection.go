package connection

import (
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/gormTest/model"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var DB *gorm.DB
var err error

func init() {
	dsn := "host=localhost user=postgres password=admin123 dbname=gorm port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to DB", err)
	}
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Migration failed", err)
	}
}

func createUser(user model.User) {
	err = DB.Create(&user).Error
	if err != nil {
		log.Fatal(err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	createUser(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	fmt.Println("Created user: ", user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	res := DB.Find(&users)

	if res.Error != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	fmt.Println("All users")
}

func updateUser(user model.User) error {
	var oldUser model.User
	res := DB.First(&oldUser, user.ID)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return res.Error
		}
		return res.Error
	}

	if user.Username == "" {
		user.Username = oldUser.Username
	}
	if user.Age == 0 {
		user.Age = oldUser.Age
	}
	if user.Job == "" {
		user.Job = oldUser.Job
	}

	res = DB.Save(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = uint(id)

	err = updateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(id uint) error {
	// So, in gorm we have deleted_at row, because of it we can't simply delete it
	// We should use UnscopedDelete instead Delete
	// Why we are passing pointer to &model.User{}? It is for gorm to determine which table to use
	result := DB.Unscoped().Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = deleteUser(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
