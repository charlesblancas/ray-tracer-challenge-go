package main

import "math"

type Tuple struct {
	X float64
	Y float64
	Z float64
	W float64
}

func tupleEqual(a Tuple, b Tuple) bool {
	return floatEqual(a.X, b.X) && floatEqual(a.Y, b.Y) && floatEqual(a.Z, b.Z) && floatEqual(a.W, b.W)
}

func tupleAdd(a Tuple, b Tuple) Tuple {
	return Tuple{a.X + b.X, a.Y + b.Y, a.Z + b.Z, a.W + b.W}
}

func tupleSubtract(a Tuple, b Tuple) Tuple {
	return Tuple{a.X - b.X, a.Y - b.Y, a.Z - b.Z, a.W - b.W}
}

func tupleScale(a Tuple, k float64) Tuple {
	return Tuple{a.X * k, a.Y * k, a.Z * k, a.W * k}
}

func tupleDivide(a Tuple, k float64) Tuple {
	return Tuple{a.X / k, a.Y / k, a.Z / k, a.W / k}
}

func tupleNegate(a Tuple) Tuple {
	return Tuple{-a.X, -a.Y, -a.Z, -a.W}
}

func point(X float64, Y float64, Z float64) Tuple {
	return Tuple{X, Y, Z, 1.0}
}

func vector(X float64, Y float64, Z float64) Tuple {
	return Tuple{X, Y, Z, 0.0}
}

func vectorMagnitude(a Tuple) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z + a.W*a.W)
}

func vectorNormalize(a Tuple) Tuple {
	m := vectorMagnitude(a)

	return Tuple{a.X / m, a.Y / m, a.Z / m, a.W / m}
}

func vectorDot(a Tuple, b Tuple) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W + b.W
}

func vectorCross(a Tuple, b Tuple) Tuple {
	return vector(
		a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X,
	)
}

func isPoint(t Tuple) bool {
	return t.W == 1.0
}

func isVector(t Tuple) bool {
	return t.W == 0.0
}
