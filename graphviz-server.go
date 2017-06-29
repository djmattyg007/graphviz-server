package main

import (
    "flag"
    "fmt"
    "encoding/base64"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "strconv"
)

func create_img(graph string, response http.ResponseWriter) {
    cmd := exec.Command("dot", "-T", "png")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer stdin.Close()
    cmd.Stdout = response

    fmt.Printf("Creating graph... ")
    if err = cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }
    response.Header().Set("Content-type", "image/png")
    fmt.Fprintf(stdin, "%s", graph)
    stdin.Close()
    cmd.Wait()
    fmt.Println("done")
}

func handle_get(response http.ResponseWriter, request *http.Request) {
    graph_encoded := request.URL.Query().Get("graph")
    if graph_encoded == "" {
        response.WriteHeader(400)
        return
    }

    bgraph, err := base64.URLEncoding.DecodeString(graph_encoded)
    if err != nil {
        response.WriteHeader(400)
        fmt.Fprintf(response, "Error: %s\n", err)
        return
    }
    // I'm not sure why this is the easiest way to convert a byte array to a string
    graph := fmt.Sprintf("%s", bgraph)
    fmt.Println(graph)

    create_img(graph, response)
}

func handle_post(response http.ResponseWriter, request *http.Request) {
    bgraph, err := ioutil.ReadAll(request.Body)
    if err != nil {
        response.WriteHeader(400)
        return
    }
    // I'm not sure why this is the easiest way to convert a byte array to a string
    graph := fmt.Sprintf("%s", bgraph)
    if graph == "" {
        response.WriteHeader(400)
        return
    }

    fmt.Println(graph)

    create_img(graph, response)
}

func handle(response http.ResponseWriter, request *http.Request) {
    fmt.Println(request.Method)

    if request.Method == "GET" {
        handle_get(response, request)
    } else if request.Method == "POST" {
        handle_post(response, request)
    } else {
        response.WriteHeader(405)
        response.Header().Set("Allow", "GET, POST")
    }
}

func main() {
    var portNumber int
    flag.IntVar(&portNumber, "port", 0, "Port to listen on")
    flag.Parse()
    if portNumber == 0 {
        portNumberStr := os.Getenv("GS_PORT")
        if portNumberStr == "" {
            portNumber = 8000 // The default port number
        } else {
            portNumberStrConv, err := strconv.Atoi(portNumberStr)
            if err == nil {
                portNumber = portNumberStrConv
            } else {
                fmt.Println(err)
                os.Exit(2)
            }
        }
    }

    http.HandleFunc("/", handle)
    fmt.Println(fmt.Sprintf("Starting webserver on port %d", portNumber))
    http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
