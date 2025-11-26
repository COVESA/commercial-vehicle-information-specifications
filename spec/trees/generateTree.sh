#!/bin/bash

usage() {
	echo "usage: ./generateTree.sh truck|car" >&2
}

truckTree() {
	echo "Generating Truck tree with Profile/Truck/ and config Config/Truck/vehicleConfig-truck.json ..."
    python3 vspecPreprocessor.py \
        -i Config/Truck/vehicleConfig-truck.json \
        -o Config/Truck/truck.vspec \
        -s Config/Truck/configScope.json \
        -v Vehicle/Car/VehicleSignalSpecification.vspec

    vspec export yaml \
        -u Vehicle/Car/units.yaml \
        -q Vehicle/Car/quantities.yaml \
        -l Profile/Truck/Chassis/Chassis.vspec \
        -l Config/Truck/truck.vspec \
        -s Vehicle/Car/VehicleSignalSpecification.vspec \
        -o cvis-truck.yaml
}

carTree() {
	echo "Generating Car tree with config Config/Car/vehicleConfig.json ..."
    python3 vspecPreprocessor.py \
        -i Config/Car/vehicleConfig.json \
        -o Config/Car/Car.vspec \
        -s Config/Car/configScope.json \
        -v Vehicle/Car/VehicleSignalSpecification.vspec

    vspec export yaml \
        -u Vehicle/Car/units.yaml \
        -q Vehicle/Car/quantities.yaml \
        -l Config/Car/Car.vspec \
        -s Vehicle/Car/VehicleSignalSpecification.vspec \
        -o cvis-car.yaml
}

case "$1" in
	truck)
		truckTree;;
	car)
		carTree;;
	*)
		usage;;
esac
