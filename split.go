package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"sync"
)

var inputFile = flag.String("infile", "foo.xml", "Input file path")
var splitTagName = flag.String("tagname", "", "Tag name to Split on")
var outputFilePrefix = flag.String("outfileprefix", "foo-out", "Output file prefix")

type node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []node     `xml:",any"`
}

func writeNodeToFile(outfileName string, node node) {
	outfile, err := os.Create(outfileName)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer outfile.Close()
	bufOut := bufio.NewWriter(outfile)

	enc := xml.NewEncoder(bufOut)
	enc.Indent("  ", "    ")
	if err := enc.Encode(node); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
func main() {
	flag.Parse()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	recordCount := 0
	var wg sync.WaitGroup
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			var node node
			if inElement == *splitTagName {
				recordCount++
				decoder.DecodeElement(&node, &se)
				outfileName := fmt.Sprintf("%v_%v.xml", *outputFilePrefix, recordCount)
				fmt.Println("Writing ", *splitTagName, recordCount, "to file", outfileName)
				go func() {
					wg.Add(1)
					writeNodeToFile(outfileName, node)
					wg.Done()
				}()
			}
		default:
		}

	}
	wg.Wait()
	fmt.Printf("Total records: %d \n", recordCount)
}
