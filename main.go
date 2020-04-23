package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bob620/qbitbaka/src/server/qbit"
	"github.com/bob620/qbitbaka/src/server/ws"
)

func main() {
	qBit := qbit.NewQBit()
	socket := ws.CreateWs()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// Read in and process the request body
		data := make([]byte, request.ContentLength)
		if _, err := io.ReadFull(request.Body, data); err != nil {
			http.ServeFile(writer, request, "public/error.html")
			fmt.Println("error:", err)
		}

		switch {
		case strings.HasPrefix(request.RequestURI, "/sockets"):
			socket.Handler(writer, request)
			break
		case strings.HasPrefix(request.RequestURI, "/assets/"):
			http.ServeFile(writer, request, "public/"+request.URL.Path)
			break
		case strings.HasPrefix(request.RequestURI, "/qbitbaka/api/data/torrents/"):
			writer.Header().Add("Content-Type", "application/json")
			hash := strings.SplitAfter(request.RequestURI, "/qbitbaka/api/data/torrents/")[1]
			output, err := json.Marshal(qBit.GetTorrent(hash))
			if err != nil {
				io.WriteString(writer, "{}")
			} else {
				writer.Write(output)
			}
			break
		case strings.HasPrefix(request.RequestURI, "/qbitbaka/api/data"):
			output, err := json.Marshal(qBit.GetData())
			if err != nil {
				io.WriteString(writer, "{}")
			} else {
				writer.Write(output)
			}
			break
		case strings.HasPrefix(request.RequestURI, "/qbitbaka"):
			http.ServeFile(writer, request, "public/index.html")
			break
		default:
			http.Redirect(writer, request, "/qbitbaka", http.StatusTemporaryRedirect)
			break
		}
	})

	fmt.Println("Listening on :9080")
	_ = http.ListenAndServe(":9080", nil)
}
