#!/bin/bash
filepathPrefix='./'
mkdir -p ./make

maketestfile="./make/test"
echo '' > $maketestfile

makebenchfile="./make/bench"
echo '' > $makebenchfile

TestNode () {
        filepath=${filepathPrefix}$1

        if [[ $filepath == *"vendor"* ]] || [[ $filepath == *"pkg"* ]] || [[ $filepath == *"resource"* ]]; then
                return
        fi

        if [ "$(find $filepath -maxdepth 1 -name '*_test.go')" ]
        then
                MODULES="$MODULES $2"
                echo "$2:" >> $maketestfile
                echo -e "\tgo test $1" >> $maketestfile
                echo "" >> $maketestfile
        fi

        for node in `ls ${filepathPrefix}$1 | sort`
        do
                if [ -d "${filepathPrefix}$1/$node" ]
                then
                        TestNode $1/$node $2-$node
                fi
        done
}

BenchNode () {
        filepath=${filepathPrefix}$1

        if [[ $filepath == *"vendor"* ]] || [[ $filepath == *"pkg"* ]] || [[ $filepath == *"resource"* ]]; then
                return
        fi

        if [ "$(find $filepath -maxdepth 1 -name '*_test.go')" ]
        then
                MODULES="$MODULES $2"
                echo "$2:" >> $makebenchfile
                # echo -e "\tcd src/$1/ && go test -bench=. -benchmem -cpuprofile profile.out" >> $makebenchfile
                echo -e "\tcd src/$1/ && go test -bench=. -benchmem " >> $makebenchfile
                echo "" >> $makebenchfile
        fi

        for node in `ls ${filepathPrefix}$1 | sort`
        do
                if [ -d "${filepathPrefix}$1/$node" ]
                then
                        BenchNode $1/$node $2-$node
                fi
        done
}

MODULES=()
TestNode "./" "test"
echo "test:$MODULES" >> $maketestfile

MODULES=()
BenchNode "./" "bench"
echo "bench:$MODULES" >> $makebenchfile
