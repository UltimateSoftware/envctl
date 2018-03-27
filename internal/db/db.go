package db

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

var (
	// StatusOff is an Environment's status when it hasn't been created yet.
	StatusOff = 0
	// StatusReady is an Environment's status when it's ready for use.
	StatusReady = 1
	// StatusError is an Environment's status when something is wrong with it.
	StatusError = 2
)

// Store is anything that can store an Environment.
type Store interface {
	Create(e Environment) error
	Read() (Environment, error)
	Delete() error
}

// Environment is just a container with its image under the hood. The container
// is really what runs it. To store it, all that needs to be tracked is the
// container and the image.
type Environment struct {
	Status    int    `json:"status"`
	Image     string `json:"image"`
	Container string `json:"container"`
}

// JSONStore implements a Store as a JSON file.
type JSONStore struct {
	File     *os.File
	basepath string
}

func NewJSONStore(basepath string) (js *JSONStore, err error) {
	js = &JSONStore{
		basepath: basepath,
	}

	if _, err = os.Stat(basepath); err != nil {
		err = os.Mkdir(basepath, os.ModePerm|os.ModeDir)
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(basepath + "envdata.json"); err != nil {
		js.File, err = os.OpenFile(
			basepath+"envdata.json",
			os.O_CREATE|os.O_RDWR,
			0666,
		)
		if err != nil {
			return
		}
	} else {
		js.File, err = os.OpenFile(
			basepath+"envdata.json",
			os.O_RDWR,
			0666,
		)
		if err != nil {
			return
		}
	}

	return
}

// Create writes an Environment to the file referenced by `js`.
func (js *JSONStore) Create(e Environment) error {
	buf, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = io.Copy(js.File, bytes.NewBuffer(buf))
	return err
}

// Read creates an Environment by reading the file referenced by `js` and
// returns it or an error if something went wrong. If the JSON Unmarshal returns
// an error, no error is returned. It's treated as an empty environment. This
// is because the subsequent call to Create will overwrite what's there when it
// writes the new Environment.
func (js *JSONStore) Read() (Environment, error) {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, js.File)
	if err != nil {
		return Environment{}, err
	}

	var e Environment
	json.Unmarshal(buf.Bytes(), &e)

	return e, err
}

// Delete writes an Environment with "status" set to StatusOff.
func (js *JSONStore) Delete() error {
	return os.RemoveAll(js.basepath)
}

// Initialized checks to see if an environment has been initialized. Initialized
// means any state that isn't "off", including "error" state.
func (e Environment) Initialized() bool {
	return e.Status != StatusOff
}
