package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager-backend/db"
	"task-manager-backend/mockDb"
	"task-manager-backend/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSignUp(t *testing.T) {

	userReq := types.DbUser{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "password123",
	}
	reqBody, err := json.Marshal(userReq)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	mockDB, cleanup := mockDb.NewMockDatabase(t)
	defer cleanup()

	if mockDB.Db == nil {
		t.Fatal("mockDB.Db is nil")
	}

	mockDB.Mock.ExpectQuery(`SELECT COUNT\(\*\) FROM information_schema.tables WHERE table_name = \?`).
		WithArgs("users").
		WillReturnRows(mockDB.Mock.NewRows([]string{"count"}).AddRow(0)) // Simulating table does not exist

	mockDB.Mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mockDB.Mock.ExpectExec(`INSERT INTO users \(username, email, password_hash\) VALUES \(\?, \?, \?\)`).
		WithArgs(userReq.Username, userReq.Email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	oldDB := db.DB.Db
	db.DB.Db = mockDB.Db
	defer func() { db.DB.Db = oldDB }()

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SignUp)

	// Serve the HTTP request to the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Decode the actual response
	var actualRes ResJson
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check if the token is present
	if actualRes.Token == "" {
		t.Error("Token was not generated in the response")
	}

	// Check the response message
	if actualRes.Message != "User created successfully" {
		t.Errorf("Unexpected response message: got %v, want %v", actualRes.Message, "User created successfully")
	}

	// Check the user ID
	if actualRes.User.Id != 1 {
		t.Errorf("Unexpected user ID: got %v, want %v", actualRes.User.Id, 1)
	}

	// Check the username and email
	if actualRes.User.Username != userReq.Username {
		t.Errorf("Unexpected username: got %v, want %v", actualRes.User.Username, userReq.Username)
	}
	if actualRes.User.Email != userReq.Email {
		t.Errorf("Unexpected email: got %v, want %v", actualRes.User.Email, userReq.Email)
	}

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Verify all expectations were met
	if err := mockDB.Mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
