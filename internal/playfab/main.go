package playfab

type Client struct {
}

func NewClient(titleId *string, secretKey *string) (*Client, error) {
	return &Client{}, nil
}

type Function struct {
	Name        string `json:"FunctionName"`
	Address     string `json:"FunctionAddress"`
	TriggerType string `json:"TriggerType"`
}

func (c *Client) GetCloudScriptFunctions() ([]*Function, error) {
	return []*Function{
		{
			Name:        "Hard Coded Fn",
			Address:     "http://localhost:12345/test",
			TriggerType: "HTTP",
		},
	}, nil
}

func (c *Client) GetCloudScriptFunction(name string) (*Function, error) {
	return &Function{
		Name:        "Hard Coded Fn",
		Address:     "http://localhost:12345/test",
		TriggerType: "HTTP",
	}, nil
}

func (c *Client) CreateCloudScriptFunction(function *Function) error {
	return nil
}
