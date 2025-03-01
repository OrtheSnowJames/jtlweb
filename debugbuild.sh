#!/bin/bash

# Navigate to the project directory
cd /home/james/Documents/jtlweb

# Build the project with debug symbols
go build -gcflags "all=-N -l" -o debugbuild