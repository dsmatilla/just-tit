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
	"log"
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
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][PORNHUB_GET]", err)
		}
		return result
	} else {
		video := pornhub.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("pornhub-video-"+videoID, string(jsonResult))
		return video
	}
}

func pornhubGetVideoEmbedCode(videoID string) pornhub.PornhubEmbedCode {
	cachedElement := getFromDB("pornhub-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result pornhub.PornhubEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][PORNHUB_EMBED]", err)
		}
		return result
	} else {
		embed := pornhub.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("pornhub-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func redtubeGetVideoByID(videoID string) redtube.RedtubeSingleVideo {
	cachedElement := getFromDB("redtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][REDTUBE_GET]", err)
		}
		return result
	} else {
		video := redtube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("redtube-video-"+videoID, string(jsonResult))
		return video
	}
}

func redtubeGetVideoEmbedCode(videoID string) redtube.RedtubeEmbedCode {
	cachedElement := getFromDB("redtube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result redtube.RedtubeEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][REDTUBE_EMBED]", err)
		}
		return result
	} else {
		embed := redtube.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("redtube-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func tube8GetVideoByID(videoID string) tube8.Tube8SingleVideo {
	cachedElement := getFromDB("tube8-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8SingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][TUBE8_GET]", err)
		}
		return result
	} else {
		video := tube8.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("tube8-video-"+videoID, string(jsonResult))
		return video
	}
}

func tube8GetVideoEmbedCode(videoID string) tube8.Tube8EmbedCode {
	cachedElement := getFromDB("tube8-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result tube8.Tube8EmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][TUBE8_EMBED]", err)
		}
		return result
	} else {
		embed := tube8.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("tube8-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func youpornGetVideoByID(videoID string) youporn.YoupornSingleVideo {
	cachedElement := getFromDB("youporn-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][YOUPORN_GET]", err)
		}
		return result
	} else {
		video := youporn.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("youporn-video-"+videoID, string(jsonResult))
		return video
	}
}

func youpornGetVideoEmbedCode(videoID string) youporn.YoupornEmbedCode {
	cachedElement := getFromDB("youporn-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result youporn.YoupornEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][YOUPORN_EMBED]", err)
		}
		return result
	} else {
		embed := youporn.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("youporn-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func xtubeGetVideoByID(videoID string) xtube.XtubeVideo {
	cachedElement := getFromDB("xtube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result xtube.XtubeVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][XTUBE_GET]", err)
		}
		return result
	} else {
		video := xtube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("xtube-video-"+videoID, string(jsonResult))
		return video
	}
}

func spankwireGetVideoByID(videoID string) spankwire.SpankwireSingleVideo {
	cachedElement := getFromDB("spankwire-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SPANKWIRE_GET]", err)
		}
		return result
	} else {
		video := spankwire.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("spankwire-video-"+videoID, string(jsonResult))
		return video
	}
}

func spankwireGetVideoEmbedCode(videoID string) spankwire.SpankwireEmbedCode {
	cachedElement := getFromDB("spankwire-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result spankwire.SpankwireEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SPANKWIRE_EMBED]", err)
		}
		return result
	} else {
		embed := spankwire.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("spankwire-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func keezmoviesGetVideoByID(videoID string) keezmovies.KeezmoviesSingleVideo {
	cachedElement := getFromDB("keezmovies-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][KEEZMOVIES_GET]", err)
		}
		return result
	} else {
		video := keezmovies.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("keezmovies-video-"+videoID, string(jsonResult))
		return video
	}
}

func keezmoviesGetVideoEmbedCode(videoID string) keezmovies.KeezmoviesEmbedCode {
	cachedElement := getFromDB("keezmovies-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result keezmovies.KeezmoviesEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][KEEZMOVIES_EMBED]", err)
		}
		return result
	} else {
		embed := keezmovies.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("keezmovies-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func extremetubeGetVideoByID(videoID string) extremetube.ExtremetubeSingleVideo {
	cachedElement := getFromDB("extremetube-video-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeSingleVideo
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][EXTREMETUBE_GET]", err)
		}
		return result
	} else {
		video := extremetube.GetVideoByID(videoID)
		jsonResult, _ := json.Marshal(video)
		putToDB("extremetube-video-"+videoID, string(jsonResult))
		return video
	}
}

func extremetubeGetVideoEmbedCode(videoID string) extremetube.ExtremetubeEmbedCode {
	cachedElement := getFromDB("extremetube-embed-" + videoID)
	if (JustTitCache{}) != cachedElement {
		var result extremetube.ExtremetubeEmbedCode
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][EXTREMETUBE_EMBED]", err)
		}
		return result
	} else {
		embed := extremetube.GetVideoEmbedCode(videoID)
		jsonResult, _ := json.Marshal(embed)
		putToDB("extremetube-embed-"+videoID, string(jsonResult))
		return embed
	}
}

func searchPornhub(search string, c chan pornhub.PornhubSearchResult) {
	defer waitGroup.Done()

	var result pornhub.PornhubSearchResult
	/*cachedElement := getFromDB("pornhub-search-" + search)
	if (JustTitCache{}) != cachedElement {
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SEARCH_PORNHUB]", err)
		}
		c <- result
	} else {
		result = pornhub.SearchVideos(search)
		jsonResult, _ := json.Marshal(result)
		putToDB("pornhub-search-"+search, string(jsonResult))
		c <- result
	}*/

	result = pornhub.SearchVideos(search)
	c <- result

	close(c)
}

func searchRedtube(search string, c chan redtube.RedtubeSearchResult) {
	defer waitGroup.Done()

	var result redtube.RedtubeSearchResult
	/*cachedElement := getFromDB("redtube-search-" + search)
	if (JustTitCache{}) != cachedElement {
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SEARCH_REDTUBE]", err)
		}
		c <- result
	} else {
		result = redtube.SearchVideos(search)
		jsonResult, _ := json.Marshal(result)
		putToDB("redtube-search-"+search, string(jsonResult))
		c <- result
	}*/


	result = redtube.SearchVideos(search)
	c <- result

	close(c)
}

func searchTube8(search string, c chan tube8.Tube8SearchResult) {
	defer waitGroup.Done()

	var result tube8.Tube8SearchResult
	/*cachedElement := getFromDB("tube8-search-" + search)
	if (JustTitCache{}) != cachedElement {
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SEARCH_TUBE8]", err)
		}
		c <- result
	} else {
		result = tube8.SearchVideos(search)
		jsonResult, _ := json.Marshal(result)
		putToDB("tube8-search-"+search, string(jsonResult))
		c <- result
	}*/

	result = tube8.SearchVideos(search)
	c <- result

	close(c)
}

func searchYouporn(search string, c chan youporn.YoupornSearchResult) {
	defer waitGroup.Done()

	var result youporn.YoupornSearchResult
	/*cachedElement := getFromDB("youporn-search-" + search)
	if (JustTitCache{}) != cachedElement {
		err := json.Unmarshal([]byte(cachedElement.Result), &result)
		if err != nil {
			log.Println("[JUST-TIT][SEARCH_YOUPORN]", err)
		}
		c <- result
	} else {
		result = youporn.SearchVideos(search)
		jsonResult, _ := json.Marshal(result)
		putToDB("youporn-search-"+search, string(jsonResult))
		c <- result
	}*/

	result = youporn.SearchVideos(search)
	c <- result

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
