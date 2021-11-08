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

func (c *Client) getConfigFilename(prefix, key string) string {
	return filepath.Join(c.cfgDir, prefix+"."+key)
}

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

func (c *Client) Store(prefix, key string, v interface{}) error {
	if v == nil || reflect.TypeOf(v) == nil {
		return gerrors.New("can't marshal nil value")
	}

	// Marshal and encrypt from structure.
	str := ""
	err := error(nil)
	if gstring.EndWith(key, ".json") {
		str, err = gjson.MarshalString(v, true)
		if err != nil {
			return err
		}
		err = gjson.Iterate(&str, true, c.encryptFn) // Encrypt JSON string if necessary.
	} else if gstring.EndWith(key, ".csv") {
		str, err = gocsv.MarshalString(v)
	} else {
		return gerrors.New("unsupported suffix of key %s", key)
	}
	if err != nil {
		return err
	}

	// Store into cache and storage.
	c.configs[strings.Join([]string{prefix, key}, ".")] = str
	return gfs.StringToFile(str, filepath.Join(c.cfgDir, prefix+"."+key))
}

func (c *Client) Load(prefix, key string, v interface{}, allowEmpty bool) error {
	// Load config string.
	cfgVal, ok := c.configs[strings.Join([]string{prefix, key}, ".")]
	if (!ok || cfgVal == "") && !allowEmpty {
		return gerrors.New("can't find config with prefix %s key %s", prefix, key)
	}

	// Decrypt and unmarshal into output structure.
	if ok && cfgVal != "" {
		if gstring.EndWith(strings.ToLower(key), ".json") {
			cfgValCopy := cfgVal
			if err := gjson.Iterate(&cfgValCopy, true, c.decryptFn); err != nil { // Decrypt JSON string if necessary.
				return err
			}
			return json.Unmarshal([]byte(cfgValCopy), v)
		} else if gstring.EndWith(strings.ToLower(key), ".csv") {
			return gocsv.UnmarshalBytes([]byte(cfgVal), v)
		} else {
			return gerrors.New("unsupported suffix of key %s", key)
		}
	}

	return nil
}

// Load config from config file first, then encrypt and store config to disk.
func (c *Client) LoadAndStore(prefix, key string, v interface{}, allowEmpty bool) error {
	if err := c.Load(prefix, key, v, allowEmpty); err != nil {
		return err
	}
	if err := c.Store(prefix, key, v); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetConfigDir() string {
	return c.cfgDir
}
