package richerror

// Type stores information that we want to show to user
type Type interface {
	String() string
}

type StringType string

func (et StringType) String() string {
	return string(et)
}
