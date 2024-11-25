package hostState

type State int

const (
	Starting State = 0
	Started  State = 1
	Stopping State = 2
	Stopped  State = 3
)

func (state State) String() string {
	switch state {
	case Starting:
		return "Starting"
	case Started:
		return "Started"
	case Stopping:
		return "Stopping"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}
