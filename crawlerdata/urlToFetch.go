package crawlerdata

import "time"

type UrlToFetch struct {
	Id        int64     `json:"id" bson:"id"`
	Url       string    `json:"url" bson:"url"`
	Interval  int       `json:"interval" bson:"interval"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
