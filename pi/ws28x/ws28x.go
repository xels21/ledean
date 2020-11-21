package ws28x

type PiWs28xConnector struct {
	spiInfo      string
	cWriteBuffer chan []byte
}

func NewPiWs28xConnector(spiInfo string) *PiWs28xConnector {
	return &PiWs28xConnector{
		spiInfo:      spiInfo,
		cWriteBuffer: make(chan []byte, 32),
	}
}

func (self *PiWs28xConnector) Write(data []byte) {
	self.cWriteBuffer <- data
}
