package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/ntamares/TrackWeaver/internal/output"
	"github.com/ntamares/TrackWeaver/internal/wayback"
)

func main() {
	outputPath := flag.String(
		"out",
		"data/raw/cdx_records.txt",
		"output file path",
	)
	flag.Parse()

	writer, err := output.NewFileWriter(*outputPath)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	prefixes := []string{
		"fm949sd.com/chill/Story.aspx?ID=",
	}

	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup

	for _, prefix := range prefixes {
		records, err := wayback.FetchCDXRecords(prefix)
		if err != nil {
			fmt.Println("CDX error: ", err)
			continue
		}

		fmt.Println("Writing to:", *outputPath)

		for _, record := range records {
			wg.Add(1)

			go func(r wayback.CDXRecord) {
				defer wg.Done()

				sem <- struct{}{}        // acquire
				defer func() { <-sem }() // release

				body, err := wayback.FetchArchivedPage(r)
				if err != nil {
					fmt.Println("fetch error:", err)
					return
				}

				// TODO parser.ParseStory(body)
				line := fmt.Sprintf("%s | %s", r.Timestamp, r.Original)

				if err := writer.WriteLine(line); err != nil {
					fmt.Println("write error:", err)
				}

				_ = body // TODO later: parse this
			}(record)
		}
	}

	wg.Wait()
}
