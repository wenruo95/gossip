package control

type ticker struct {
}

func newTicker() *ticker {
	t := new(ticker)
	return t
}

func (t *ticker) Serve() error {
	return nil
}
