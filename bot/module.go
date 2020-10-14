package bot

type ModuleType int

const (
	MOTNil ModuleType = iota
	MOTPrivate
	MOTChannel
	MOTSuperGroup
	MOTGroup
)

func (ct ModuleType) String() string {
	switch ct {
	case MOTPrivate:
		return "private"
	case MOTChannel:
		return "channel"
	case MOTSuperGroup:
		return "supergroup"
	case MOTGroup:
		return "group"
	default:
		return ""
	}
}

func String2ModuleType(str string) ModuleType {
	switch str {
	case "private":
		return MOTPrivate
	case "channel":
		return MOTChannel
	case "supergroup":
		return MOTSuperGroup
	case "group":
		return MOTGroup
	default:
		return MOTNil
	}
}
