package main

// IName service generating names
type IName interface {
	Gen() string
}

// Name service implementation
type Name struct {
	random IRandom
}

// NewName is a IName factory
func NewName(randoms []IRandom) (IName, error) {
	return &Name{
		random: randoms[0],
	}, nil
}

// Gen generate name
func (n *Name) Gen() string {
	return "Noskov " + n.random.Name()
}
