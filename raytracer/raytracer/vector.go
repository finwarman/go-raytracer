package raytracer

import "math"

// Vector3f
type Vector3f struct {
	X, Y, Z float64
}

func (a Vector3f) Add(b Vector3f) Vector3f {
	return Vector3f{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

func (a Vector3f) Sub(b Vector3f) Vector3f {
	return Vector3f{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
	}
}

func (a Vector3f) Multiply(s float64) Vector3f {
	return Vector3f{
		a.X * s,
		a.Y * s,
		a.Z * s,
	}
}

func (a Vector3f) Dot(b Vector3f) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector3f) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Vector3f) Cross(b Vector3f) Vector3f {
	return Vector3f{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func (a Vector3f) Norm() Vector3f {
	return a.Multiply(1.0 / a.Length())
}

// Vector2f
type Vector2f struct {
	X, Y float64
}

func (a Vector2f) Add(b Vector2f) Vector2f {
	return Vector2f{
		a.X + b.X,
		a.Y + b.Y,
	}
}

func (a Vector2f) Sub(b Vector2f) Vector2f {
	return Vector2f{
		a.X - b.X,
		a.Y - b.Y,
	}
}

func (a Vector2f) Multiply(s float64) Vector2f {
	return Vector2f{
		a.X * s,
		a.Y * s,
	}
}

func (a Vector2f) Dot(b Vector2f) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2f) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Vector2f) Norm() Vector2f {
	return a.Multiply(1.0 / a.Length())
}
