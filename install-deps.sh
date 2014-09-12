#!/bin/bash

function main {
    git clone https://github.com/lightning/liblightning.git
    cd liblightning
    ./install-opus.sh
    ./install-check.sh
    make
    sudo make install
    go get github.com/bmizerany/assert        \
        github.com/gorilla/mux                \
        github.com/gorilla/websocket          \
        github.com/hypebeast/go-osc/osc
}

main "$@"

