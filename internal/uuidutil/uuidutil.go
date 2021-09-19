package uuidutil

import "github.com/google/uuid"

func RandomString() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uid.String()
}
