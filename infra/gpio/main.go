package gpio

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/warthog618/go-gpiocdev"
)

var Chip *gpiocdev.Chip
var chipContext context.Context

func Initialize(ctx context.Context) error {
	c, err := gpiocdev.NewChip("gpiochip0", gpiocdev.WithConsumer(configs.Conf.AppName))
	if err != nil {
		log.Println("Error opening GPIO chip:", err.Error())
	} else {
		Chip = c
		chipContext = ctx
		if ctx.Err() != nil {
			log.Println("Closing GPIO chip")
			_ = c.Close()
			return ctx.Err()
		}
	}
	return nil
}
func CheckChip() bool {
	return Chip != nil
}
func setOutput(pin int, value int) error {
	l, err := Chip.RequestLine(pin, gpiocdev.AsOutput(value))
	if err != nil {
		log.Println("GPIO: Error setting output:", err.Error())
		return err
	}
	log.Printf("GPIO: Pin %d set to %d", pin, value)
	stat, _ := l.Info()

	log.Printf("GPIO: Pin %d value: %d status: %+v", pin, value, stat)

	if chipContext.Err() != nil {
		log.Printf("GPIO: Closing pin %d", pin)
		_ = l.Close()
	}
	// defer func() {
	// 	_ = l.Close()
	// }()
	err = db.SetPin(pin, value)
	if err != nil {
		return fmt.Errorf("GPIO: Error setting pin value: %s", err.Error())
	}
	return nil
}
func setInput(pin int) error {
	l, err := Chip.RequestLine(pin, gpiocdev.AsInput)
	if err != nil {
		return err
	}
	if chipContext.Err() != nil {
		log.Printf("GPIO: Closing pin %d", pin)
		_ = l.Close()
	}
	err = db.SetPin(pin, 0)
	if err != nil {
		return err
	}

	val, err := l.Value()
	if err != nil {
		return err
	}
	log.Printf("Pin %d value: %d", pin, val)
	return nil
}

func SetBool(pin int, direction string, value int) error {
	if direction != "in" && direction != "out" {
		return errors.New("invalid direction")
	} else if direction == "out" {
		return setOutput(pin, value)
	} else {
		return setInput(pin)
	}

}

func GetAll() (db.Map, error) {
	if !CheckChip() {
		return nil, errors.New("GPIO chip not initialized")
	}

	gpios, err := db.GetAllPin()
	if err != nil {
		return nil, fmt.Errorf("error getting GPIO list: %s", err.Error())
	}
	return gpios, nil
}

func GetGpioAll() (map[int]string, error) {
	if !CheckChip() {
		return nil, errors.New("GPIO chip not initialized")
	}

	// Obtém o número total de linhas no chip GPIO
	numLines := Chip.Lines()

	usedPins := make(map[int]string)

	// Itera sobre todas as linhas possíveis
	for offset := 0; offset < numLines; offset++ {
		// Tenta obter informações sobre a linha
		info, err := Chip.LineInfo(offset)
		if err != nil {
			log.Printf("Erro ao obter informações da linha %d: %v\n", offset, err)
			continue
		}

		// Verifica se a linha está em uso
		if info.Consumer != "" {
			usedPins[offset] = info.Consumer
		}
	}

	return usedPins, nil

}
