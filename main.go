package main
import _ "github.com/lib/pq"
import (
	"github.com/SlyShamrock/Gator/internal/config"
	"github.com/SlyShamrock/Gator/internal/database"
	"fmt"
	"os"
	"errors"
	"database/sql"
	
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %s", err)
		os.Exit(1)
	}
	
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("failed to open database: %s", err)
	}

	dbQueries := database.New(db)
	
	newState := state{
		db: dbQueries,
		cfg: cfg,
	}

	cmds := commands{
		handlers : make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
		
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