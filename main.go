package main
import (
	"github.com/SlyShamrock/Gator/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %s", err)
	}

	err = cfg.SetUser("slyshamrock")
	if err != nil {
		fmt.Printf("failed to set new username: %s", err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %s", err)
	}

	fmt.Printf("%s\n%s\n", cfg.DBURL, cfg.CurrentUserName)
}