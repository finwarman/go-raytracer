package raytracer

import "math"

// Vector3f
type Vector3f struct {
	X, Y, Z float64
}

func (a *Vector3f) Set(index int, value float64) {
	switch index {
	case 0:
		a.X = value
	case 1:
		a.Y = value
	case 2:
		a.Z = value
	default:
		panic("Index out of range")
	}
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

func (a Vector3f) Normalised() Vector3f {
	return a.Multiply(1.0 / a.Norm())
}

// alias for length
func (a Vector3f) Norm() float64 {
	return a.Length()
}

// Vector4f
type Vector4f struct {
	X, Y, Z, W float64
}

func (a *Vector4f) Set(index int, value float64) {
	switch index {
	case 0:
		a.X = value
	case 1:
		a.Y = value
	case 2:
		a.Z = value
	case 3:
		a.W = value
	default:
		panic("Index out of range")
	}
}

func (a Vector4f) Add(b Vector4f) Vector4f {
	return Vector4f{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
		a.W + b.W,
	}
}

func (a Vector4f) Sub(b Vector4f) Vector4f {
	return Vector4f{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
		a.W - b.W,
	}
}

func (a Vector4f) Multiply(s float64) Vector4f {
	return Vector4f{
		a.X * s,
		a.Y * s,
		a.Z * s,
		a.W * s,
	}
}

func (a Vector4f) Dot(b Vector4f) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func (a Vector4f) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Vector4f) Cross(b Vector4f) Vector4f {
	panic("cross product is not defined for 4-vectors")
}

func (a Vector4f) Normalised() Vector4f {
	return a.Multiply(1.0 / a.Norm())
}

// alias for length
func (a Vector4f) Norm() float64 {
	return a.Length()
}

// Multiply4 multiplies a 3D vector by a 4D vector, returning a new 3D vector
// The fourth component of the result may need to be divided by the
// result's W component to get the correct 3-vector result.
func (a Vector3f) Multiply4(b Vector4f) Vector3f {
	// Construct a 4-vector with W=0 to represent a 3D vector.
	a4 := Vector4f{X: a.X, Y: a.Y, Z: a.Z, W: 0}

	// Multiply the 4-vectors element-wise.
	result := Vector4f{
		X: a4.X * b.X,
		Y: a4.Y * b.Y,
		Z: a4.Z * b.Z,
		W: a4.W * b.W,
	}

	// Divide the first three components by the fourth component,
	// if the fourth component is not zero.
	if result.W != 0 {
		result.X /= result.W
		result.Y /= result.W
		result.Z /= result.W
	}

	// Return the result as a 3D vector.
	return Vector3f{X: result.X, Y: result.Y, Z: result.Z}
}

// Matrix4x4 is a 4x4 matrix
type Matrix4x4 [4][4]float64

// Multiply the vector by a 4x4 matrix
func (v Vector3f) MultiplyMatrix4x4(m [4][4]float64) Vector3f {
	// Create a Vector4f with w=1
	v4 := Vector4f{v.X, v.Y, v.Z, 1.0}

	// Multiply the 4-vector by the matrix
	res := Vector4f{}
	for i := 0; i < 4; i++ {
		res.Set(i, v4.X*m[i][0]+v4.Y*m[i][1]+v4.Z*m[i][2]+v4.W*m[i][3])
	}

	// Divide the resulting 4-vector by its fourth component
	return Vector3f{res.X / res.W, res.Y / res.W, res.Z / res.W}
}
