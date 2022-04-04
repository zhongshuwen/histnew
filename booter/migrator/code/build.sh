#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BROWN='\033[0;33m'
NC='\033[0m'

BUILD_SUFFIX=${1}

printf "${BROWN}Compiling ${BUILD_SUFFIX}${NC}\n"

mkdir -p $ROOT/build${BUILD_SUFFIX}
zswhq-cpp -I /usr/local/Cellar/zswhq.cdt/1.7.0/opt/zswhq.cdt/include/zswhqlib/capi  ./migrator.cpp -o ./build/migrator.wasm


