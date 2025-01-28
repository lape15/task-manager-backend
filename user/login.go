package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"task-manager-backend/db"
	"task-manager-backend/types"
)

func Login(res http.ResponseWriter, req *http.Request) {
	var user types.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		log.Println("Email and password are required")
		http.Error(res, "Email and password are required", http.StatusBadRequest)
		return
	}
	var userAction db.DatabaseInterface = db.DB

	dbUser, err := userAction.ScanDb(user.Email, user.Username)
	if err != nil {
		log.Printf("Failed to scan database: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to retreive user: %v", err), http.StatusInternalServerError)
		return
	}

	if dbUser == (types.DbUser{}) {
		log.Println("User not found")
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	err = VerifyPassword(dbUser.PasswordHash, user.Password)
	if err != nil {
		log.Println("Invalid password")
		http.Error(res, "Invalid password", http.StatusUnauthorized)
		return
	}

	credential := types.Credential{
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
	}

	tokenString, err := GenerateJWT(credential)
	if err != nil {
		log.Printf("Failed to generate JWT: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to generate JWT: %v", err), http.StatusInternalServerError)
		return
	}
	UserJson := types.DbUser{
		Email:    dbUser.Email,
		Id:       dbUser.Id,
		Username: dbUser.Username,
	}

	resJson := ResJson{
		Token:   tokenString,
		User:    UserJson,
		Message: "Login successfuL",
	}

	jsonRes, err := json.Marshal(resJson)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Login successfully done")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonRes)
}
