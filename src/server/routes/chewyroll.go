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
	default:
		http.ServeFile(writer, request, "public/index.html")
		break
	}
}

func (cRoute *ChewyrollRoute) ChewyrollApi(path string) []byte {
	switch {
	case strings.HasPrefix(path, "/series/lookup/"):
		url := strings.SplitAfter(path, "/series/lookup/")[1]
		output, err := json.Marshal(cRoute.chewyroll.SeriesLookup("https://crunchyroll.com/" + url))
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
	}
	return []byte("{}")
}
