//go:build !tinygo

package dbdriver

import "github.com/sdomino/scribble"

type DbDriver struct {
	driver *scribble.Driver
}

func NewDbDriver(dir string) (*DbDriver, error) {
	self := DbDriver{}

	driver, err := scribble.New(dir, nil)
	self.driver = driver
	return &self, err
}

func (self *DbDriver) Read(collection, resource string, v interface{}) error {
	return self.driver.Read(collection, resource, v)
}

func (self *DbDriver) Write(collection, resource string, v interface{}) error {
	return self.driver.Write(collection, resource, v)
}
