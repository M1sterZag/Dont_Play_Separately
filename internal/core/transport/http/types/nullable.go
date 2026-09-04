package core_http_types

import (
	"encoding/json"

	"github.com/M1sterZag/Dont_Play_Separately/internal/core/domain"
)

// Nullable — транспортная обёртка над domain.Nullable.
// Позволяет при декодировании JSON различить: поле отсутствует / null / значение.
type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
