package interfaces

type StringerE interface {
	String() (string, error)
}
type StringParser interface {
	FromString(string)
}
type StringParserE interface {
	FromString(string) error
}

type ArrayParser interface {
	FromArray([]string)
}
type ArrayParserE interface {
	FromArray([]string) error
}

type MapParser[T any] interface {
	FromMap(map[string]T)
}
type MapParserE[T any] interface {
	FromMap(map[string]T) error
}

func GetStringParser(iface any) func(string) error {
	if val, ok := iface.(StringParser); ok {
		return func(s string) error {
			val.FromString(s)
			return nil
		}
	} else if val, ok := iface.(StringParserE); ok {
		return val.FromString
	}
	return nil
}
