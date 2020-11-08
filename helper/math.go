package helper

// Max returns the larger of x or y.
func MaxByte(x, y byte) byte {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func MinByte(x, y byte) byte {
	if x > y {
		return y
	}
	return x
}

// Max returns the larger of x or y.
func MaxFloat32(x, y float32) float32 {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func MinFloat32(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
}
