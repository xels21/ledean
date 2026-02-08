package helper

import "math/rand"

type Rand struct {
	r *rand.Rand
}

func NewRand(seed uint64) *Rand {
	return &Rand{r: rand.New(rand.NewSource(int64(seed)))}
}

func (r *Rand) Seed(seed uint64) {
	r.r.Seed(int64(seed))
}

func (r *Rand) Uint64() uint64 {
	return uint64(r.r.Int63())>>31 | uint64(r.r.Int63())<<32
}

func (r *Rand) Uint32() uint32 {
	return r.r.Uint32()
}

func (r *Rand) Int() int {
	return r.r.Int()
}

func (r *Rand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return r.r.Intn(n)
}

func (r *Rand) Float64() float64 {
	return r.r.Float64()
}

func (r *Rand) Float32() float32 {
	return r.r.Float32()
}
