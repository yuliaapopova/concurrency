package compute

type Command uint32

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
