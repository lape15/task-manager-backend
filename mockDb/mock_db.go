package mockDb

import (
	"database/sql"
	"task-manager-backend/types"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockDatabase struct {
	Db   *sql.DB
	Mock sqlmock.Sqlmock
}

func NewMockDatabase(t *testing.T) (*MockDatabase, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	return &MockDatabase{
			Db:   db,
			Mock: mock,
		}, func() {
			db.Close()
		}
}

func (m *MockDatabase) DoesTableExist(tableName string) (bool, error) {
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	m.Mock.ExpectQuery("SELECT EXISTS").WithArgs(tableName).WillReturnRows(rows)
	return true, nil
}

func (m *MockDatabase) CreateUsersTable() error {
	m.Mock.ExpectExec(`CREATE TABLE IF NOT EXISTS users`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	return nil
}

func (m *MockDatabase) CreateUser(user types.DbUser) (int64, error) {
	// rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	m.Mock.ExpectQuery(`INSERT INTO users\s*\(.+\) VALUES \(.+\)`).
		WithArgs(user.Username, user.PasswordHash, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	return 1, nil
}

func (m *MockDatabase) Close() {
	m.Db.Close()
}

func (m *MockDatabase) ScanDb(email, userName string) (types.DbUser, error) {
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash"}).
		AddRow(1, "testuser", "test@example.com", "hashedpassword")
	m.Mock.ExpectQuery("SELECT").WithArgs(email, userName).WillReturnRows(rows)
	return types.DbUser{
		Id:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}, nil
}

func (m *MockDatabase) DoesUserExist(email string, userName string) bool {
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	m.Mock.ExpectQuery("SELECT EXISTS").WithArgs(email, userName).WillReturnRows(rows)
	return true
}

func (m *MockDatabase) CreateTaskTable() error {
	m.Mock.ExpectExec("CREATE TABLE tasks").WillReturnResult(sqlmock.NewResult(0, 0))
	return nil
}

func (m *MockDatabase) CreateTask(task types.Task) (int64, error) {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	m.Mock.ExpectQuery("INSERT INTO tasks").WithArgs(task.Title, task.Description, task.Completed).WillReturnRows(rows)
	return 1, nil
}
