package gconfig

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/crypto/gfes"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"github.com/gocarina/gocsv"
	"path/filepath"
	"reflect"
	"strings"
)

type (
	Client struct {
		configs  map[string]string
		password string
		salt     string
		cfgDir   string
	}
)

// NewClient creates new config client.
func NewClient(customConfigDir string) (*Client, error) {
	c := &Client{configs: map[string]string{}}

	cfgDir := ""
	err := error(nil)
	if customConfigDir != "" {
		cfgDir = customConfigDir
	} else {
		cfgDir, err = gsysinfo.GetHomeDir()
		if err != nil {
			return nil, err
		}
		cfgDir = filepath.Join(cfgDir, "config")
	}
	c.cfgDir = cfgDir

	_, files, err := gfs.ListDir(cfgDir)
	if err != nil {
		return nil, err
	}
	for _, filename := range files {
		cfgKey := gfs.PathBase(filename)
		if cfgKey == "" {
			return nil, gerrors.New("invalid config key %s", cfgKey)
		}
		cfgVal, err := gfs.FileToString(filename)
		if err != nil {
			return nil, err
		}
		c.configs[cfgKey] = cfgVal
	}

	return c, nil
}

// getConfigFilePath returns config file full path.
func (c *Client) getConfigFilePath(configFileName string) string {
	return filepath.Join(c.cfgDir, configFileName)
}

// SetPassword sets password for config.
func (c *Client) SetPassword(password, salt string) {
	c.password = password
	c.salt = salt
}

func (c *Client) cryptFn(method, key string, val interface{}) (newVal interface{}, modified bool, err error) {
	if !gstring.EndWith(key, "EncryptMe") { // Don't encrypt/decrypt key doesn't end with "EncryptMe".
		return nil, false, nil
	}

	if reflect.TypeOf(val) != nil && reflect.TypeOf(val).Kind() == reflect.String {
		if val.(string) == "" { // Don't encrypt/decrypt empty string.
			return nil, false, nil
		}
		if method == "encrypt" {
			cipher, err := gfes.TriMartolodEncrypt(val.(string), c.password, c.salt)
			if err != nil {
				return nil, false, gerrors.Wrap(err, fmt.Sprintf("%s key %s", method, key))
			} else {
				return cipher, true, nil
			}
		} else if method == "decrypt" {
			plain, err := gfes.TriMartolodDecrypt(val.(string), c.password, c.salt)
			if err == nil {
				return plain, true, nil
			}
		} else {
			return nil, false, gerrors.New("unsupported cryptFn method %s", method)
		}
	}
	return nil, false, nil
}

func (c *Client) encryptFn(key string, val interface{}) (newVal interface{}, modified bool, err error) {
	return c.cryptFn("encrypt", key, val)
}

func (c *Client) decryptFn(key string, val interface{}) (newVal interface{}, modified bool, err error) {
	return c.cryptFn("decrypt", key, val)
}

// Store writes config content `v` into file `configFileName`.
// configFileName is config file short name with suffix, for example `myapp.json`.
func (c *Client) Store(configFileName string, v interface{}) error {
	if v == nil || reflect.TypeOf(v) == nil {
		return gerrors.New("can't marshal nil value")
	}

	// Marshal and encrypt from structure.
	str := ""
	err := error(nil)
	if gstring.EndWith(configFileName, ".json") {
		str, err = gjson.MarshalString(v, true)
		if err != nil {
			return err
		}
		err = gjson.Iterate(&str, true, c.encryptFn) // Encrypt JSON string if necessary.
	} else if gstring.EndWith(configFileName, ".csv") {
		str, err = gocsv.MarshalString(v)
	} else {
		return gerrors.New("unsupported suffix of config %s", configFileName)
	}
	if err != nil {
		return err
	}

	// Store into cache and storage.
	c.configs[configFileName] = str
	return gfs.StringToFile(str, filepath.Join(c.cfgDir, configFileName))
}

// Load loads config and Unmarshal it into `v`.
// configFileName is config file short name with suffix, for example `myapp.json`.
func (c *Client) Load(configFileName string, v interface{}, allowEmpty bool) error {
	// Load config string.
	cfgVal, ok := c.configs[configFileName]
	if (!ok || cfgVal == "") && !allowEmpty {
		return gerrors.New("can't find config with %s", configFileName)
	}

	// Decrypt and unmarshal into output structure.
	if ok && cfgVal != "" {
		if gstring.EndWith(strings.ToLower(configFileName), ".json") {
			cfgValCopy := cfgVal
			if err := gjson.Iterate(&cfgValCopy, true, c.decryptFn); err != nil { // Decrypt JSON string if necessary.
				return err
			}
			return json.Unmarshal([]byte(cfgValCopy), v)
		} else if gstring.EndWith(strings.ToLower(configFileName), ".csv") {
			return gocsv.UnmarshalBytes([]byte(cfgVal), v)
		} else {
			return gerrors.New("unsupported suffix of config %s", configFileName)
		}
	}

	return nil
}

// LoadAndStore loads config from config file first, then encrypt and store config to disk.
// configFileName is config file short name with suffix, for example `myapp.json`.
func (c *Client) LoadAndStore(configFileName string, v interface{}, allowEmpty bool) error {
	if err := c.Load(configFileName, v, allowEmpty); err != nil {
		return err
	}
	if err := c.Store(configFileName, v); err != nil {
		return err
	}
	return nil
}

// ConfigFileExists returns true if config file exists.
// configFileName is config file short name with suffix, for example `myapp.json`.
func (c *Client) ConfigFileExists(configFileName string) bool {
	fileName := filepath.Join(c.cfgDir, configFileName)
	return gfs.FileExits(fileName)
}

func (c *Client) ConfigDir() string {
	return c.cfgDir
}
