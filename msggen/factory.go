package msggen

import "github.com/orb15/gostream/opts"
import "log"

type MessageGenerator interface {
	GenerateMessage() string
	GenerateByteMessage() []byte
}

func GetInstance(indexChannel opts.IndexChannel) MessageGenerator {
	switch(indexChannel) {
		case opts.INDEX_TEST:
			return newTestChannel()
		default:
			log.Panicf("Unable to GetInstance for provided IndexChannel %v\nThis is a code error and should never happen", indexChannel)
	}
	return nil
}