package ssp

func (c *Connection) DisplayOff() (*Response, error) {
	p, err := c.execute(Command{Code: DisplayOff})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}

func (c *Connection) DisplayOn() (*Response, error) {
	p, err := c.execute(Command{Code: DisplayOn})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}
