package superset

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	BaseURL       url.URL
	AdminUser     string
	AdminPassword string
}

type Client struct {
	httpClient *http.Client
	config     Config
	token      string
	csrfToken  string
	logger     *zap.Logger
}

type GuestDescriptor struct {
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

func NewClient(config Config, logger *zap.Logger) *Client {
	cookieJar, _ := cookiejar.New(nil)
	client := &Client{
		httpClient: &http.Client{Timeout: time.Second * 10, Jar: cookieJar}, config: config, logger: logger}
	err := client.UpdateSession()
	if err != nil {
		panic(err) // FIXME: something very smelly
	}
	return client
}

func (c *Client) Authenticate() error {
	bodyRequest := map[string]any{
		"username": c.config.AdminUser,
		"password": c.config.AdminPassword,
		"refresh":  false,
		"provider": "db",
	}

	b, err := json.Marshal(bodyRequest)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		c.config.BaseURL.JoinPath("api", "v1", "security", "login").String(),
		bytes.NewBuffer(b),
	)

	req.Header.Set("Content-Type", "application/json")

	do, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		resp, _ := io.ReadAll(do.Body)
		return fmt.Errorf("authenticate: invalid status code: %d with body %s", do.StatusCode, string(resp))
	}
	var bodyResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	err = json.NewDecoder(do.Body).Decode(&bodyResponse)
	if err != nil {
		return err
	}
	c.token = bodyResponse.AccessToken
	return nil
}

func (c *Client) UpdateCSRFToken() error {
	req, _ := http.NewRequest(
		http.MethodGet,
		c.config.BaseURL.JoinPath("api", "v1", "security", "csrf_token").String(),
		nil,
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	do, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		resp, _ := io.ReadAll(do.Body)
		return fmt.Errorf("get csrf token: invalid status code: %d with body %s", do.StatusCode, string(resp))
	}
	var bodyResponse struct {
		Result string `json:"result"`
	}

	err = json.NewDecoder(do.Body).Decode(&bodyResponse)
	if err != nil {
		return err
	}

	c.csrfToken = bodyResponse.Result
	return nil
}

func (c *Client) UpdateSession() error {
	err := c.Authenticate()
	if err != nil {
		return err
	}
	return c.UpdateCSRFToken()
}

func (c *Client) GetGuestToken(dashboardId string, guest GuestDescriptor) (string, error) {
	token, err := c.doGetGuestToken(dashboardId, guest)
	if err == nil {
		return token, nil
	}

	if !errors.Is(err, ErrUnauthorized) {
		return "", err
	}

	if err := c.UpdateSession(); err != nil {
		return "", err
	}

	return c.doGetGuestToken(dashboardId, guest)
}

func (c *Client) doGetGuestToken(dashboardId string, guest GuestDescriptor) (string, error) {
	bodyRequest := map[string]any{
		"user": map[string]any{
			"username":   guest.Username,
			"first_name": guest.Firstname,
			"last_name":  guest.Lastname,
		},
		"rls": []map[string]string{},
		"resources": []map[string]string{
			{
				"type": "dashboard",
				"id":   dashboardId,
			},
		},
	}

	b, err := json.Marshal(bodyRequest)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		c.config.BaseURL.JoinPath("api", "v1", "security", "guest_token").String(),
		bytes.NewBuffer(b),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-CSRFToken", c.csrfToken)

	do, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

	if do.StatusCode != http.StatusOK {
		resp, _ := io.ReadAll(do.Body)
		return "", fmt.Errorf("get guest token: invalid status code: %d with body %s", do.StatusCode, string(resp))
	}

	var bodyResponse struct {
		Token string `json:"token"`
	}

	err = json.NewDecoder(do.Body).Decode(&bodyResponse)
	if err != nil {
		return "", err
	}

	return bodyResponse.Token, nil
}
