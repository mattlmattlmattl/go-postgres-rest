package main

import "os"

// saving sensitive data in environment variables is a good way to keep
// them relatively safe.  The name of the database, the username and
// password are stored in the environment and accessed here.

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("GPR_DB_USERNAME"),
		os.Getenv("GPR_DB_PASSWORD"),
		os.Getenv("GPR_DB_NAME"))

	a.Run(":8080")
}
