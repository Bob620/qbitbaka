package main

type MessageInput struct {
	MessageType string `json:"type"`
	Data []byte `json:"data"`
}

type 