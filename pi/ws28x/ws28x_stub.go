// +build !linux

package ws28x

import log "github.com/sirupsen/logrus"

func (self *PiWs28xConnector) Connect() error {
	go self.listen()
	return nil
}

func (self *PiWs28xConnector) listen() {
	for {
		log.Trace(<-self.cWriteBuffer)
	}
}
