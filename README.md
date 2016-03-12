# kat-compose

Compose to Aurora / JSON / Thrift API

### Current Status

This works

    go run cmd/kat-compose.go load -p myproj -f examples/saas-config.yml

This will dump out the Compose yaml file in JSON, after the yaml file has been parsed into internal data structures.
The next step is to support flexible parsing of the V1 and V2 format, and possibly the V3 (Swarm V2) format.

Conversion to the Aurora JSON format will use templates to generate the input JSON for the Aurora client.

Given the docker-compose.yml file from the voting app,

    go run cmd/kat-compose.go -logtostderr convert -p myproj -f examples/voting-app.yml  -dump

will generate a JSON of multiple jobs for Aurora.

### TODO:

  1. Not parsing networks in V2
  2. Not tested against actual Aurora client yet
  3. Not calculating dependencies.

### Ideas:

   1. This may be best included in a server where a yml file is POST to the endpoint. The server then
   goes through a pipeline of handling special cases for production environment, like volumes, ports and networks.
   The user's input can just be what works in her local development environment.
   2. This should be posted in a single document as opposed to multiple calls to an API from the client.