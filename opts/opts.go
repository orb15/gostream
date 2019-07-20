package opts

import "errors"

type IndexChannel int
const (
	INDEX_TEST = iota
)

type OptionsDef struct {
	Index IndexChannel	// the elastic index to write messages
	Count int			// the number of messages to write to the index
}

func StringToIndexChannel(rawIndex string) (IndexChannel,error) {
	switch rawIndex {
		case "INDEX_TEST": return INDEX_TEST, nil
		default: return -1, errors.New("Index does not exist")
	}
}

func IndexChannelToString(rawIndex IndexChannel) (string, error) {
	switch rawIndex {
	case INDEX_TEST: return "INDEX_TEST", nil
	default: return "DOES NOT EXIST", errors.New("Should never happen")
	}
}