package pangolin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Pangolin struct {
	token      string
	host       string
	org        string
	httpClient http.Client
}

type PangolinResource struct {
	ResourceID *int    `json:"resourceId"`
	Name       *string `json:"name"`
	Enabled    *bool   `json:"enabled"`
	FullDomain *string `json:"fullDomain"`
}

type PangolinRule struct {
	RuleID     *int    `json:"ruleId,omitempty"`
	ResourceID *int    `json:"resourceId,omitempty"`
	Enabled    *bool   `json:"enabled,omitempty"`
	Priority   *int    `json:"priority,omitempty"`
	Action     *string `json:"action,omitempty"`
	Match      *string `json:"match,omitempty"`
	Value      *string `json:"value,omitempty"`
}

type PangolinRules struct {
	Rules []PangolinRule `json:"rules"`
}

type PangolinResources struct {
	Resources []PangolinResource `json:"resources"`
}

type PangolinResponse[T any] struct {
	Data    T      `json:"data"`
	Success bool   `json:"success"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func InitPangolin(token, host, org string, httpClient http.Client) *Pangolin {
	return &Pangolin{token, host, org, httpClient}
}

func (p Pangolin) baseURL() string {
	return p.host + "/v1"
}

func (p Pangolin) request(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+p.token)
	req.Header.Add("Content-Type", "application/json")
	return p.httpClient.Do(req)
}

func (p Pangolin) GetResources() (*PangolinResources, error) {
	url := p.baseURL() + "/org/" + p.org + "/resources"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := p.request(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("error status code: " + res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp PangolinResponse[PangolinResources]
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (p Pangolin) CreateRule(rule PangolinRule, resourceID int) (*PangolinRule, error) {
	url := p.baseURL() + "/resource/" + strconv.Itoa(resourceID) + "/rule"

	reqBody, err := json.Marshal(rule)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	res, err := p.request(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp PangolinResponse[PangolinRule]
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error {
		return nil, errors.New("error: " + resp.Message)
	}

	return &resp.Data, nil
}

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}
