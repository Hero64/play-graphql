package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type User struct {
	Id       int
	Name     string
	LastName string
	Age      string
	Email    string
}


var users []User

func main() {
	u1 := User{
		Id:       1,
		Name:     "Aníbal",
		LastName: "Jorquera",
		Age:      "2016-01-01",
		Email:    "anibal.jorquera@outlook.com",
	}
	u2 := User{
		Id:       2,
		Name:     "Emilio",
		LastName: "Corejo",
		Age:      "2010-06-06",
		Email:    "anibal.jorquera@gmail.com",
	}

	users := append(users, u1, u2)

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Users",
		Fields: graphql.Fields{
			"Id": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"LastName": &graphql.Field{
				Type: graphql.String,
			},
			"Age": &graphql.Field{
				Type: graphql.String,
			},
			"Email": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	userMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "UserMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type:        userType,
				Description: "create new user",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lastName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					lastName, _ := params.Args["lastName"].(string)
					age, _ := params.Args["age"].(string)

					newUser := User{
						Id:       1,
						Name:     name,
						LastName: lastName,
						Age:      age,
					}
					users = append(users, newUser)
					return newUser, nil
				},
			},
		},
	})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "rootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        userType,
				Description: "get a single user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, ok := params.Args["id"].(int)
					if ok {
						for _, user := range users {

							if user.Id == id {
								return user, nil
							}
						}
					}
					return User{}, nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "All users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
		},
	})
	shema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: userMutation,
	})
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.POST("/api", func(c echo.Context) error {
		query := c.QueryParam("query")
		result := graphql.Do(graphql.Params{
			Schema:        shema,
			RequestString: query,
		})
		if len(result.Errors) > 0 {
			fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
		}
		return c.JSON(http.StatusOK, result)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
