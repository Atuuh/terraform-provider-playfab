package playfab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	titleId   string
	secretKey string
	basePath  string
	data      map[string]any
}

func NewClient(titleId *string, secretKey *string) (*Client, error) {
	return &Client{
		titleId:   *titleId,
		secretKey: *secretKey,
		basePath:  fmt.Sprintf("https://%s.playfabapi.com", *titleId),
		data:      map[string]any{},
	}, nil
}

type GetEntityTokenResponse struct {
	Data struct {
		EntityToken string `json:"EntityToken"`
	} `json:"data"`
}

func (c *Client) SetData(key string, value any) {
	c.data[key] = value
}

func (c *Client) GetTitleEntityId() (string, error) {
	endpoint := fmt.Sprintf("%s/Authentication/GetEntityToken", c.basePath)
	// fmt.Printf("calling endpoint %s", endpoint)
	req, err := http.NewRequest("POST", endpoint, http.NoBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-SecretKey", c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	// fmt.Printf("making request %+v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d with body %s", resp.StatusCode, body)
	}
	if err != nil {
		return "", err
	}
	// fmt.Printf("body %s", body)

	var entityResponse GetEntityTokenResponse
	err = json.Unmarshal(body, &entityResponse)
	if err != nil {
		return "", err
	}

	return entityResponse.Data.EntityToken, nil
}

type Function struct {
	Name        string `json:"FunctionName"`
	Address     string `json:"FunctionAddress"`
	TriggerType string `json:"TriggerType"`
}

type ListFunctionsResult struct {
	Data struct {
		Functions []*Function `json:"Functions"`
	} `json:"data"`
}

func (c *Client) GetCloudScriptFunctions() ([]*Function, error) {
	endpoint := fmt.Sprintf("%s/CloudScript/ListFunctions", c.basePath)
	req, err := http.NewRequest("POST", endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	titleEntityKey, ok := c.data["title_entity_token"]
	if !ok {
		return nil, fmt.Errorf("no title entity key")
	}
	req.Header.Set("X-EntityToken", titleEntityKey.(string))
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("sending req %+v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d with body %s", resp.StatusCode, body)
	}
	var response ListFunctionsResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("got body %s", body)

	return response.Data.Functions, nil
}

type GetFunctionResult struct {
	Data *Function `json:"data"`
}

func (c *Client) GetCloudScriptFunction(name string) (*Function, error) {
	endpoint := fmt.Sprintf("%s/CloudScript/GetFunction", c.basePath)
	reqBody := fmt.Sprintf(`{ "FunctionName": "%s" }`, name)
	fmt.Printf("req body %s", reqBody)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(reqBody))
	if err != nil {
		return nil, err
	}
	titleEntityKey, ok := c.data["title_entity_token"]
	if !ok {
		return nil, fmt.Errorf("no title entity key")
	}
	req.Header.Set("X-EntityToken", titleEntityKey.(string))
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("sending req %+v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d with body %s", resp.StatusCode, body)
	}
	var response GetFunctionResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	fmt.Printf("got body %s", body)
	response.Data.Name = name
	return response.Data, nil
}

func (c *Client) CreateCloudScriptFunction(function *Function) error {
	endpoint := fmt.Sprintf("%s/CloudScript/RegisterHttpFunction", c.basePath)
	reqBody := fmt.Sprintf(`{ "FunctionName": "%s", "FunctionUrl": "%s" }`, function.Name, function.Address)
	fmt.Printf("req body %s", reqBody)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(reqBody))
	if err != nil {
		return err
	}
	titleEntityKey, ok := c.data["title_entity_token"]
	if !ok {
		return fmt.Errorf("no title entity key")
	}
	req.Header.Set("X-EntityToken", titleEntityKey.(string))
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("sending req %+v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d with body %s", resp.StatusCode, body)
	}
	return nil
}
