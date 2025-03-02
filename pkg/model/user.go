package model

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User ユーザーデータを表す構造体
type User struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}

func GenerateUsers(numberOfUsers int) []User {
	var users []User
	for i := range numberOfUsers {
		users = append(users, User{
			ID:   bson.NewObjectID(),
			Name: fmt.Sprintf("user%v", i),
		})
	}
	return users
}
