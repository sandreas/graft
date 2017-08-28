package action

import "github.com/urfave/cli"

type ConnectionInterface interface {
	Connect(url string)
	Disconnect() error
}


type CliActionInterface interface {
	Execute(c *cli.Context) error
}