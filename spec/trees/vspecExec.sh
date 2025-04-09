#! /usr/bin/bash

COMMON_ARGS="-u ../units.yaml -q ../quantities.yaml"
VALIDATE="-e validate"
# Fix to prevent the overlay string to be truncated at the space char
OVERLAY=$3

function export_yaml {
    echo vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.yaml > exporterData/exporterLog.txt
    vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.yaml >> exporterData/exporterLog.txt
}

function export_json {
    echo vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.json > exporterData/exporterLog.txt
    vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.json >> exporterData/exporterLog.txt
}

function export_binary {
    echo vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.binary > exporterData/exporterLog.txt
    vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.binary >> exporterData/exporterLog.txt
}

function export_csv {
    echo vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.csv > exporterData/exporterLog.txt
    vspec export $1 $COMMON_ARGS $OVERLAY -s $2 $VALIDATE -o exporterData/cvis.csv >> exporterData/exporterLog.txt
}

if [ $1 = yaml ]; then
    export_yaml $1 $2
else
    if [ $1 = json ]; then
        export_json $1 $2
    else
        if [ $1 = binary ]; then
            export_binary $1 $2
        else
            if [ $1 = csv ]; then
                export_csv $1 $2
            else
                if [ $1 = all ]; then
                    export_yaml yaml $2
                    export_json json $2
                    export_binary binary $2
                    export_csv csv $2
                fi
            fi
        fi
    fi
fi
