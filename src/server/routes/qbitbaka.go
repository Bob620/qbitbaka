package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bob620/qbitbaka/src/server/qbit"
)

type QbitBakaRoute struct {
	qBit *qbit.QBit
}

func MakeQbitBakaRoute(qbit *qbit.QBit) *QbitBakaRoute {
	return &QbitBakaRoute{qbit}
}

func (qRoute *QbitBakaRoute) QbitBaka(path string, writer http.ResponseWriter, request *http.Request) {
	switch {
	case strings.HasPrefix(path, "/api"):
		writer.Header().Add("Content-Type", "application/json")
		writer.Write(qRoute.QbitBakaApi(path[4:]))
		break
	default:
		http.ServeFile(writer, request, "html/qbit/index.html")
		break
	}
}

func (qRoute *QbitBakaRoute) QbitBakaApi(path string) []byte {
	switch {
	case strings.HasPrefix(path, "/data/torrents/"):
		hash := strings.SplitAfter(path, "/data/torrents/")[1]
		output, err := json.Marshal(qRoute.qBit.GetTorrent(hash))
		if err == nil {
			return output
		}
		break
	case strings.HasPrefix(path, "/data"):
		output, err := json.Marshal(qRoute.qBit.GetData())
		if err == nil {
			return output
		}
		break
	}
	return []byte("{}")
}
