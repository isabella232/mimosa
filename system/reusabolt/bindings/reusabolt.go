package bindings

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Reusabolt represents a deployed Reusabolt instance
type Reusabolt struct {
	URL string
}

// NewReusabolt creates a Reusabolt instance
func NewReusabolt(url string) *Reusabolt {
	return &Reusabolt{
		URL: url,
	}
}

// RunTask via this Reusabolt instance
func (r *Reusabolt) RunTask(te TaskExecution) (int, string, error) {
	data, err := json.Marshal(te)
	if err != nil {
		return 0, "", err
	}
	response, err := http.Post(r.URL+"/v1/task", "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, "", err
	}
	return response.StatusCode, string(data), nil
}

// TaskExecution contains all the context necessary to run Reusabolt
type TaskExecution struct {
	Targets   []string          `json:"targets,omitempty"`
	Name      string            `json:"name,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
	Inventory inventory         `json:"inventory,omitempty"`
}

type inventory struct {
	Nodes []inventoryNode `json:"nodes,omitempty"`
}

type inventoryNode struct {
	Name   string          `json:"name,omitempty"`
	Config inventoryConfig `json:"config,omitempty"`
}

type inventoryConfig struct {
	Transport string          `json:"transport"`
	WinRM     *inventoryWinRM `json:"winrm,omitempty"`
	SSH       *inventorySSH   `json:"ssh,omitempty"`
}

type inventoryWinRM struct {
	User              string `json:"user,omitempty"`
	Password          string `json:"password,omitempty"`
	SSL               bool   `json:"ssl"`                 // Don't omit empty
	AllowHTTPFallback bool   `json:"allow_http_fallback"` // Attempt http if https fails
	SSLVerify         bool   `json:"ssl-verify"`          // Don't omit empty
}

type inventorySSH struct {
	User         string               `json:"user,omitempty"`
	Password     string               `json:"password,omitempty"`
	PrivateKey   *inventoryPrivateKey `json:"private-key,omitempty"`
	HostKeyCheck bool                 `json:"host-key-check"` // Don't omit empty
}

type inventoryPrivateKey struct {
	KeyData string `json:"key-data,omitempty"`
}
