#!/bin/bash

function main {
    git clone https://github.com/lightning/liblightning.git
    cd liblightning
    ./install-opus.sh
    ./install-check.sh
    make
    sudo make install
    cd .. && ./install-pkg.sh
}

main "$@"

