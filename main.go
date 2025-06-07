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

type State struct {
	Config *config.Config
	DbQueries *database.Queries
}

type Command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*State, Command) error
}

var cmds = commands{
	handlers: map[string]func(*State, Command) error{
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

	state := State{
		Config: gatorConfig,
		DbQueries: dbQueries,
	}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No arguments provided\n")
		os.Exit(1)
	}

	command := Command{
		name: args[0],
		args: args[1:],
	}

	err = cmds.run(&state, command)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func loginHandler (state *State, command Command) error {
	if len(command.args) == 0 {
		return errors.New("no username entered") 
	}
	if len(command.args) > 1 {
		return errors.New("no the username cannot have spaces")
	}

	username := command.args[0]

	err := loginUser(state, username)

	if err != nil {
		return err
	}

	return nil
}

func loginUser(state *State, username string) error {
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

func registerHandler(state *State, command Command) error {
	if len(command.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	_, err := state.DbQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: command.args[0],
	})

	if err != nil {
		return fmt.Errorf("failed to create new user %s. user already exists", command.args[0])
	}

	err = loginUser(state, command.args[0])

	if err != nil {
		return err
	}

	return nil
}

func resetHandler(state *State, command Command) error {
	err := state.DbQueries.Reset(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func listHandler(state *State, command Command) error {
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

func (cmds *commands) run(state *State, command Command) error {
	handler, ok := cmds.handlers[command.name]

	if !ok {
		return fmt.Errorf("error: command %s not found", command.name)
	}

	return handler(state, command)
}
