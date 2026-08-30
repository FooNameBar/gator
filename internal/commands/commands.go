package commands

import (
	"fmt"

	"github.com/FooNameBar/gator/internal/data"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*data.State, Command) error
}

func (c *Commands) Run(s *data.State, cmd Command) error {
	com, exists := c.List[cmd.Name]
	if !exists {
		return fmt.Errorf("The command %s doesn't exist\n", cmd.Name)
	}

	return com(s, cmd)
}

func (c *Commands) Register(name string, f func(*data.State, Command) error) {
	c.List[name] = f
}
