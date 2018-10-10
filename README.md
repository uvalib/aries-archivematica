# Aries API implementation for Archivematica

This is a web service to implement the Aries API for Archivematica.
It supports the following endpoints:

* / : returns version information
* /archivematica/[ID] : returns a JSON object with some information about the AIP referenced by ID

### System Requirements

* GO version 1.9.2 or greater
* DEP (https://golang.github.io/dep/) version 0.4.1 or greater
