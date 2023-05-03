package ssp

type VersionResponse struct {
	*Response
	FirmwareVersion string
	DeviceVersion   string
	ReleaseVersion  string
	BetaVersion     string
}

func NewVersionResponse(p *Packet) *VersionResponse {
	return &VersionResponse{
		Response:        p.Response(),
		FirmwareVersion: string(p.args),
		DeviceVersion:   string(p.args[0:6]),
		ReleaseVersion:  string(p.args[9:13]),
		BetaVersion:     string(p.args[13:16]),
	}
}

func (c *Connection) GetFirmwareVersion() (*VersionResponse, error) {
	p, err := c.execute(Command{Code: GetFirmwareVersion})
	if err != nil {
		return nil, err
	}
	return NewVersionResponse(p), nil
}

func (c *Connection) GetDatasetVersion() (string, error) {
	p, err := c.execute(Command{Code: GetDatasetVersion})
	if err != nil {
		return "", nil
	}
	return string(p.args), nil
}
