package url

import "time"

type Url struct {
	Id          string    `db:"id"`
	OriginalUrl string    `db:"original_url"`
	ClickCount  uint64    `db:"click_count"`
	Date        time.Time `db:"date"`
}
