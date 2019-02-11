#!/bin/bash
mkdir -p ./make

maketestfile="./make/test"
echo '' > $maketestfile

makebenchfile="./make/bench"
echo '' > $makebenchfile

TestNode () {
        filepathPrefix=$1
        filepath=${filepathPrefix}$2/
        if [ "$(find $filepath -maxdepth 1 -name '*_test.go')" ]
        then
                MODULES="$MODULES $3"
                echo "$3:" >> $maketestfile
                echo -e "\tgo test $2" >> $maketestfile
                echo "" >> $maketestfile
        fi

        for node in `ls ${filepathPrefix}$2/ | sort`
        do
                if [ -d "${filepathPrefix}$2/$node" ]
                then
                        TestNode ${filepathPrefix} $2/$node $3-$node
                fi
        done
}

BenchNode () {
        filepathPrefix=$1
        filepath=${filepathPrefix}$2/
        if [ "$(find $filepath -maxdepth 1 -name '*_test.go')" ]
        then
                MODULES="$MODULES $3"
                echo "$3:" >> $makebenchfile
                # echo -e "\tcd src/$2/ && go test -bench=. -benchmem -cpuprofile profile.out" >> $makebenchfile
                echo -e "\tcd src/$2/ && go test -bench=. -benchmem " >> $makebenchfile
                echo "" >> $makebenchfile
        fi

        for node in `ls ${filepathPrefix}$2/ | sort`
        do
                if [ -d "${filepathPrefix}$2/$node" ]
                then
                        BenchNode ${filepathPrefix} $2/$node $3-$node
                fi
        done
}

MODULES=()
TestNode "./lib/" "soloos" "test-soloos"
TestNode "./src/" "libsdfs" "test-libsdfs"
echo "test:$MODULES" >> $maketestfile

MODULES=()
BenchNode "./lib/" "soloos" "bench-soloos"
BenchNode "./src/" "libsdfs" "bench-libsdfs"
echo "bench:$MODULES" >> $makebenchfile
