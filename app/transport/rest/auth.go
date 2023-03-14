package rest

import "os"

// isAuth checks credentials.
// It's generally recommended to store hashed passwords in a database.
// For the sake of simplicity environment variables address the problem without creating a new entity.
func isAuth(userName, password string) bool {
	return os.Getenv("USER") == userName && os.Getenv("PASSWORD") == password
}
