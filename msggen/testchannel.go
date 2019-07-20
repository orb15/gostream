package msggen

import "encoding/json"
import "math/rand"
import "log"
import "time"
import "github.com/gofrs/uuid"

type testChannel struct {
	Index string
	Id string
	RandomVal int
	Classification string
}

var classificationMap = map[int]string{
	0: "Cat-0",
	1: "Cat-1",
	2: "Cat-2",
	3: "Cat-3",
	4: "Cat-4",
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func newTestChannel() MessageGenerator {
	a := testChannel{Index: "INDEX_TEST"}
	return &a
}

func (a *testChannel) GenerateMessage() string {

	a.buildBaseMessage()
	return string(a.buildJsonBytes())
}

func (a *testChannel) GenerateByteMessage() []byte {
	a.buildBaseMessage()
	return a.buildJsonBytes()
}

func (a *testChannel) buildJsonBytes() []byte {

	jsonBytes, err := json.Marshal(a)
	if err != nil {
		log.Fatalf("Unable to marshal a testChannel instance: %v", err)
	}
	return jsonBytes
} 

func (a *testChannel) buildBaseMessage() {
	a.RandomVal = rnd.Intn(19) + 1
	a.Classification = classificationMap[rnd.Intn(5)]

	uuid, err := uuid.NewV4()
	if(err != nil){
		log.Fatalf("Unable to generate a UUID: %v", err)
	}

	a.Id = uuid.String()
}


