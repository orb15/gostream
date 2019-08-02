// Package opts handles command line options
package opts

import "errors"

// IndexChannel defines the ElasticSearch index to which messages can be written. Indirectly,
// the IndexChannel defines the 'shape' (format, contents) of the message written to elastic
type IndexChannel int
const (

	// INDEX_TEST is a default message format
	INDEX_TEST = iota
)

// OptionsDef defines the options for the program
type OptionsDef struct {
	Index IndexChannel	// the elastic index to write messages
	Count int			// the number of messages to write to the index
}

// StringToIndexChannel converts a provided Index Channel name (as a string) into the actual IndexChannel
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