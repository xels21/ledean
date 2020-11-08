package mode

type Mode interface {
	Activate()
	Deactivate()
	Randomize()
}
