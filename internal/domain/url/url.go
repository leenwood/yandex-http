package url

import "time"

type Url struct {
	Id          uint64
	OriginalUrl string
	ShortUrl    string
	Date        time.Time
}
