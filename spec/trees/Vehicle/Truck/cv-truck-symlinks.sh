#!/bin/bash

# This shell file is executed in the directory of the root vspecfile for the domain.
# The set of symlink commands below are for the trees/cv/truck directory.

OBJECTPATH=../../../../objects/

echo "Updating Truck symlinks."
echo "In Vehicle directory..."

rm Vehicle/CurrentLocation.vspec
rm Vehicle/GenericVehicle.vspec
rm Vehicle/LowVoltageBattery.vspec
rm Vehicle/Trailer.vspec
rm Vehicle/VehicleIdentification.vspec
rm Vehicle/Version.vspec

ln -s ${OBJECTPATH}Vehicle/CurrentLocation.vspec Vehicle/CurrentLocation.vspec
ln -s ${OBJECTPATH}Vehicle/GenericVehicle.vspec Vehicle/GenericVehicle.vspec
ln -s ${OBJECTPATH}Vehicle/LowVoltageBattery.vspec Vehicle/LowVoltageBattery.vspec
ln -s ${OBJECTPATH}Vehicle/Trailer.vspec Vehicle/Trailer.vspec
ln -s ${OBJECTPATH}Vehicle/VehicleIdentification.vspec Vehicle/VehicleIdentification.vspec
ln -s ${OBJECTPATH}Vehicle/Version.vspec Vehicle/Version.vspec

echo "In Powertrain directory..."

rm Powertrain/BatteryConditioning.vspec
rm Powertrain/CombustionEngine.vspec
rm Powertrain/ElectricMotor.vspec
rm Powertrain/EnergyManagement.vspec
rm Powertrain/EngineCoolant.vspec
rm Powertrain/FuelSystem.vspec
rm Powertrain/TractionBattery.vspec
rm Powertrain/Transmission.vspec

ln -s ${OBJECTPATH}Powertrain/BatteryConditioning.vspec Powertrain/BatteryConditioning.vspec
ln -s ${OBJECTPATH}Powertrain/CombustionEngine.vspec Powertrain/CombustionEngine.vspec
ln -s ${OBJECTPATH}Powertrain/ElectricMotor.vspec Powertrain/ElectricMotor.vspec
ln -s ${OBJECTPATH}Powertrain/EnergyManagement.vspec Powertrain/EnergyManagement.vspec
ln -s ${OBJECTPATH}Powertrain/EngineCoolant.vspec Powertrain/EngineCoolant.vspec
ln -s ${OBJECTPATH}Powertrain/FuelSystem.vspec Powertrain/FuelSystem.vspec
ln -s ${OBJECTPATH}Powertrain/TractionBattery.vspec Powertrain/TractionBattery.vspec
ln -s ${OBJECTPATH}Powertrain/Transmission.vspec Powertrain/Transmission.vspec

echo "In Chassis directory..."

rm Chassis/Accelerator.vspec
rm Chassis/Brake.vspec
rm Chassis/ParkingBrake.vspec
rm Chassis/SteeringWheel.vspec
rm Chassis/Axle.vspec2
rm Chassis/Wheel.vspec

ln -s ${OBJECTPATH}Chassis/Accelerator.vspec Chassis/Accelerator.vspec
ln -s ${OBJECTPATH}Chassis/Brake.vspec Chassis/Brake.vspec
ln -s ${OBJECTPATH}Chassis/ParkingBrake.vspec Chassis/ParkingBrake.vspec
ln -s ${OBJECTPATH}Chassis/SteeringWheel.vspec Chassis/SteeringWheel.vspec
ln -s ${OBJECTPATH}Chassis/Axle.vspec2 Chassis/Axle.vspec2
ln -s ${OBJECTPATH}Chassis/Wheel.vspec Chassis/Wheel.vspec

echo "In include directory..."

rm include/ItemHeatingCooling.vspec
rm include/LockableMovableItem.vspec
rm include/MovableItem.vspec
rm include/PowerOptimize.vspec

ln -s ${OBJECTPATH}include/ItemHeatingCooling.vspec include/ItemHeatingCooling.vspec
ln -s ${OBJECTPATH}include/LockableMovableItem.vspec include/LockableMovableItem.vspec
ln -s ${OBJECTPATH}include/MovableItem.vspec include/MovableItem.vspec
ln -s ${OBJECTPATH}include/PowerOptimize.vspec include/PowerOptimize.vspec

echo "Truck symlinks update done."

