package mode

type Mode interface {
	Activate()
	Deactivate()
	Randomize()
	SetParameter(parm interface{})
	GetParameterJson() []byte
	GetFriendlyName() string
}
