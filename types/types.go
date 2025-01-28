package types

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Credential struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password,omitempty"`
	UserId       int64  `json:"id,omitempty"`
}

type DbUser struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password,omitempty"`
	Email        string `json:"email"`
	Id           int    `json:"id,omitempty"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	User        int    `json:"user_id,omitempty"`
	Tags        []Tag  `json:"tags,omitempty"`
}
