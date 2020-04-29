package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bob620/qbitbaka/src/server/chewyroll"
	"github.com/bob620/qbitbaka/src/server/qbit"
	"github.com/bob620/qbitbaka/src/server/routes"
	"github.com/bob620/qbitbaka/src/server/ws"
)

func main() {
	qBit := qbit.NewQBit()
	cr := chewyroll.MakeChewyroll()
	socket := ws.CreateWs()

	qRoute := routes.MakeQbitBakaRoute(qBit)
	cRoute := routes.MakeChewyrollRoute(cr)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// Read in and process the request body
		data := make([]byte, request.ContentLength)
		if _, err := io.ReadFull(request.Body, data); err != nil {
			http.ServeFile(writer, request, "public/error.html")
			fmt.Println("error:", err)
		}

		url := request.RequestURI

		// Root /
		switch {
		// /socket
		case strings.HasPrefix(url, "/sockets"):
			socket.Handler(writer, request)
			break
		// /assets
		case strings.HasPrefix(url, "/assets/"):
			http.ServeFile(writer, request, "public/"+request.URL.Path)
			break
		// /qbitbaka
		case strings.HasPrefix(url, "/qbitbaka"):
			url = url[9:]
			qRoute.QbitBaka(url, writer, request)
			break
		case strings.HasPrefix(url, "/chewyroll"):
			url = url[10:]
			cRoute.Chewyroll(url, writer, request)
			break
		default:
			http.Redirect(writer, request, "/qbitbaka", http.StatusTemporaryRedirect)
			break
		}
	})

	fmt.Println("Listening on :9080")
	_ = http.ListenAndServe(":9080", nil)
}
