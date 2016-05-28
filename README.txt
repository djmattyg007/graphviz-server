Graphviz as a Service

This is a simple program to run a webserver that responds to all GET requests
that provide the correct requests data with rendered graphs in PNG format.

By default, the program listens on port 8000. It looks for a query string
parameter named "graph", where it expects to find a URL-encoded base64
representation of a graph in the dot language used by graphviz's dot program.

Here is an example:

curl -s http://localhost:8000/?graph=ZGlncmFwaCBncmFwaG5hbWUgewogICAgYSAtPiBiIC0-IGM7CiAgICBiIC0-IGQ7Cn0K

This will return a fully-formed PNG image with the graph rendered by graphviz.

This software is released into the public domain without any warranty.
