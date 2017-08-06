package newtransfer

type TransferStrategyInterface interface {
	Transfer(s, d string) error
	CleanUp(transferredDirs []string) error
}