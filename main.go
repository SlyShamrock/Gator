package main
import (
	"github.com/SlyShamrock/Gator/internal/config"
	"fmt"
	"os"
	"errors"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %s", err)
		os.Exit(1)
	}

	newState := state{
		cfg: cfg,
	}

	cmds := commands{
		handlers : make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	
	args := os.Args
	if len(args) < 2 {
		err := errors.New("user must provide a commmand")
		fmt.Println(err)
		os.Exit(1)
	}
	newCmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&newState, newCmd)
	if err != nil {
		fmt.Printf("failed to run command: %s", err)
		os.Exit(1)
	}	
}