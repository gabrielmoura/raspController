// Package db implementa funcionalidades para interagir com um banco de dados.
package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/rosedblabs/rosedb/v2"
)

// DB representa o banco de dados.
var DB *rosedb.DB

type Map map[string]interface{}

// Initialize inicializa o banco de dados.
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

// Set insere um par chave-valor no banco de dados.
func Set(name, value interface{}) error {
	return DB.Put([]byte(name.(string)), []byte(value.(string)))
}
func SetJson(name string, value interface{}) error {
	jsonValue, _ := json.Marshal(value)
	return DB.Put([]byte(name), jsonValue)
}

// Get recupera o valor associado a uma chave do banco de dados.
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

// SetPin define o valor de um pino no banco de dados.
func SetPin(pin int, value int) error {
	gpios := make(Map)
	// Pode ser mesclado caso o gpio_list já exista.
	// Ou seja, se o gpio_list já existir, ele será mesclado com o novo valor.
	// Se não existir, ele será criado.
	// Busque o valor atual do gpio_list.
	err := GetJson("gpio_list", &gpios)
	if err != nil {
		// Se não existir, crie um novo.
		log.Println("DB: gpio_list not found")
	}

	// Adicione o novo valor ao mapa.
	gpios[string(rune(pin))] = value
	// Salve o novo valor.
	return SetJson("gpio_list", gpios)
}

// GetPin obtém o valor de um pino do banco de dados.
func GetPin(pin int) (int, error) {
	gpios := make(Map)
	err := GetJson("gpio_list", &gpios)
	if err != nil {
		return 0, fmt.Errorf("error getting pin %d: %s", pin, err.Error())
	}
	value, ok := gpios[string(rune(pin))]
	if !ok {
		return 0, errors.New("pin not found")
	}
	return value.(int), nil
}

func GetAllPin() (Map, error) {
	gpios := make(Map)
	err := GetJson("gpio_list", &gpios)
	if err != nil {
		return nil, err
	}
	return gpios, nil
}
