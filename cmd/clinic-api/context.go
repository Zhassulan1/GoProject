package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
)

type contextKey string

// userContextKey is used as a key for getting and setting user information in the request
// context.
const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request with the provided User struct added to the
// context.
func (app *application) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only time that
// this helper should be used is when we logically expect there to be a User struct value
// in the context, and if it doesn't exist it will firmly be an 'unexpected' error, upon we panic.
func (app *application) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)

	log.Print("\n\nIs there user\n\n")
	// log.Print(user)
	authToken := r.Header.Get("Authorization")

	if len(authToken) < 8 {
		log.Print("Error getting user \nAuth: ", authToken)
		return model.AnonymousUser
	}

	var err error
	if !ok {
		user, err = app.models.Users.GetUserByToken(authToken[7:])
		if err != nil {
			log.Print("Error getting user \nAuth: ", authToken)
			log.Println(err)
			log.Print("\n\n authToken[7:] = ", authToken[7:], "\n\n")
			panic("could not get user")
		}
	}

	log.Print("Auth: ", authToken)

	// if !ok {
	// 	panic("missing user value in request context")
	// }

	return user
}
