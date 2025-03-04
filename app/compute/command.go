package compute

type Command int

const (
	GET Command = iota
	SET
	DEL
	UNKNOWN
)

func (c Command) String() string {
	switch c {
	case GET:
		return "GET"
	case SET:
		return "SET"
	case DEL:
		return "DEL"
	default:
		return "UNKNOWN"
	}
}

func (c Command) Int() int {
	switch c {
	case GET:
		return int(GET)
	case SET:
		return int(SET)
	case DEL:
		return int(DEL)
	default:
		return int(UNKNOWN)
	}
}
