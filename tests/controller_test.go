package controllers

import (
	"just-tit/controllers"
	"testing"

	"github.com/astaxie/beego/cache"
)

func checkSearchResponse(t *testing.T, resp []controllers.JTVideo) {
	if len(resp) > 0 {
		t.Logf("ID: %+v\n", resp[0].ID)
		t.Logf("Domain: %+v\n", resp[0].Domain)
		t.Logf("Title: %+v\n", resp[0].Title)
		t.Logf("Description: %+v\n", resp[0].Description)
		t.Logf("Thumb: %+v\n", resp[0].Thumb)
		t.Logf("Thumbs: %+v\n", resp[0].Thumbs)
		t.Logf("Embed: %+v\n", resp[0].Embed)
		t.Logf("URL: %+v\n", resp[0].URL)
		t.Logf("Provider: %+v\n", resp[0].Provider)
		t.Logf("Rating: %+v\n", resp[0].Rating)
		t.Logf("Ratings: %+v\n", resp[0].Ratings)
		t.Logf("Duration: %+v\n", resp[0].Duration)
		t.Logf("Views: %+v\n", resp[0].Views)
		t.Logf("Width: %+v\n", resp[0].Width)
		t.Logf("Height: %+v\n", resp[0].Height)
		t.Logf("Segment: %+v\n", resp[0].Segment)
		t.Logf("PublishDate: %+v\n", resp[0].PublishDate)
		t.Logf("Type: %+v\n", resp[0].Type)
		t.Logf("Tags: %+v\n", resp[0].Tags)
		t.Logf("Categories: %+v\n", resp[0].Categories)
		t.Logf("Pornstars: %+v\n", resp[0].Pornstars)
		t.Logf("ExternalURL: %+v\n", resp[0].ExternalURL)
		t.Logf("ExternalID: %+v\n", resp[0].ExternalID)	
	}
	t.Logf("Total results: %+v\n", len(resp))
}

// TestPornhubSearch Test PornhubSearch method
func TestPornhubSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.PornhubSearch(search)
	checkSearchResponse(t, result)
}

func TestRedtubeSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.RedtubeSearch(search)
	checkSearchResponse(t, result)
}

func TestTube8Search(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.Tube8Search(search)
	checkSearchResponse(t, result)
}

func TestYoupornSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.YoupornSearch(search)
	checkSearchResponse(t, result)
}

func TestXtubeSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.XtubeSearch(search)
	checkSearchResponse(t, result)
}

func TestSpankwireSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.SpankwireSearch(search)
	checkSearchResponse(t, result)
}

func TestExtremetubeSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.ExtremetubeSearch(search)
	checkSearchResponse(t, result)
}

func TestKeezmoviesSearch(t *testing.T) {
	controllers.JTCache, _ = cache.NewCache("memory", `{"interval":60}`)
	var search = "amateur"
	result := controllers.KeezmoviesSearch(search)
	checkSearchResponse(t, result)
}