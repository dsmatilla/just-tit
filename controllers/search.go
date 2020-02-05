package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"sync"
)

type SearchController struct {
	beego.Controller
}

func (c *SearchController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	search := strings.Replace(aux, "/", "", -1)
	result := doSearch(search)

	c.Data["PageTitle"] = fmt.Sprintf("Search results for %s", search)
	c.Data["Result"] = result
	c.Data["PageMetaDesc"] = fmt.Sprintf("Search results for %s", search)
	c.Data["Search"] = search
	c.TplName = "index.html"
}

type searchResult struct {
	Pornhub PornhubSearchResult
	Redtube RedtubeSearchResult
	Tube8   Tube8SearchResult
	Youporn YoupornSearchResult
	Flag    bool
}

var waitGroup sync.WaitGroup

func doSearch(search string) searchResult {
	//var cached searchResult
	//if err := cache.Get(search, &cached); err == nil {
	//	return cached
	//} else {
		//log.Print("Cache NOT found")
		waitGroup.Add(4)

		PornhubChannel := make(chan PornhubSearchResult)
		RedtubeChannel := make(chan RedtubeSearchResult)
		Tube8Channel := make(chan Tube8SearchResult)
		YoupornChannel := make(chan YoupornSearchResult)

		go searchPornhub(search, PornhubChannel)
		go searchRedtube(search, RedtubeChannel)
		go searchTube8(search, Tube8Channel)
		go searchYouporn(search, YoupornChannel)
		result := searchResult{<-PornhubChannel, <-RedtubeChannel, <-Tube8Channel, <-YoupornChannel, true}

		waitGroup.Wait()

		//go cache.Set(search, result, 5 * time.Minute)
		return result
	//}
}

func searchPornhub(search string, c chan PornhubSearchResult) {
	defer waitGroup.Done()
	var result PornhubSearchResult
	result = PornhubSearchVideos(search)
	c <- result
	close(c)
}

func searchRedtube(search string, c chan RedtubeSearchResult) {
	defer waitGroup.Done()
	var result RedtubeSearchResult
	result = RedtubeSearchVideos(search)
	c <- result
	close(c)
}

func searchTube8(search string, c chan Tube8SearchResult) {
	defer waitGroup.Done()
	var result Tube8SearchResult
	result = Tube8SearchVideos(search)
	c <- result
	close(c)
}

func searchYouporn(search string, c chan YoupornSearchResult) {
	defer waitGroup.Done()
	var result YoupornSearchResult
	result = YoupornSearchVideos(search)
	c <- result
	close(c)
}