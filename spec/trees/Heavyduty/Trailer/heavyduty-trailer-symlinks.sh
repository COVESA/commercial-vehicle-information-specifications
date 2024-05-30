#!/bin/bash

# This shell file is executed in the directory of the root vspecfile for the domain.
# The set of symlink commands below are for the trees/heavyduty/truck directory.

OBJECTPATH=../../../../objects/

echo "Updating trailer symlinks."
echo "In Vehicle directory..."

#rm Vehicle/CurrentLocation.yaml
rm Vehicle/GenericVehicle.yaml
rm Vehicle/VehicleIdentification.yaml
rm Vehicle/Version.yaml

#ln -s ${OBJECTPATH}Vehicle/CurrentLocation.yaml Vehicle/CurrentLocation.yaml
ln -s ${OBJECTPATH}Vehicle/GenericVehicle.yaml Vehicle/GenericVehicle.yaml
ln -s ${OBJECTPATH}Vehicle/VehicleIdentification.yaml Vehicle/VehicleIdentification.yaml
ln -s ${OBJECTPATH}Vehicle/Version.yaml Vehicle/Version.yaml

echo "In Body directory..."

rm Body/CargoBody.vspec
rm Body/FreezerBody.vspec

ln -s ${OBJECTPATH}Heavyduty/Trailer/CargoBody.vspec Body/CargoBody.vspec
ln -s ${OBJECTPATH}Heavyduty/Trailer/FreezerBody.vspec Body/FreezerBody.vspec

echo "In Chassis directory..."

rm Chassis/Accelerator.yaml
rm Chassis/Brake.yaml
rm Chassis/ParkingBrake.yaml
rm Chassis/Axle.vspec
rm Chassis/Wheel.yaml

ln -s ${OBJECTPATH}Chassis/Accelerator.yaml Chassis/Accelerator.yaml
ln -s ${OBJECTPATH}Chassis/Brake.yaml Chassis/Brake.yaml
ln -s ${OBJECTPATH}Chassis/ParkingBrake.yaml Chassis/ParkingBrake.yaml
ln -s ${OBJECTPATH}Chassis/Axle.vspec Chassis/Axle.vspec
ln -s ${OBJECTPATH}Chassis/Wheel.yaml Chassis/Wheel.yaml

echo "In Exterior directory..."

#rm Exterior/xxx.yaml

#ln -s ${OBJECTPATH}Exterior/xxx.yaml Exterior/xxx.yaml

echo "Trailer symlinks update done."

