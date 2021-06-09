package controller

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TransportError struct {
	Status   int
	Response Response
}
