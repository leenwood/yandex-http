package url

import "time"

type Url struct {
	Id          string    `db:"id"`
	OriginalUrl string    `db:"original_url"`
	ClickCount  uint64    `db:"click_count"`
	CreatedDate time.Time `db:"created_date"`
}
