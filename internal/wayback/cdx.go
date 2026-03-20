package wayback

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchCDXRecords(prefix string) ([]CDXRecord, error) {
	url := fmt.Sprintf(
		"https://web.archive.org/cdx/search/cdx?url=%s&matchType=prefix&output=json&fl=original,timestamp,mimetype&filter=mimetype:text/html",
		prefix,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw [][]string
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var records []CDXRecord
	for i, row := range raw {
		if i == 0 {
			continue
		}

		records = append(records, CDXRecord{
			Original:  row[0],
			Timestamp: row[1],
			MimeType:  row[2],
		})
	}

	return records, nil
}
