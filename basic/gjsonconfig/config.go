package gjsonconfig

import (
	"encoding/json"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ginterface"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"io/ioutil"
	"runtime"
	"strings"
)

// DefaultConfig same dir same name.
func DefaultConfig() (string, error) {
	surffix := ".json"
	fn, err := gproc.SelfPath()
	if err != nil {
		return "", err
	}
	fnLower := strings.ToLower(fn)
	if (runtime.GOOS == "windows" && gstring.EndWith(fnLower, ".exe")) ||
		(runtime.GOOS == "darwin" && gstring.EndWith(fnLower, ".app")) ||
		(runtime.GOOS == "linux" && gstring.EndWith(fnLower, ".bin")) {
		fn = fn[0 : len(fn)-4]
	}
	return fn + surffix, nil
}

// DefaultUnmarshal read same dir same name .json config file and unmarshal to a struct
// v is a pointer to structure
func DefaultUnmarshal(v interface{}) error {
	filename, err := DefaultConfig()
	if err != nil {
		return err
	}
	return Unmarshal(filename, v)
}

// Unmarshal read same dir json conf file and unmarshal to a struct
// shortfn example: "*.json"
func Unmarshal(filename string, v interface{}) error {
	typeName, isPtr := ginterface.Parse(v)
	if !isPtr {
		return gerrors.Errorf("pointer needed for json.Unmarshal, but type %s is NOT a pointer", typeName)
	}

	pi, err := gfs.GetPathInfo(filename)
	if err != nil {
		return err
	}
	if pi.IsFolder || !pi.Exist {
		return gerrors.Errorf("config file '%s' is folder or not exist", filename)
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return gerrors.Errorf("config file '%s' content is empty", filename)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// DefaultMarshal set same dir same name .json config file and marshal from a struct
// v is a pointer to structure
func DefaultMarshal(v interface{}) error {
	filename, err := DefaultConfig()
	if err != nil {
		return err
	}
	return Marshal(filename, v)
}

func Marshal(filename string, v interface{}) error {
	pi, err := gfs.GetPathInfo(filename)
	if err != nil {
		return err
	}
	if pi.IsFolder || !pi.Exist {
		return gerrors.Errorf("config file '%s' is folder or not exist", filename)
	}

	str, err := gjson.MarshalString(v, true)
	if err != nil {
		return err
	}
	if err := gfs.StringToFile(str, filename); err != nil {
		return err
	}
	return nil
}
