package fn

import "math"

func GetManhattenDistance(fromX int, fromY int, toX int, toY int) int {
	return int(math.Abs(float64(fromX-toX))) + int(math.Abs(float64(fromY-toY)))
}
