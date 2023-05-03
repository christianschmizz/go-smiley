package ssp

func (c *Connection) Enable() (*Response, error) {
	p, err := c.execute(Command{Code: Enable})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}

func (c *Connection) Disable() (*Response, error) {
	p, err := c.execute(Command{Code: Disable})
	if err != nil {
		return nil, err
	}
	return p.Response(), nil
}
