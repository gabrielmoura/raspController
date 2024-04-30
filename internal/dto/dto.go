package dto

import "errors"

type PinMode struct {
	Pin       int    `json:"pin"`
	Value     int    `json:"value"`
	Direction string `json:"direction"` // in or out
	Active    string `json:"active"`    // low or hight
}

// Constants para direção e ativação válidos.
const (
	Input  = "in"
	Output = "out"
	Low    = "low"
	High   = "high"
)

// Validation valida a estrutura PinMode.
func (p *PinMode) Validation() error {
	if len(p.Direction) > 0 && p.Direction != Input && p.Direction != Output {
		return errors.New("invalid direction, use 'in' or 'out'")
	}
	if len(p.Active) > 0 && p.Active != Low && p.Active != High {
		return errors.New("invalid active state, use 'low' or 'high'")
	}
	return nil
}
