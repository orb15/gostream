package main

import (
"os"
"log"
"strconv"

"github.com/orb15/gostream/opts"
"github.com/orb15/gostream/msggen"
"github.com/orb15/gostream/elastic"
)

func main() {

	log.Println("Execution beginning...")
	
	//fetch command line args
	options := getOptionsFromArgs(os.Args[1:])
	
	//get the appropriate message generator
	gen := msggen.GetInstance(options.Index)

	//use it
	elastic.Fill(options, gen)
	

	log.Println("Execution complete")
}

func getOptionsFromArgs(args []string) opts.OptionsDef {

	options := opts.OptionsDef{}

	//check for proper number of arguments
	if(len(args) < 2) {
		log.Fatal("Too few arguments")
	}

	//NaN
	count, err := strconv.Atoi(args[1])
	if(err != nil){
		log.Fatalf("Number of messages is not actually a valid number: %s", args[2])
	}

	//too many messages?
	if(count <= 0 || count > 10000000) {
		log.Fatal("Do not request more than 10 million messages please")
	}

	//set up the options
	options.Count = count
	options.Index, err = opts.StringToIndexChannel(args[0])
	if(err != nil) {
		log.Fatalf("Unknown Index: %s", args[0])
	} 

	return options
}