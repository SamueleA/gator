package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/samuelea/gator/internal/config"
)

type State struct {
	Config *config.Config
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
	},
}

func main() {
	gatorConfig, err := config.Read()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	state := State{
		Config: gatorConfig,
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

	err := state.Config.SetUser(username)

	if err != nil {
		return err
	}

	fmt.Printf("username %s logged in!\n", username)

	return nil
}

func (cmds *commands) run(state *State, command Command) error {
	handler, ok := cmds.handlers[command.name]

	if !ok {
		return fmt.Errorf("error: command %s not found", command.name)
	}

	return handler(state, command)
}

func (cmds *commands) register(name string, f func(*State, Command) error) {
	cmds.handlers[name] = f
}