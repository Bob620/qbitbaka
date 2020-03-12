package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func main() {
	qBit := NewQBit()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// Read in and process the request body
		data := make([]byte, request.ContentLength)
		if _, err := io.ReadFull(request.Body, data); err != nil {
			http.ServeFile(writer, request, "public/error.html")
			fmt.Println("error:", err)
		}

		switch {
		case strings.HasPrefix(request.RequestURI, "/assets/"):
			http.ServeFile(writer, request, "public/"+request.URL.Path)
			break
		case strings.HasPrefix(request.RequestURI, "/qbitbaka"):
			http.ServeFile(writer, request, "public/index.html")
			break
		case strings.HasPrefix(request.RequestURI, "/api/data/torrents/"):
			writer.Header().Add("Content-Type", "application/json")
			hash := strings.SplitAfter(request.RequestURI, "/api/data/torrents/")[1]
			output, err := json.Marshal(qBit.GetTorrent(hash))
			if err != nil {
				io.WriteString(writer, "{}")
			} else {
				writer.Write(output)
			}
			break
		case strings.HasPrefix(request.RequestURI, "/api/data"):
			output, err := json.Marshal(qBit.GetData())
			if err != nil {
				io.WriteString(writer, "{}")
			} else {
				writer.Write(output)
			}
			break
		default:
			http.Redirect(writer, request, "/qbitbaka", http.StatusTemporaryRedirect)
			break
		}
	})

	// Spawn goroutine to handle websockets
	go func() {
		ln, _ := net.Listen("tcp", ":9081")

		for {
			conn, _ := ln.Accept()

			// Spawn goroutine for each connection made to the websocket
			go func() {
				for {
					data := make([]byte, 1024)

					length, err := conn.Read(data)
					if err != nil {
						// Kill a dead connection or restart an existing one from an EOF error loop
						// It should work for the provided python3 client
						_ = conn.Close()
						break
					}

					// If we have a request
					if length > 0 {
						json.Unmarshal(data)
					}
				}
			}()
		}
	}()

	fmt.Println("Listening on :9080")
	_ = http.ListenAndServe(":9080", nil)
}
