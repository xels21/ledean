//go:build tinygo && poi
// +build tinygo,poi

package ledean

import (
	"ledean/log"
	"ledean/mode"
	"time"

	"tinygo.org/x/espradio"
)

const (
	poiRadioMagic0       = 0x4c
	poiRadioMagic1       = 0x44
	poiRadioVersion      = 1
	poiRadioHello        = 1
	poiRadioMode         = 2
	poiRadioMaster       = 3
	poiRadioHeaderLength = 9
	poiRadioChunkLength  = espradio.ESPNowMaxDataLength - poiRadioHeaderLength
	poiRadioElectionTime = 1200 * time.Millisecond
)

type poiRadioPeer struct {
	addr espradio.ESPNowAddr
	seen bool
}

func startPoiRadio(modeController *mode.ModeController) {
	go runPoiRadio(modeController)
}

func runPoiRadio(modeController *mode.ModeController) {
	if err := espradio.Enable(espradio.Config{}); err != nil {
		log.Error("Could not enable ESP-NOW: ", err)
		return
	}
	if err := espradio.Start(); err != nil {
		log.Error("Could not start ESP-NOW: ", err)
		return
	}

	now, err := espradio.NewESPNow(espradio.ESPNowConfig{})
	if err != nil {
		log.Error("Could not initialize ESP-NOW: ", err)
		return
	}
	broadcast, err := now.Broadcast()
	if err != nil {
		log.Error("Could not initialize ESP-NOW broadcast: ", err)
		_ = now.Close()
		return
	}
	defer now.Close()
	defer broadcast.Close()

	localAddress := broadcast.LocalESPNowAddr()
	peers := make(map[espradio.ESPNowAddr]poiRadioPeer)
	deadline := time.Now().Add(poiRadioElectionTime)
	masterSeen := false
	for time.Now().Before(deadline) {
		_ = sendPoiRadioPacket(broadcast, poiRadioHello, 0, 0, 1, []byte(localAddress[:]))
		_ = broadcast.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		if packet, address, readErr := broadcast.ReadPacket(); readErr == nil {
			if sender, ok := address.(espradio.ESPNowAddr); ok {
				if packetType, _, _, _, payload, ok := decodePoiRadioPacket(packet); ok && packetType == poiRadioMaster {
					masterSeen = true
				} else if packetType == poiRadioHello && len(payload) == espradio.ESPNowAddressLength {
					var peerAddress espradio.ESPNowAddr
					copy(peerAddress[:], payload)
					peers[sender] = poiRadioPeer{addr: peerAddress, seen: true}
				}
			}
		}
	}

	isMaster := !masterSeen
	for _, peer := range peers {
		if poiRadioAddressLess(peer.addr, localAddress) {
			isMaster = false
			break
		}
	}
	if isMaster {
		log.Info("ESP-NOW role: master")
	} else {
		log.Info("ESP-NOW role: slave")
	}

	go receivePoiRadioModes(broadcast, modeController, isMaster)
	if isMaster {
		for {
			_ = sendPoiRadioPacket(broadcast, poiRadioMaster, 0, 0, 1, []byte(localAddress[:]))
			sendCurrentPoiMode(broadcast, modeController)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func poiRadioAddressLess(left, right espradio.ESPNowAddr) bool {
	for i := range left {
		if left[i] != right[i] {
			return left[i] < right[i]
		}
	}
	return false
}

func receivePoiRadioModes(broadcast *espradio.Peer, modeController *mode.ModeController, isMaster bool) {
	var message []byte
	for {
		if err := broadcast.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			return
		}
		packet, _, err := broadcast.ReadPacket()
		if err != nil {
			continue
		}
		packetType, sequence, index, total, payload, ok := decodePoiRadioPacket(packet)
		if !ok || packetType != poiRadioMode || isMaster {
			continue
		}
		if index == 0 {
			message = make([]byte, 0, int(total)*poiRadioChunkLength)
		}
		if sequence == 0 && index == 0 {
			message = message[:0]
		}
		message = append(message, payload...)
		if index+1 == total {
			if err := modeController.ApplyModeBinary(message); err != nil {
				log.Error("Could not apply synchronized mode: ", err)
			}
		}
	}
}

func sendCurrentPoiMode(broadcast *espradio.Peer, modeController *mode.ModeController) {
	data, err := modeController.GetCurrentModeBinary()
	if err != nil {
		log.Error("Could not serialize current mode: ", err)
		return
	}
	sequence := uint8(time.Now().UnixNano())
	for index, offset := 0, 0; offset < len(data); index++ {
		end := offset + poiRadioChunkLength
		if end > len(data) {
			end = len(data)
		}
		total := uint8((len(data) + poiRadioChunkLength - 1) / poiRadioChunkLength)
		if err := sendPoiRadioPacket(broadcast, poiRadioMode, sequence, uint8(index), total, data[offset:end]); err != nil {
			log.Error("Could not send synchronized mode: ", err)
			return
		}
		offset = end
	}
}

func sendPoiRadioPacket(peer *espradio.Peer, packetType, sequence, index, total uint8, payload []byte) error {
	packet := make([]byte, poiRadioHeaderLength+len(payload))
	packet[0] = poiRadioMagic0
	packet[1] = poiRadioMagic1
	packet[2] = poiRadioVersion
	packet[3] = packetType
	packet[4] = sequence
	packet[5] = index
	packet[6] = total
	packet[7] = uint8(len(payload) >> 8)
	packet[8] = uint8(len(payload))
	copy(packet[poiRadioHeaderLength:], payload)
	_, err := peer.WriteTo(packet, nil)
	return err
}

func decodePoiRadioPacket(packet []byte) (byte, byte, byte, byte, []byte, bool) {
	if len(packet) < poiRadioHeaderLength || packet[0] != poiRadioMagic0 || packet[1] != poiRadioMagic1 || packet[2] != poiRadioVersion {
		return 0, 0, 0, 0, nil, false
	}
	payloadLength := int(packet[7])<<8 | int(packet[8])
	if payloadLength != len(packet)-poiRadioHeaderLength || packet[6] == 0 {
		return 0, 0, 0, 0, nil, false
	}
	return packet[3], packet[4], packet[5], packet[6], packet[poiRadioHeaderLength:], true
}
