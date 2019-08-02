//Package elastic handles communication with ElasticSearch
package elastic

import (
"log"
"net/http"
"bytes"
"strings"

"github.com/orb15/gostream/opts"
"github.com/orb15/gostream/msggen"
)

const cNONBULKLOADLIMIT int = 100

//I was using localhost here, but received DNS errors when running at scale. Need
//to revist reasoning behind this when I get a second
const cURL string = "http://127.0.0.1:9200/"
const cJSONMIME = "application/json"

// Fill loads an elasticsearch instance with messages based on some MessageGenerator. 
// Fill automatically switches between using a one-POST-per-message approach and using
// elastic's bulk fill POST methods depending on the number of messages desired
func Fill(options opts.OptionsDef, gen msggen.MessageGenerator) {

  url := getBaseURL(options)
  log.Printf("Filling elasticsearch at: %v", url)

  totalErrors := 0
  if options.Count <= cNONBULKLOADLIMIT {
    log.Printf("Simple fill starting for: %v documents", options.Count)
    totalErrors = simpleFill(url, options.Count, gen)
  } else {
    log.Printf("Bulk fill starting for: %v documents", options.Count)
    totalErrors = bulkFill(url, options.Count, gen)
  }
  log.Printf("Fill complete. Total errors: %v", totalErrors)
}

//POSTs one message to elasticsearch
func simpleFill(url string, count int, gen msggen.MessageGenerator) int {

  totalErrors := 0
  for i:=1; i <= count; i++ {
    resp, err := http.Post(url, cJSONMIME, bytes.NewBuffer(gen.GenerateByteMessage()))
    if err != nil {
      log.Printf("Unable to POST new document: %v", err)
      totalErrors++
    } else if resp.StatusCode >= 400 {
      log.Printf("Bad Request! HTTP Status code: %v", resp.StatusCode)
      totalErrors++
    }
  }
  return totalErrors
}


//fills elastic with a bulk load
func bulkFill(url string, count int, gen msggen.MessageGenerator) int {
  return 0
}

func getBaseURL(options opts.OptionsDef) string {
  indexAsString, err := opts.IndexChannelToString(options.Index)
  if err != nil {
    log.Fatalf("Unable to convert Index to string: %v", err)
  }
  return cURL + strings.ToLower(indexAsString) + "/msg"
}