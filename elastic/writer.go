package elastic

import "github.com/orb15/gostream/opts"
import "github.com/orb15/gostream/msggen"
import "log"
import "net/http"
import "bytes"
import "strings"
import "runtime"

const c_SINGLE_THREAD_THRESHOLD int = 100
const c_MULTITHREAD_THRESHOLD int = 250000
const c_THREADS_PER_CPU int = 8

const c_URL string = "http://127.0.0.1:9200/"
const c_JSON_MIME = "application/json"

func Fill(options opts.OptionsDef, gen msggen.MessageGenerator) {

  url := getBaseURL(options)
  log.Printf("Filling elasticsearch at: %v", url)

  totalErrors := 0
  if options.Count <= c_SINGLE_THREAD_THRESHOLD {
    log.Printf("Single thread fill starting for: %v documents", options.Count)
    totalErrors = singleThreadFill(url, options.Count, gen, nil)
  } else if options.Count < c_MULTITHREAD_THRESHOLD {
    log.Printf("Multithread fill starting for: %v documents", options.Count)
    totalErrors = multithreadFill(url, options.Count, gen)
  } else {
    log.Printf("Bulk fill starting for: %v documents", options.Count)
    totalErrors = bulkFill(url, options.Count, gen)
  }
  log.Printf("Fill complete. Total errors: %v", totalErrors)
}

func singleThreadFill(url string, count int, gen msggen.MessageGenerator, ch chan int) int {

  totalErrors := 0
  for i:=1; i <= count; i++ {
    resp, err := http.Post(url, c_JSON_MIME, bytes.NewBuffer(gen.GenerateByteMessage()))
    if err != nil {
      log.Printf("Unable to POST new document: %v", err)
      totalErrors++
    } else if resp.StatusCode >= 400 {
      log.Printf("Bad Request! HTTP Status code: %v", resp.StatusCode)
      totalErrors++
    }
  }

  //tell the controlling thread we are done (if neccessary)
  if ch != nil {
    ch <- totalErrors
  }

  return totalErrors
}

func multithreadFill(url string, count int, gen msggen.MessageGenerator) int {
  
  cpuCount := runtime.NumCPU()
  desiredThreadCount := cpuCount * c_THREADS_PER_CPU
  log.Printf("Spinning up %v threads for %v vCPUs", desiredThreadCount, cpuCount)

  //determine number of messages that give us an even workload
  unbalancedMsgs := count % desiredThreadCount
  msgsPerThread := (count - unbalancedMsgs) / desiredThreadCount
  log.Printf("Messages per thread: %v with extra: %v", msgsPerThread, unbalancedMsgs)

  //before gertting underway, build an array of channels our other threads can chat on
  channels := make([]chan int,desiredThreadCount)

  //start the other threads to handle the bulk of the messages
  //Note that the channel here has a buffer of 1, allowing it to run to completion and store
  //results when complete. 
  for i := 0; i < desiredThreadCount; i++ {
    channels[i] = make(chan int, 1)
    go singleThreadFill(url, msgsPerThread, gen, channels[i])
  }

  //use this (main) thread to handle the extra messages
  //this also gives this thread something to do rather than dropping to the 
  //wait and getting switched out
  allThreadSendErrors := 0
  if unbalancedMsgs > 0 {
    allThreadSendErrors = singleThreadFill(url, unbalancedMsgs, gen, nil)
  }

  //gather and sum all errors. Requirwes all channels to have reported before proceding
  for i := 0; i < desiredThreadCount; i++ {
    allThreadSendErrors += <- channels[i]
  }

  return allThreadSendErrors
  
}

func bulkFill(url string, count int, gen msggen.MessageGenerator) int {
  return 0
}

func getBaseURL(options opts.OptionsDef) string {
  indexAsString, err := opts.IndexChannelToString(options.Index)
  if err != nil {
    log.Fatalf("Unable to convert Index to string: %v", err)
  }
  return c_URL + strings.ToLower(indexAsString) + "/msg"
}