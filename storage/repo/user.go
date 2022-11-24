package repo

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStorageI interface {
	Create(user *User) (*User, error)
	GetByEmail(email string) (*User, error)
}

type InMemoryStorageI interface {
	SetWithTTL(key string, value string, n int64) error
	Get(key string) (string, error)
}
