package expression

type Status int64

const (
	Parsing Status = iota
	Error
	Waiting
	Ready
	InProgress
	Done
)

func (s Status) ToString() string {
	switch s {
	case Parsing:
		return "parsing"
	case Error:
		return "error"
	case Waiting:
		return "waiting"
	case Ready:
		return "ready"
	case InProgress:
		return "inProgress"
	case Done:
		return "done"
	}
	return "undefined"
}

func ToStatus(status string) Status {
	switch status {
	case "parsing":
		return Parsing
	case "error":
		return Error
	case "waiting":
		return Waiting
	case "ready":
		return Ready
	case "inProgress":
		return InProgress
	case "done":
		return Done
	}
	return -1
}
