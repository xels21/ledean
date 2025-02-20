//go:build tinygo
// +build tinygo

package dbdriver

import "errors"

type DbDriver struct {
}

func NewDbDriver(dir string) (DbDriver, error) {
	self := DbDriver{}
	return self, nil
}

func (self *DbDriver) Read(collection, resource string, v interface{}) error {
	return errors.New("DbDriver stubbed")
}

func (self *DbDriver) Write(collection, resource string, v interface{}) error {
	return errors.New("DbDriver stubbed")
}
