package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/samuelea/gator/internal/config"
	"github.com/samuelea/gator/internal/database"

	_ "github.com/lib/pq"
)

var cmds = config.Commands{
	Handlers: map[string]func(*config.State, config.Command) error{
		"login": loginHandler,
		"register": registerHandler,
		"reset": resetHandler,
		"users": listHandler,
	},
}

func main() {
	gatorConfig, err := config.Read()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", gatorConfig.DBUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	dbQueries := database.New(db)

	state := config.State{
		Config: gatorConfig,
		DbQueries: dbQueries,
	}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No arguments provided\n")
		os.Exit(1)
	}

	command := config.Command{
		Name: args[0],
		Args: args[1:],
	}

	err = cmds.Run(&state, command)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func loginHandler (state *config.State, command config.Command) error {
	if len(command.Args) == 0 {
		return errors.New("no username entered") 
	}
	if len(command.Args) > 1 {
		return errors.New("no the username cannot have spaces")
	}

	username := command.Args[0]

	err := loginUser(state, username)

	if err != nil {
		return err
	}

	return nil
}

func loginUser(state *config.State, username string) error {
	_, err := state.DbQueries.GetUser(context.Background(), username)

	if err != nil {
		return err
	}
	
	err = state.Config.SetUser(username)

	if err != nil {
		return err
	}

	fmt.Printf("username %s logged in!\n", username)

	return nil
} 

func registerHandler(state *config.State, command config.Command) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("no username provided")
	}

	_, err := state.DbQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: command.Args[0],
	})

	if err != nil {
		return fmt.Errorf("failed to create new user %s. user already exists", command.Args[0])
	}

	err = loginUser(state, command.Args[0])

	if err != nil {
		return err
	}

	return nil
}

func resetHandler(state *config.State, command config.Command) error {
	err := state.DbQueries.Reset(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func listHandler(state *config.State, command config.Command) error {
	users, err := state.DbQueries.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == state.Config.CurrentUserName {
			fmt.Printf("* %s (current) \n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
