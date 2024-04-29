// Package db implementa funcionalidades para interagir com um banco de dados.
package db

import (
	"encoding/json"
	"errors"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/rosedblabs/rosedb/v2"
)

// DB representa o banco de dados.
var DB *rosedb.DB

// Initialize inicializa o banco de dados.
func Initialize() error {
	options := rosedb.DefaultOptions
	options.DirPath = configs.Conf.DBDir
	db, err := rosedb.Open(options)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()
	DB = db

	return nil
}

// Set insere um par chave-valor no banco de dados.
func Set(name, value interface{}) error {
	return DB.Put([]byte(name.(string)), []byte(value.(string)))
}

// Get recupera o valor associado a uma chave do banco de dados.
func Get(name string) (string, error) {
	value, err := DB.Get([]byte(name))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// SetPin define o valor de um pino no banco de dados.
func SetPin(pin int, value int) error {
	// Pode ser mesclado caso o gpio_list já exista.
	// Ou seja, se o gpio_list já existir, ele será mesclado com o novo valor.
	// Se não existir, ele será criado.
	// Busque o valor atual do gpio_list.
	currentValue, err := Get("gpio_list")
	if err != nil {
		// Se não existir, crie um novo.
		currentValue = "{}"
	}
	// Converta o valor atual para um mapa.
	var gpioList map[string]int
	_ = json.Unmarshal([]byte(currentValue), &gpioList)
	// Adicione o novo valor ao mapa.
	gpioList[string(pin)] = value
	// Converta o mapa para uma string.
	newValue, _ := json.Marshal(gpioList)
	// Salve o novo valor.
	return Set("gpio_list", string(newValue))
}

// GetPin obtém o valor de um pino do banco de dados.
func GetPin(pin int) (int, error) {
	// Busque o valor atual do gpio_list.
	currentValue, err := Get("gpio_list")
	if err != nil {
		return 0, err
	}
	// Converta o valor atual para um mapa.
	var gpioList map[string]int
	_ = json.Unmarshal([]byte(currentValue), &gpioList)
	// Busque o valor do pino.
	value, ok := gpioList[string(pin)]
	if !ok {
		return 0, errors.New("pin value not found")
	}
	return value, nil
}
