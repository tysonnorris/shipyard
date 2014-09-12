package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/citadel/citadel"
	"github.com/shipyard/shipyard"
)

type (
	Manager struct {
		baseUrl string
		config  *ShipyardConfig
	}
)

func NewManager(cfg *ShipyardConfig) *Manager {
	m := &Manager{
		config: cfg,
	}
	return m
}

func (m *Manager) buildUrl(path string) string {
	return fmt.Sprintf("%s%s", m.config.Url, path)
}

func (m *Manager) doRequest(path string, method string, expectedStatus int, b []byte) (*http.Response, error) {
	url := m.buildUrl(path)
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if m.config.ServiceKey != "" {
		req.Header.Add("X-Service-Key", m.config.ServiceKey)
	} else {
		req.Header.Add("X-Access-Token", fmt.Sprintf("%s:%s", m.config.Username, m.config.Token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return resp, shipyard.ErrUnauthorized
	}

	if resp.StatusCode != expectedStatus {
		c, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return resp, errors.New(string(c))
	}
	return resp, nil
}

func (m *Manager) Containers() ([]*citadel.Container, error) {
	containers := []*citadel.Container{}
	resp, err := m.doRequest("/api/containers", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (m *Manager) Run(image *citadel.Image, count int, pull bool) ([]*citadel.Container, error) {
	b, err := json.Marshal(image)
	if err != nil {
		return nil, err
	}
	var containers []*citadel.Container
	resp, err := m.doRequest(fmt.Sprintf("/api/containers?count=%d&pull=%v", count, pull), "POST", 201, b)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (m *Manager) Destroy(container *citadel.Container) error {
	b, err := json.Marshal(container)
	if err != nil {
		return err
	}
	if _, err := m.doRequest(fmt.Sprintf("/api/containers/%s", container.ID), "DELETE", 204, b); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Engines() ([]*shipyard.Engine, error) {
	engines := []*shipyard.Engine{}
	resp, err := m.doRequest("/api/engines", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&engines); err != nil {
		return nil, err
	}
	return engines, nil
}

func (m *Manager) AddEngine(engine *shipyard.Engine) error {
	b, err := json.Marshal(engine)
	if err != nil {
		return err
	}
	if _, err := m.doRequest("/api/engines", "POST", 201, b); err != nil {
		return err
	}
	return nil
}

func (m *Manager) RemoveEngine(engine *shipyard.Engine) error {
	if _, err := m.doRequest(fmt.Sprintf("/api/engines/%s", engine.Engine.ID), "DELETE", 204, nil); err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetContainer(id string) (*citadel.Container, error) {
	var container *citadel.Container
	resp, err := m.doRequest(fmt.Sprintf("/api/containers/%s", id), "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
		return nil, err
	}
	return container, nil
}

func (m *Manager) GetEngine(id string) (*shipyard.Engine, error) {
	var engine *shipyard.Engine
	resp, err := m.doRequest(fmt.Sprintf("/api/engines/%s", id), "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&engine); err != nil {
		return nil, err
	}
	return engine, nil
}

func (m *Manager) Info() (*citadel.ClusterInfo, error) {
	var info *citadel.ClusterInfo
	resp, err := m.doRequest("/api/cluster/info", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return info, nil
}

func (m *Manager) Events() ([]*shipyard.Event, error) {
	events := []*shipyard.Event{}
	resp, err := m.doRequest("/api/events", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

func (m *Manager) Accounts() ([]*shipyard.Account, error) {
	accounts := []*shipyard.Account{}
	resp, err := m.doRequest("/api/accounts", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (m *Manager) Roles() ([]*shipyard.Role, error) {
	roles := []*shipyard.Role{}
	resp, err := m.doRequest("/api/roles", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (m *Manager) Role(name string) (*shipyard.Role, error) {
	role := &shipyard.Role{}
	resp, err := m.doRequest(fmt.Sprintf("/api/roles/%s", name), "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return role, nil
}

func (m *Manager) AddAccount(account *shipyard.Account) error {
	b, err := json.Marshal(account)
	if err != nil {
		return err
	}
	if _, err := m.doRequest("/api/accounts", "POST", 204, b); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DeleteAccount(account *shipyard.Account) error {
	b, err := json.Marshal(account)
	if err != nil {
		return err
	}
	if _, err := m.doRequest("/api/accounts", "DELETE", 204, b); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Login(username, password string) (*shipyard.AuthToken, error) {
	creds := map[string]string{}
	creds["username"] = username
	creds["password"] = password
	b, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	resp, err := m.doRequest("/auth/login", "POST", 200, b)
	if err != nil {
		return nil, err
	}
	var token *shipyard.AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Manager) ChangePassword(password string) error {
	creds := map[string]string{}
	creds["password"] = password
	b, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	if _, err := m.doRequest("/account/changepassword", "POST", 200, b); err != nil {
		return err
	}
	return nil
}

func (m *Manager) ServiceKeys() ([]*shipyard.ServiceKey, error) {
	keys := []*shipyard.ServiceKey{}
	resp, err := m.doRequest("/api/servicekeys", "GET", 200, nil)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (m *Manager) NewServiceKey(description string) (*shipyard.ServiceKey, error) {
	k := &shipyard.ServiceKey{
		Description: description,
	}
	b, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	resp, err := m.doRequest("/api/servicekeys", "POST", 200, b)
	if err != nil {
		return nil, err
	}
	var key *shipyard.ServiceKey
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, err
	}
	return key, nil
}

func (m *Manager) RemoveServiceKey(key *shipyard.ServiceKey) error {
	b, err := json.Marshal(key)
	if err != nil {
		return err
	}
	if _, err := m.doRequest("/api/servicekeys", "DELETE", 204, b); err != nil {
		return err
	}
	return nil
}
