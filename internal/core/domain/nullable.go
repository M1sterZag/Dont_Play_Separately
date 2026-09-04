package domain

// Nullable представляет поле, которое может быть:
//   - не передано вовсе (Set == false);
//   - передано со значением null (Set == true, Value == nil);
//   - передано со значением (Set == true, Value != nil).
type Nullable[T any] struct {
	Value *T
	Set bool
}