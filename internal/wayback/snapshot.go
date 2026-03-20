package wayback

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func buildSnapshotURL(record CDXRecord) string {
	return fmt.Sprintf(
		"https://web.archive.org/web/%s/%s",
		record.Timestamp,
		record.Original,
	)
}

func FetchSnapshot(record CDXRecord) ([]byte, error) {
	url := buildSnapshotURL(record)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
