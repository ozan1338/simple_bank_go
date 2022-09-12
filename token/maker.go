package token

import "time"

//Maker is an interface managing token
type Maker interface {
	//Create Token creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string , error)

	//Verify Token checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}