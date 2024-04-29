package gpio

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/warthog618/go-gpiocdev"
)

var Chip *gpiocdev.Chip

func Initialize() {
	c, err := gpiocdev.NewChip("gpiochip0", gpiocdev.WithConsumer(configs.Conf.AppName))
	if err != nil {
		log.Println("Error opening GPIO chip:", err)
	} else {
		Chip = c
		defer func() {
			_ = c.Close()
		}()
	}
}
func CheckChip() bool {
	return Chip != nil
}

func SetBool(pin int, direction string, value int) (int, error) {
	switch direction {
	case "in":
		l, err := Chip.RequestLine(pin, gpiocdev.AsInput)
		if err != nil {
			return 0, err
		}
		defer func() {
			_ = l.Close()
		}()
		err = db.SetPin(pin, value)
		if err != nil {
			return 0, err
		}

		val, err := l.Value()
		if err != nil {
			return 0, err
		}
		return val, nil
	case "out":
		l, err := Chip.RequestLine(pin, gpiocdev.AsOutput(value))
		if err != nil {
			return 0, err
		}
		defer func() {
			_ = l.Close()
		}()
		err = db.SetPin(pin, value)
		if err != nil {
			return 0, err
		}
		return value, nil
	default:
		return 0, errors.New("invalid direction provided")
	}
}

func GetAll() (map[int]int, error) {
	if !CheckChip() {
		return nil, errors.New("GPIO chip not initialized")
	}
	list, err := db.Get("gpio_list")
	if err != nil {
		return nil, err
	}
	var gpioList map[int]int
	_ = json.Unmarshal([]byte(list), &gpioList)
	return gpioList, nil
}
