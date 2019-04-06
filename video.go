package main

import (
	"encoding/json"
	"github.com/dsmatilla/extremetube"
	"github.com/dsmatilla/keezmovies"
	"github.com/dsmatilla/pornhub"
	"github.com/dsmatilla/redtube"
	"github.com/dsmatilla/spankwire"
	"github.com/dsmatilla/tube8"
	"github.com/dsmatilla/xtube"
	"github.com/dsmatilla/youporn"
	"sync"
)

type searchResult struct {
	Pornhub pornhub.PornhubSearchResult
	Redtube redtube.RedtubeSearchResult
	Tube8   tube8.Tube8SearchResult
	Youporn youporn.YoupornSearchResult
	Flag    bool
}

var waitGroup sync.WaitGroup

func pornhubGetVideoByID(videoID string) pornhub.PornhubSingleVideo {
	cachedElement := getFromDB("pornhub-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := pornhub.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("pornhub-video-"+videoID, string(json))
		return video
	}
}

func pornhubGetVideoEmbedCode(videoID string) pornhub.PornhubEmbedCode {
	cachedElement := getFromDB("pornhub-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := pornhub.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("pornhub-embed-"+videoID, string(json))
		return embed
	}
}

func redtubeGetVideoByID(videoID string) redtube.RedtubeSingleVideo {
	cachedElement := getFromDB("redtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := redtube.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("redtube-video-"+videoID, string(json))
		return video
	}
}

func redtubeGetVideoEmbedCode(videoID string) redtube.RedtubeEmbedCode {
	cachedElement := getFromDB("redtube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := redtube.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("redtube-embed-"+videoID, string(json))
		return embed
	}
}

func tube8GetVideoByID(videoID string) tube8.Tube8SingleVideo {
	cachedElement := getFromDB("tube8-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8SingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := tube8.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("tube8-video-"+videoID, string(json))
		return video
	}
}

func tube8GetVideoEmbedCode(videoID string) tube8.Tube8EmbedCode {
	cachedElement := getFromDB("tube8-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8EmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := tube8.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("tube8-embed-"+videoID, string(json))
		return embed
	}
}

func youpornGetVideoByID(videoID string) youporn.YoupornSingleVideo {
	cachedElement := getFromDB("youporn-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := youporn.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("youporn-video-"+videoID, string(json))
		return video
	}
}

func youpornGetVideoEmbedCode(videoID string) youporn.YoupornEmbedCode {
	cachedElement := getFromDB("youporn-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := youporn.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("youporn-embed-"+videoID, string(json))
		return embed
	}
}

func xtubeGetVideoByID(videoID string) xtube.XtubeVideo {
	cachedElement := getFromDB("xtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result xtube.XtubeVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := xtube.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("xtube-video-"+videoID, string(json))
		return video
	}
}

func spankwireGetVideoByID(videoID string) spankwire.SpankwireSingleVideo {
	cachedElement := getFromDB("spankwire-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := spankwire.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("spankwire-video-"+videoID, string(json))
		return video
	}
}

func spankwireGetVideoEmbedCode(videoID string) spankwire.SpankwireEmbedCode {
	cachedElement := getFromDB("spankwire-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := spankwire.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("spankwire-embed-"+videoID, string(json))
		return embed
	}
}

func keezmoviesGetVideoByID(videoID string) keezmovies.KeezmoviesSingleVideo {
	cachedElement := getFromDB("keezmovies-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := keezmovies.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("keezmovies-video-"+videoID, string(json))
		return video
	}
}

func keezmoviesGetVideoEmbedCode(videoID string) keezmovies.KeezmoviesEmbedCode {
	cachedElement := getFromDB("keezmovies-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := keezmovies.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("keezmovies-embed-"+videoID, string(json))
		return embed
	}
}

func extremetubeGetVideoByID(videoID string) extremetube.ExtremetubeSingleVideo {
	cachedElement := getFromDB("extremetube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeSingleVideo
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		video := extremetube.GetVideoByID(videoID)
		json, _ := json.Marshal(video)
		putToDB("extremetube-video-"+videoID, string(json))
		return video
	}
}

func extremetubeGetVideoEmbedCode(videoID string) extremetube.ExtremetubeEmbedCode {
	cachedElement := getFromDB("extremetube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeEmbedCode
		json.Unmarshal([]byte(cachedElement.Result), &result)
		return result
	} else {
		embed := extremetube.GetVideoEmbedCode(videoID)
		json, _ := json.Marshal(embed)
		putToDB("extremetube-embed-"+videoID, string(json))
		return embed
	}
}

func searchPornhub(search string, c chan pornhub.PornhubSearchResult) {
	defer waitGroup.Done()
	c <- pornhub.SearchVideos(search)
	close(c)
}

func searchRedtube(search string, c chan redtube.RedtubeSearchResult) {
	defer waitGroup.Done()
	c <- redtube.SearchVideos(search)
	close(c)
}

func searchTube8(search string, c chan tube8.Tube8SearchResult) {
	defer waitGroup.Done()
	c <- tube8.SearchVideos(search)
	close(c)
}

func searchYouporn(search string, c chan youporn.YoupornSearchResult) {
	defer waitGroup.Done()
	c <- youporn.SearchVideos(search)
	close(c)
}

func doSearch(search string) searchResult {
	waitGroup.Add(4)

	PornhubChannel := make(chan pornhub.PornhubSearchResult)
	RedtubeChannel := make(chan redtube.RedtubeSearchResult)
	Tube8Channel := make(chan tube8.Tube8SearchResult)
	YoupornChannel := make(chan youporn.YoupornSearchResult)

	go searchPornhub(search, PornhubChannel)
	go searchRedtube(search, RedtubeChannel)
	go searchTube8(search, Tube8Channel)
	go searchYouporn(search, YoupornChannel)

	result := searchResult{<-PornhubChannel, <-RedtubeChannel, <-Tube8Channel, <-YoupornChannel, true}

	waitGroup.Wait()
	return result
}
