package newtransfer

type TransferStrategyInterface interface {
	Transfer(s, d string) error
}