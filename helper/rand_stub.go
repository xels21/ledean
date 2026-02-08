//go:build false
// +build false

package helper

const randUint64Mul = 2685821657736338717
const randFloatDiv = 4294967296.0

type Rand struct {
	state uint64
}

func NewRand(seed uint64) *Rand {
	if seed == 0 {
		seed = 1
	}
	return &Rand{state: seed}
}

func (r *Rand) Seed(seed uint64) {
	if seed == 0 {
		seed = 1
	}
	r.state = seed
}

func (r *Rand) Uint64() uint64 {
	x := r.state
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.state = x
	return x * randUint64Mul
}

func (r *Rand) Uint32() uint32 {
	return uint32(r.Uint64() >> 32)
}

func (r *Rand) Int() int {
	return int(r.Uint32())
}

func (r *Rand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.Uint32() % uint32(n))
}

func (r *Rand) Float64() float64 {
	return float64(r.Uint32()) / randFloatDiv
}

func (r *Rand) Float32() float32 {
	return float32(r.Uint32()) / float32(randFloatDiv)
}
