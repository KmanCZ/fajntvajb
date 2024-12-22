package database

type User struct {
	ID          int
	Username    string
	DisplayName string
	Password    string
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := db.connection.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) CreateUser(username, displayName, password string) (*User, error) {
	user := &User{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	}
	_, err := db.connection.NamedExec("INSERT INTO users (username, display_name, password) VALUES (:username, :display_name, :password)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
