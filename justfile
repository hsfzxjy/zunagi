export GOBIN := "bin/"

set shell := ["bash", "-c"]

binary_name := "zunagi"

default:
    just --list

build version os:
    #!/usr/bin/env bash
    export GOOS={{os}}
    ext={{ if os == "windows" { ".exe" } else { "" } }}
    ldflags={{ if os == "windows" {"'-ldflags -H=windowsgui'"} else {""} }}
    go build -v -x $ldflags -o bin/{{binary_name}}-{{version}}-{{os}}$ext ./cmd/{{version}}

dev version os:
    #!/usr/bin/env bash
    export GOOS={{os}}
    ldflags={{ if os == "windows" { "'-ldflags -H=windowsgui'" } else { "" } }}
    go run -v -x $ldflags ./cmd/{{version}}