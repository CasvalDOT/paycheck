package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"paycheck/helpers"
)

type config struct {
	SvCloudUsername string `json:"svcloud_username"`
	SvCloudPassword string `json:"svcloud_password"`
	SvCloudBaseURL  string `json:"svcloud_base_url"`
	SvCloudEndpoint string `json:"svcloud_endpoint"`
	BoxToken        string `json:"box_token"`
	BoxTargetID     string `json:"box_target_id"`
	repo            string
	pubKey          string
	path            string
}

// IConfig is the main interface
type IConfig interface {
	GetSvCloudUsername() string
	GetSvCloudPassword() string
	GetBoxToken() string
	GetBoxTargetID() string
	GetSvCloudBaseURL() string
	GetSvCloudEndpoint() string
	GetRepo() string
	GetPubKey() string
}

func (c *config) check() {
	stats, err := os.Stat(c.path)
	if err != nil {
		log.Fatalln(err)
	}

	mode := stats.Mode()
	if mode.Perm() != os.FileMode(0o600) {
		helpers.MessageWarning(fmt.Sprintf("%s permits set to %s instead of 600", c.path, mode.Perm()))
	}
}

func (c *config) checkRepository() {
	stats, err := os.Stat(c.repo)
	if err != nil {
		log.Fatalln(err)
	}

	mode := stats.Mode()
	if mode.Perm() != os.FileMode(0o700) {
		helpers.MessageWarning(fmt.Sprintf("%s permits set to %s instead of 700", c.repo, mode.Perm()))
	}
}

func (c *config) load() (*config, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	basePath := homePath + "/.config/paycheck"
	filePath := basePath + "/config.json"

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	c.repo = basePath + "/repository"
	c.pubKey = basePath + "/paycheck.asc"
	c.path = filePath

	return c, nil
}

func (c *config) GetSvCloudUsername() string {
	return c.SvCloudUsername
}

func (c *config) GetSvCloudPassword() string {
	return c.SvCloudPassword
}

func (c *config) GetBoxToken() string {
	return c.BoxToken
}

func (c *config) GetBoxTargetID() string {
	return c.BoxTargetID
}

func (c *config) GetSvCloudBaseURL() string {
	return c.SvCloudBaseURL
}

func (c *config) GetSvCloudEndpoint() string {
	return c.SvCloudEndpoint
}

func (c *config) GetRepo() string {
	return c.repo
}

func (c *config) GetPubKey() string {
	return c.pubKey
}

// New configuration
func New() (IConfig, error) {
	c := config{}
	_, err := c.load()
	if err != nil {
		return &c, err
	}

	c.check()
	c.checkRepository()

	return &c, nil
}
