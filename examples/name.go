package main

// IName service generating names
type IName interface {
	Gen() string
}

// Name service implementaion
type Name struct {
	random IRandom
}

// NewName is a Name serivice implementation factory
func NewName(random IRandom) IName {
	return &Name{
		random: random,
	}
}

// Gen generate name
func (n *Name) Gen() string {
	return "Noskov " + n.random.Name()
}
