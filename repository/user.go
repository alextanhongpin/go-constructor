package repository

import "github.com/alextanhongpin/go-constructor/model"

type User struct {
	ID   int64
	Name string
	Age  int64
}

// NewUser maps repository User to model User.
func NewUser(u User) *model.User {
	return &model.User{
		ID:   u.ID,
		Name: u.Name,
		//Age:  u.Age,
	}
}
