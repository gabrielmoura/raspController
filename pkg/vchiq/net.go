package vchiq

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type INetStatistic struct {
	Interface string `json:"interface"`
	Mac       string `json:"mac"`
	RxBytes   int    `json:"rx_bytes"`
	TxBytes   int    `json:"tx_bytes"`
}

func GetNetStatistic() ([]INetStatistic, error) {
	files, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, errors.New("permission Denied")
	}

	var nets []INetStatistic
	for _, file := range files {
		if file.Type() == os.ModeSymlink {
			rx_bytes, err := os.ReadFile("/sys/class/net/" + file.Name() + "/statistics/rx_bytes")
			if err != nil {
				return nil, errors.New("permission Denied")
			}
			tx_bytes, err := os.ReadFile("/sys/class/net/" + file.Name() + "/statistics/tx_bytes")
			if err != nil {
				return nil, errors.New("permission Denied")
			}
			rx, err := strconv.Atoi(strings.TrimSpace(string(rx_bytes)))
			if err != nil {
				return nil, errors.New("error converting rx_bytes to int")
			}
			tx, err := strconv.Atoi(strings.TrimSpace(string(tx_bytes)))
			if err != nil {
				return nil, errors.New("error converting tx_bytes to int")
			}
			mac, err := os.ReadFile("/sys/class/net/" + file.Name() + "/address")
			if err != nil {
				return nil, errors.New("permission Denied")
			}
			nets = append(nets, INetStatistic{
				Interface: file.Name(),
				RxBytes:   rx,
				TxBytes:   tx,
				Mac:       strings.TrimSpace(string(mac)),
			})
		}
	}
	return nets, nil
}
