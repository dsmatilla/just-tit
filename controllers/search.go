package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"sync"
	"sort"
)

// SearchController Beego Controller
type SearchController struct {
	beego.Controller
}

// Get results of the search using different providers with go routines
func (c *SearchController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	search := strings.Replace(aux, "/", "", -1)
	result := doSearch(search)

	c.Data["PageTitle"] = fmt.Sprintf("Search results for %s", search)
	c.Data["Result"] = result
	c.Data["PageMetaDesc"] = fmt.Sprintf("Search results for %s", search)
	c.Data["Search"] = search

	//c.Data["video"] = video
	c.Data["debug"] = result
	c.Layout = "index.tpl"
	c.TplName = "search.tpl"
}

var waitGroup sync.WaitGroup

func doSearch(search string) []JTVideo {
	waitGroup.Add(5)

	ChannelPornhub := make(chan []JTVideo)
	ChannelRedtube := make(chan []JTVideo)
	ChannelYouporn := make(chan []JTVideo)
	ChannelTube8 := make(chan []JTVideo)
	ChannelKeezmovies := make(chan []JTVideo)

	go searchPornhub(search, ChannelPornhub)
	go searchRedtube(search, ChannelRedtube)
	go searchYouporn(search, ChannelYouporn)
	go searchTube8(search, ChannelTube8)
	go searchKeezmovies(search, ChannelKeezmovies)

	resultPornhub := <-ChannelPornhub
	resultRedtube := <-ChannelRedtube
	resultYouporn := <-ChannelYouporn
	resultTube8 := <-ChannelTube8
	resultKeezmovies := <-ChannelKeezmovies

	waitGroup.Wait()

	var result []JTVideo
	result = append(result, resultPornhub...)
	result = append(result, resultRedtube...)
	result = append(result, resultYouporn...)
	result = append(result, resultTube8...)
	result = append(result, resultKeezmovies...)

    sort.Slice(result, func(p, q int) bool {  
		return result[p].Rating > result[q].Rating }) 

	return result
}

func searchPornhub(search string, c chan []JTVideo) {
	defer waitGroup.Done()
	var result []JTVideo
	result = PornhubSearch(search)
	c <- result
	close(c)
}

func searchRedtube(search string, c chan []JTVideo) {
	defer waitGroup.Done()
	var result []JTVideo
	result = RedtubeSearch(search)
	c <- result
	close(c)
}

func searchYouporn(search string, c chan []JTVideo) {
	defer waitGroup.Done()
	var result []JTVideo
	result = YoupornSearch(search)
	c <- result
	close(c)
}

func searchTube8(search string, c chan []JTVideo) {
	defer waitGroup.Done()
	var result []JTVideo
	result = Tube8Search(search)
	c <- result
	close(c)
}

func searchKeezmovies(search string, c chan []JTVideo) {
	defer waitGroup.Done()
	var result []JTVideo
	result = KeezmoviesSearch(search)
	c <- result
	close(c)
}