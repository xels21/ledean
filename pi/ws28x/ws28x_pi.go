// +build linux

package ws28x

import (
	log "github.com/sirupsen/logrus"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/experimental/devices/nrzled"
)

func (self *PiWs28xConnector) Connect(ledCount int) error {

	portCloser, err := spireg.Open(self.spiInfo)
	// defer portCloser.Close()
	if err != nil {
		log.Panic(err)
	}
	opts := &nrzled.Opts{
		Channels:  3,
		NumPixels: ledCount,
		Freq:      2500 * physic.KiloHertz, //should be 800khz...

	}
	var dev *nrzled.Dev
	dev, err = nrzled.NewSPI(portCloser, opts)

	if err != nil {
		log.Fatal(err)
	}

	go self.listen(dev)

	return nil
}

func (self *PiWs28xConnector) listen(dev *nrzled.Dev) {
	var err error
	var pixels []byte

	for {
		pixels = <-self.cWriteBuffer

		log.Trace("write: ", pixels)

		_, err = (*dev).Write(pixels)

		if err != nil {
			log.Fatal(err)
		}

		// time.Sleep(50 * time.Microsecond)

	}
}
