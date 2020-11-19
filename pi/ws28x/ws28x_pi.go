// +build linux

package ws28x

import (
	"time"

	log "github.com/sirupsen/logrus"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

func (self *PiWs28xConnector) Connect() error {
	var portCloser spi.PortCloser
	var err error
	portCloser, err = spireg.Open(self.spiInfo)
	// defer portCloser.Close()
	if err != nil {
		log.Panic(err)
	}

	if p, ok := portCloser.(spi.Pins); ok {
		// log.Trace("  CLK : %s", p.CLK())
		log.Trace("  MOSI: ", p.MOSI())
		// log.Trace("  MISO: %s", p.MISO())
		// log.Trace("  CS  : %s", p.CS())
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	// conn, err := portCloser.Connect(3*400*physic.Hertz, spi.Mode0, 8)
	// conn, err := portCloser.Connect(3*physic.MegaHertz, spi.Mode0, 8)

	// CAUTION:
	// It was observed on RPi3 hardware to have a one clock delay between each packet.
	conn, err := portCloser.Connect(3*800*physic.KiloHertz, spi.Mode0|spi.NoCS, 8)
	// conn, err := portCloser.Connect(800*physic.KiloHertz, spi.Mode0, 8)
	if err != nil {
		log.Fatal(err)
	}

	go self.listen(conn)

	return nil
}

func (self *PiWs28xConnector) listen(conn conn.Conn) {
	var err error
	var write, read []byte
	// r := make([]byte, 0)
	for {
		write = <-self.cWriteBuffer
		// write = []byte{0, 1, 2, 3, 4, 5}

		// read := make([]byte, len(write))
		log.Trace("write: ", write)

		err = conn.Tx(write, read)
		if err != nil {
			log.Panic("something went wrong: ", err) // Todo reinit spi conn
		}
		// log.Trace("1")
		time.Sleep(200 * time.Microsecond) //RESET - 50us - just tested TODO remove when wiring is done
		// log.Trace("2")

		// err = conn.Tx([]byte{0, 0, 0}, read)
		// if err != nil {
		// 	log.Panic("something went wrong: ", err) // Todo reinit spi conn
		// }
	}
}
