package test_pkg

import (
	"io/ioutil"
	"os"
)

func HijackStdout(f func()) (<-chan []byte, <-chan error) {
	errch := make(chan error)
	outch := make(chan []byte)

	oldstdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		go func() { errch <- err }()
		return outch, errch
	}

	os.Stdout = w

	f()

	os.Stdout = oldstdout

	go func() {
		hijackedStdout, err := ioutil.ReadAll(r)
		defer r.Close()
		if err != nil {
			errch <- err
			return
		}

		outch <- hijackedStdout
	}()

	w.Close()

	return outch, errch
}
