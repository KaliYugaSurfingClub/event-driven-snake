package snake

type Direction int8

const (
	DirectionUp Direction = iota + 1
	DirectionRight
	DirectionDown
	DirectionLeft
)

func isOppositeDirections(a, b Direction) bool {
	return a%2 == b%2
}
