#!/usr/bin/env bash

SOLOOS_COMMON_PATH=/opt/soloos/common
# adjust GOPATH
case ":$GOPATH:" in
        *":$SOLOOS_COMMON_PATH:"*) :;;
        *) GOPATH=$SOLOOS_COMMON_PATH:$GOPATH;;
esac

SOLOOS_SDBONE_PATH=/opt/soloos/sdbone
# adjust GOPATH
case ":$GOPATH:" in
        *":$SOLOOS_SDBONE_PATH:"*) :;;
        *) GOPATH=$SOLOOS_SDBONE_PATH:$GOPATH;;
esac

# adjust GOPATH
case ":$GOPATH:" in
    *":$(pwd):"*) :;;
    *) GOPATH=$(pwd):$GOPATH;;
esac

export GOPATH

# adjust PATH
while IFS=':' read -ra ADDR; do
    for i in "${ADDR[@]}"; do
        case ":$PATH:" in
            *":$i/bin:"*) :;;
            *) PATH=$i/bin:$PATH
        esac
    done
done <<< "$GOPATH"
export PATH

if [ ! -d "$(pwd)/bin" ];
then
    mkdir -p "$(pwd)/bin"
fi

if [ ! -d "$(pwd)/src" ];
then
    mkdir -p "$(pwd)/src"
fi

# mkdir -p go.build/cache
# mkdir -p go.build/tmp
for folder in `ls ./app/`
do 
        if [ ! -e "$(pwd)/src/$folder" ];
        then
                mkdir -p "$(pwd)/src/"
                ln -s "$(pwd)/app/$folder" "$(pwd)/src/$folder"
        fi
done

for folder in `ls ./lib/`
do
        if [ ! -e "$(pwd)/src/$folder" ];
        then
                ln -s "$(pwd)/lib/$folder" "$(pwd)/src/$folder"
        fi
done

mkdir -p 3rdlib
for folder in `ls ./3rdlib/`
do
        if [ ! -e "$(pwd)/src/$folder" ];
        then
                ln -s "$(pwd)/3rdlib/$folder" "$(pwd)/src/$folder"
        fi
done
