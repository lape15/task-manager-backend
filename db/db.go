package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"task-manager-backend/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Database struct {
	Db *sql.DB
}

var DB = &Database{}

// var DB DatabaseInterface

type DatabaseInterface interface {
	DoesTableExist(tableName string) (bool, error)
	CreateUsersTable() error
	CreateUser(user types.DbUser) (int64, error)
	Close()
	ScanDb(email, userName string) (types.DbUser, error)
	DoesUserExist(email string, userName string) bool
	CreateTaskTable() error
	CreateTask(task types.Task) (int64, error)
}

func (db *Database) Close() {
	if db.Db != nil {
		db.Db.Close()
	}
}

func (db *Database) DoesTableExist(tableName string) (bool, error) {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	var count int
	err := db.Db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		log.Printf("Error checking if table exists: %v", err)
		return false, err
	}

	return count > 0, nil

}

func (db *Database) CreateUsersTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) CreateUser(user types.DbUser) (int64, error) {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	result, err := db.Db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return 0, err
	}
	id, err := getlastId(result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *Database) CreateTaskTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id INT,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`
	_, err := db.Db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func getlastId(result sql.Result) (int64, error) {

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) CreateTask(task types.Task) (int64, error) {
	query := "INSERT INTO tasks (title, description,completed, user_id) VALUES (?, ?, ?, ?)"
	result, err := db.Db.Exec(query, task.Title, task.Description, task.Completed, task.User)
	if err != nil {
		return 0, err
	}
	id, err := getlastId(result)
	if err != nil {
		return 0, err
	}
	return id, nil

}
func (db *Database) DoesUserExist(email, userName string) bool {
	dbUser, _ := DB.ScanDb(email, userName)

	return dbUser.Email == ""
}

func (db *Database) ScanDb(email, userName string) (types.DbUser, error) {
	var user types.DbUser

	err := db.Db.QueryRow(
		"SELECT id, email, username, password_hash FROM users WHERE email = ? OR username = ? LIMIT 1",
		email, userName,
	).Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {

			return types.DbUser{}, nil
		}

		log.Printf("Error querying database: %v", err)
		return types.DbUser{}, err
	}

	return user, nil
}

func ConnectDb() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	// dbHost := "mysql-container" // mysql image container
	dbHost := "localhost"
	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Environment variables (MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE) are not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	var err error

	DB.Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.Db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected to MySQL!")

}
