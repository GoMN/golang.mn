package models

import(
	"bootstrap"
	"services/conf"
	"services/meetup"
	"time"
)

type Page struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Members   []bootstrap.Member `json:"members"`
	Groups    []meetup.Group `json:"groups"`
	MapsKey   string `json:"mapsKey"`
	Year      int `json:"year"`
	Version   string `json:version`
}

func AppData() Page{
	appdata := Page{}
	appdata.Title = "Go(lang)MN"
	appdata.Subtitle = "Minnesota Go Language Meetup"
	appdata.MapsKey = conf.Config.Maps.Key
	// TODO: replace with proper versioning
	appdata.Version = conf.Config.Version + time.Now().String()
	appdata.Year = time.Now().Year()
	return appdata
}
