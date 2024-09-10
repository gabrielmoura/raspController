package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/internal/dto"
	"github.com/rosedblabs/rosedb/v2"
)

// DB represents the database
var DB *rosedb.DB

type Map map[string]interface{}
type PinMap map[int]dto.PinMode

// Initialize initializes the database.
func Initialize(ctx context.Context) error {
	options := rosedb.DefaultOptions
	options.DirPath = configs.Conf.DBDir
	db, err := rosedb.Open(options)
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		_ = db.Close()
		return ctx.Err()
	}
	DB = db

	return nil
}

// Set inserts a key-value pair into the database.
func Set(name, value interface{}) error {
	return DB.Put([]byte(name.(string)), []byte(value.(string)))
}
func SetJson(name string, value interface{}) error {
	jsonValue, _ := json.Marshal(value)
	return DB.Put([]byte(name), jsonValue)
}

// Get retrieves the value associated with a database key.
func Get(name string) (string, error) {
	value, err := DB.Get([]byte(name))
	if err != nil {
		return "", err
	}
	return string(value), nil
}
func GetJson(name string, value *Map) error {
	jsonValue, err := DB.Get([]byte(name))
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonValue, value)
}
func GetJsonPin(name string, value *PinMap) error {
	jsonValue, err := DB.Get([]byte(name))
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonValue, value)
}

// SetPin sets the value of a pin in the database.
func SetPin(pin dto.PinMode) error {
	gpios := make(PinMap)
	err := GetJsonPin("gpio_list", &gpios)
	if err != nil {
		// Se não existir, crie um novo.
		log.Println("DB: gpio_list not found")
	}

	gpios[pin.Pin] = pin

	return SetJson("gpio_list", gpios)
}

// GetPin gets the value of a pin from the database.
func GetPin(pin int) (int, error) {
	gpios := make(PinMap)
	err := GetJsonPin("gpio_list", &gpios)
	if err != nil {
		return 0, fmt.Errorf("error getting pin %d: %s", pin, err.Error())
	}
	value, ok := gpios[pin]
	if !ok {
		return 0, errors.New("pin not found")
	}
	return value.Pin, nil
}

func GetAllPin() (Map, error) {
	gpios := make(Map)
	err := GetJson("gpio_list", &gpios)
	if err != nil {
		return nil, err
	}
	return gpios, nil
}
