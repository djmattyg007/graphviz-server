package main

import (
    "fmt"
    "encoding/base64"
    "net/http"
    "os/exec"
)

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

    cmd := exec.Command("dot", "-T", "png")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer stdin.Close()
    cmd.Stdout = response

    fmt.Printf("Creating graph...")
    if err = cmd.Start(); err != nil {
        response.WriteHeader(500)
        fmt.Println("An error occurred: ", err)
        return
    }
    fmt.Fprintf(stdin, "%s", graph)
    stdin.Close()
    cmd.Wait()
    fmt.Println("done")
}

func main() {
    http.HandleFunc("/", handle_get)
    fmt.Println("Starting webserver")
    http.ListenAndServe(":8000", nil)
}
