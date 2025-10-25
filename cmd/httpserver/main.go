package main

import (
	// "fmt"
	"ibrahemassa/http_bootdev/internal/headers"
	"ibrahemassa/http_bootdev/internal/request"
	"ibrahemassa/http_bootdev/internal/response"
	"ibrahemassa/http_bootdev/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func funnyHandler(w *response.Writer, req *request.Request){
	badRequest := []byte(`
	<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>
	`)

	internalServerError := []byte(`
	<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>
	`)

	successResponse := []byte(`
	<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>
	`)
	var statusCode response.StatusCode
	var h headers.Headers 
	var body []byte
	chunked := false
	if req.RequestLine.RequestTarget == "/yourproblem"{
		statusCode = response.BadRequest
		body = badRequest
		// return &server.HandlerError{
		// 	StatusCode: response.BadRequest,
		// 	Body: "Your problem is not my problem\n",
		// }
	} else if req.RequestLine.RequestTarget == "/myproblem"{
		statusCode = response.InternalServerError
		body = internalServerError
		// return &server.HandlerError{
		// 	StatusCode: response.InternalServerError,
		// 	Body: "Woopsie, my bad\n",
		// }
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		res, err := http.Get("https://httpbin.org/" + req.RequestLine.RequestTarget[len("/httpbin/"):])
		if err != nil{
			log.Fatal(err)
		}
		statusCode = response.OK

		w.WriteStatusLine(response.OK)
		h = response.GetDefaultHeader(len(body), "text/plain")
		h.Set("transfer-encodin", "chunked")
		h.Delete("content-length")
		w.WriteHeaders(h)
		for {
			data := make([]byte, 32)
			_, err := res.Body.Read(data)
			if err != nil{
				break
			}
			w.WriteChunckedBody(data)
		}
		w.WriteChunckedBodyDone()

	} else{
		statusCode = response.OK
		body = successResponse
	}

	if chunked{
		return
	}
	h = response.GetDefaultHeader(len(body), "text/html")
	err := w.WriteStatusLine(statusCode)
	if err != nil{
		log.Fatal(err)
		return
	}

	err = w.WriteHeaders(h)
	if err != nil{
		log.Fatal(err)
		return
	}

	// _, err = w.WriteBody(body)
	_, err = w.WriteChunckedBody(body)
	if err != nil{
		log.Fatal(err)
		return
	}
}

func main() {
	server, err := server.Serve(port, funnyHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
