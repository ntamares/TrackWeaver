package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

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

		for idx, record := range records {
			wg.Add(1)

			// todo remove after parsing HTML
			go func(idx int, r wayback.CDXRecord) {
				defer wg.Done()
				sem <- struct{}{}        // acquire
				defer func() { <-sem }() // release

				time.Sleep(500 * time.Millisecond)

				body, err := wayback.FetchSnapshot(r)
				if err != nil {
					fmt.Println("fetch error:", err)
					return
				}

				// TODO remove after test parsing
				if idx < 2 {
					filename := fmt.Sprintf("snapshot_%d.html", idx)
					if err := os.WriteFile(filename, body, 0644); err != nil {
						fmt.Println("write error:", err)
					}
				}

				// TODO parser.ParseStory(body)
				line := fmt.Sprintf("%s | %s", r.Timestamp, r.Original)

				if err := writer.WriteLine(line); err != nil {
					fmt.Println("write error:", err)
				}

			}(idx, record)
		}
	}

	wg.Wait()
}
