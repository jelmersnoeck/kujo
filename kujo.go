package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jelmersnoeck/kujo/pkg/kujo"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	output, err := kujo.SuffixJobs(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output[:len(output)-1]))
}
