package dto

type PinMode struct {
	Pin       int    `json:"pin"`
	Value     int    `json:"value"`
	Direction string `json:"direction"` // in or out
}
