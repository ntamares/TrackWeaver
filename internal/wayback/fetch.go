package wayback

import (
	"fmt"
	"io"
	"net/http"
)

func FetchArchivedPage(record CDXRecord) ([]byte, error) {
	url := fmt.Sprintf(
		"https://web.archive.org/web/%s/%s",
		record.Timestamp,
		record.Original,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
