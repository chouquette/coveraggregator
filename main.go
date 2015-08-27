package main

import (
	"flag"
	"fmt"
	"github.com/chouquette/coveraggregator/aggregator"
	"os"
)

var (
	outputFile string
)

func init() {
	flag.StringVar(&outputFile, "o", "", "Aggregated profiles output file")
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 || outputFile == "" {
		fmt.Printf("usage: %s -o aggregatedprofile profile1 profile2 [profileN...]\n", os.Args[0])
		fmt.Println("At least 2 input profile files and an output file are required to run this tool")
		os.Exit(1)
	}
	p := aggregator.CoverProfile{}
	for _, f := range flag.Args() {
		fmt.Println("Aggregating", f)
		err := p.Aggregate(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
	fmt.Println("Writing aggregate to", outputFile)
	err := p.Write(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
