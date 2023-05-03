package ssp

func (c *Connection) Sync() (*Response, error) {
	c.seqFlag = true
	p, err := c.execute(Command{Code: Sync})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}


func (c *Connection) Reset() (*Response, error) {
	p, err := c.execute(Command{Code: Reset})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}
