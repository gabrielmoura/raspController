package gpio

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/gabrielmoura/raspController/internal/dto"
	"github.com/warthog618/go-gpiocdev"
)

var (
	Chip  *gpiocdev.Chip
	lines = make(map[int]*gpiocdev.Line)
	mu    sync.RWMutex // RWMutex permite múltiplas leituras concorrentes
	once  sync.Once    // Garante que a inicialização ocorra apenas uma vez
)

func initializeChip(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	c, err := gpiocdev.NewChip("gpiochip0", gpiocdev.WithConsumer(configs.Conf.AppName))
	if err != nil {
		return fmt.Errorf("Error opening GPIO chip: %w", err)
	}

	Chip = c

	select {
	case <-ctx.Done():
		log.Println("Closing GPIO chip due to context cancellation")
		_ = c.Close()
		return ctx.Err()
	default:
		log.Println("GPIO chip initialized")
	}
	return nil
}

func Initialize(ctx context.Context) error {
	var initErr error
	once.Do(func() {
		initErr = initializeChip(ctx)
	})
	return initErr
}

func CheckChip() bool {
	mu.RLock()
	defer mu.RUnlock()
	return Chip != nil
}

func setPinMode(pin dto.PinMode, asOutput bool) error {
	mu.Lock()
	defer mu.Unlock()

	if lines[pin.Pin] != nil {
		_ = lines[pin.Pin].Close()
	}

	var l *gpiocdev.Line
	var err error

	if asOutput {
		l, err = Chip.RequestLine(pin.Pin, gpiocdev.AsOutput(pin.Value))
	} else {
		l, err = Chip.RequestLine(pin.Pin, gpiocdev.AsInput)
	}

	if err != nil {
		return fmt.Errorf("GPIO: Error requesting line for pin %d: %w", pin.Pin, err)
	}

	lines[pin.Pin] = l

	err = db.SetPin(pin)
	if err != nil {
		return fmt.Errorf("GPIO: Error setting pin value in database: %w", err)
	}

	if !asOutput {
		val, err := l.Value()
		if err != nil {
			return fmt.Errorf("GPIO: Error reading pin %d value: %w", pin.Pin, err)
		}
		log.Printf("Pin %d value: %d", pin.Pin, val)
	}

	log.Printf("GPIO: Pin %d set to %d", pin.Pin, pin.Value)
	return nil
}

func setOutput(pin dto.PinMode) error {
	return setPinMode(pin, true)
}

func setInput(pin dto.PinMode) error {
	return setPinMode(pin, false)
}

func SetBool(pin dto.PinMode) error {
	if pin.Direction != "in" && pin.Direction != "out" {
		return errors.New("invalid direction")
	} else if pin.Direction == "out" {
		return setOutput(pin)
	} else {
		return setInput(pin)
	}
}

func GetAll() (db.Map, error) {
	if !CheckChip() {
		return nil, errors.New("GPIO chip not initialized")
	}

	return db.GetAllPin()
}

func GetGpioAll() (map[int]string, error) {
	if !CheckChip() {
		return nil, errors.New("GPIO chip not initialized")
	}

	mu.RLock()
	defer mu.RUnlock()

	numLines := Chip.Lines()
	usedPins := make(map[int]string)

	for offset := 0; offset < numLines; offset++ {
		info, err := Chip.LineInfo(offset)
		if err != nil {
			log.Printf("Error retrieving line info for line %d: %v\n", offset, err)
			continue
		}

		if info.Consumer != "" {
			usedPins[offset] = info.Consumer
		}
	}

	return usedPins, nil
}
