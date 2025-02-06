package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"task-manager-backend/db"
	"task-manager-backend/types"
)

type ResJson struct {
	Token   string       `json:"token"`
	Message string       `json:"message"`
	User    types.DbUser `json:"user,omitempty"`
}

func SignUp(res http.ResponseWriter, req *http.Request) {
	log.Println("Received sign-up request")

	var userReq types.User
	err := json.NewDecoder(req.Body).Decode(&userReq)
	if err != nil {
		log.Printf("Failed to decode request body: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if userReq.Email == "" || userReq.Password == "" {
		log.Println("Email and password are required")
		http.Error(res, "Email and password are required", http.StatusBadRequest)
		return
	}

	var userTab db.DatabaseInterface = db.DB
	if exists, _ := userTab.DoesTableExist("users"); !exists {
		err = userTab.CreateUsersTable()
		if err != nil {
			log.Printf("Failed to create users table: %v\n", err)
			http.Error(res, fmt.Sprintf("Failed to create users table: %v", err), http.StatusInternalServerError)
			return
		}
	}

	hashedPassword, err := HashPassword(userReq.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to hash password: %v", err), http.StatusInternalServerError)
		return
	}

	credential := types.Credential{
		Username:     userReq.Username,
		PasswordHash: hashedPassword,
	}
	userId, err := userTab.CreateUser(types.DbUser{
		Username:     userReq.Username,
		PasswordHash: hashedPassword,
		Email:        userReq.Email,
	})
	if err != nil {
		log.Printf("Failed to create user: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		res.Write([]byte("User already exists!"))
		return
	}
	credential.UserId = userId

	tokenString, err := GenerateJWT(credential)
	if err != nil {
		log.Printf("Failed to generate JWT: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to generate JWT: %v", err), http.StatusInternalServerError)
		return
	}

	resJson := ResJson{
		Token:   tokenString,
		Message: "User created successfully",
		User: types.DbUser{
			Id:       int(userId),
			Email:    userReq.Email,
			Username: userReq.Username,
		},
	}

	jsonRes, err := json.Marshal(resJson)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v\n", err)
		http.Error(res, fmt.Sprintf("Failed to marshal JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("User created successfully")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonRes)
}
