package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bob620/qbitbaka/src/server/chewyroll"
)

type ChewyrollRoute struct {
	chewyroll *chewyroll.Chewyroll
}

func MakeChewyrollRoute(cr *chewyroll.Chewyroll) *ChewyrollRoute {
	return &ChewyrollRoute{cr}
}

func (cRoute *ChewyrollRoute) Chewyroll(path string, writer http.ResponseWriter, request *http.Request) {
	switch {
	case strings.HasPrefix(path, "/api"):
		writer.Header().Add("Content-Type", "application/json")
		writer.Write(cRoute.ChewyrollApi(path[4:]))
		break
		//	case strings.HasPrefix(path, "/"):
		//		writer.Header().Add("Content-Type", "text/html; charset=utf-8")
		//		break
	default:
		http.ServeFile(writer, request, "html/cr/index.html")
		break
	}
}

func (cRoute *ChewyrollRoute) ChewyrollApi(path string) []byte {
	switch {
	case strings.HasPrefix(path, "/series/search/name"):
		name := strings.SplitAfter(path, "/series/search/name/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesSearchByName(name))
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/series/search/id"):
		id := strings.SplitAfter(path, "/series/search/id/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesSearchById(id))
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/series/lookup/id"):
		id := strings.SplitAfter(path, "/series/lookup/id/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesLookupById(id))
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/series/lookup/uuid"):
		uuid := strings.SplitAfter(path, "/series/lookup/uuid/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesLookupByUuid(uuid))
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/series/lookup/url"):
		url := strings.SplitAfter(path, "/series/lookup/url/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesLookupByUrl("https://crunchyroll.com/" + url))
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/series/update/"):
		uuid := strings.SplitAfter(path, "/series/update/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesUpdate(uuid))
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/series/download/"):
		uuid := strings.SplitAfter(path, "/series/download/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesDownload(uuid))
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/episodes/lookup/"):
		uuid := strings.SplitAfter(path, "/episodes/lookup/")[1]
		output, err := json.Marshal(cRoute.chewyroll.EpisodesLookup(uuid))
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/episodes/update/"):
		uuid := strings.SplitAfter(path, "/episodes/update/")[1]
		output, err := json.Marshal(cRoute.chewyroll.EpisodesUpdate(uuid))
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/episodes/download/"):
		uuid := strings.SplitAfter(path, "/episodes/download/")[1]
		output, err := json.Marshal(cRoute.chewyroll.EpisodesDownload(uuid))
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/cr/listing/update"):
		output, err := json.Marshal(cRoute.chewyroll.CRListingUpdate())
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/cr/season/update"):
		output, err := json.Marshal(cRoute.chewyroll.CRSeasonUpdate())
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/cr/season/get"):
		output, err := json.Marshal(cRoute.chewyroll.CRSeasonGet())
		if err == nil {
			return output
		}

	case strings.HasPrefix(path, "/mal/season/update"):
		output, err := json.Marshal(cRoute.chewyroll.MalSeasonUpdate())
		if err == nil {
			return output
		}
	case strings.HasPrefix(path, "/mal/season/get"):
		output, err := json.Marshal(cRoute.chewyroll.MalSeasonGet())
		if err == nil {
			return output
		}
	}
	return []byte("{}")
}
