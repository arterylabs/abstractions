package builders

import "abstractions/db/models"

type UserBuilder struct {
	email    string
	name     string
	password string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) SetEmail(email string) *UserBuilder {
	b.email = email
	return b
}

func (b *UserBuilder) SetName(name string) *UserBuilder {
	b.name = name
	return b
}

func (b *UserBuilder) SetPassword(password string) *UserBuilder {
	b.password = password
	return b
}

func (b *UserBuilder) Build() models.User {
	user := models.User{
		Email: b.email,
		Name:  b.name,
	}

	return user
}
