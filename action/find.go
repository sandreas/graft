package action

import (
	"fmt"

	"github.com/urfave/cli"
)

type FindAction struct {
	AbstractAction
}

func (act *FindAction) Execute(c *cli.Context) error {
	act.PrepareExecution(c, 1)
	fmt.Println("find action")
	// fmt.Println(act.MinAge)
	return nil
}
