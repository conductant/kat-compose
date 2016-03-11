# kat-compose

Compose to Aurora

Current Status
==============

This works

   go run cmd/kat-compose.go load -p myproj -f examples/saas-config.yml

This will dump out the Compose yaml file in JSON, after the yaml file has been parsed into internal data structures.
The next step is to support flexible parsing of the V1 and V2 format, and possibly the V3 (Swarm V2) format.

Conversion to the Aurora JSON format will use templates to generate the input JSON for the Aurora client.