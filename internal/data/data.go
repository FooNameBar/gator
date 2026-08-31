package data

import (
	c "github.com/FooNameBar/gator/internal/config"
	"github.com/FooNameBar/gator/internal/database"
)

type State struct {
	Config *c.Config
	DB     *database.Queries
}
