package interfaces

type Initializer interface {
	Initialize()
}

type Validator interface {
	Validate() error
}
