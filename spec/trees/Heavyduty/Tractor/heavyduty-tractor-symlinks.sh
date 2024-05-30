#!/bin/bash

# This shell file is executed in the directory of the root vspecfile for the domain.
# The set of symlink commands below are for the trees/heavyduty/truck directory.

OBJECTPATH=../../../../objects/

echo "Updating Tractor symlinks."
echo "In Vehicle directory..."

rm Vehicle/CurrentLocation.yaml
rm Vehicle/GenericVehicle.yaml
rm Vehicle/LowVoltageBattery.yaml
rm Vehicle/Trailer.yaml
rm Vehicle/VehicleIdentification.yaml
rm Vehicle/Version.yaml

ln -s ${OBJECTPATH}Vehicle/CurrentLocation.yaml Vehicle/CurrentLocation.yaml
ln -s ${OBJECTPATH}Vehicle/GenericVehicle.yaml Vehicle/GenericVehicle.yaml
ln -s ${OBJECTPATH}Vehicle/LowVoltageBattery.yaml Vehicle/LowVoltageBattery.yaml
ln -s ${OBJECTPATH}Vehicle/Trailer.yaml Vehicle/Trailer.yaml
ln -s ${OBJECTPATH}Vehicle/VehicleIdentification.yaml Vehicle/VehicleIdentification.yaml
ln -s ${OBJECTPATH}Vehicle/Version.yaml Vehicle/Version.yaml

echo "In Powertrain directory..."

rm Powertrain/BatteryConditioning.yaml
rm Powertrain/CombustionEngine.yaml
rm Powertrain/ElectricMotor.yaml
rm Powertrain/EnergyManagement.yaml
rm Powertrain/EngineCoolant.yaml
rm Powertrain/FuelSystem.yaml
rm Powertrain/TractionBattery.yaml
rm Powertrain/Transmission.yaml

ln -s ${OBJECTPATH}Powertrain/BatteryConditioning.yaml Powertrain/BatteryConditioning.yaml
ln -s ${OBJECTPATH}Powertrain/CombustionEngine.yaml Powertrain/CombustionEngine.yaml
ln -s ${OBJECTPATH}Powertrain/ElectricMotor.yaml Powertrain/ElectricMotor.yaml
ln -s ${OBJECTPATH}Powertrain/EnergyManagement.yaml Powertrain/EnergyManagement.yaml
ln -s ${OBJECTPATH}Powertrain/EngineCoolant.yaml Powertrain/EngineCoolant.yaml
ln -s ${OBJECTPATH}Powertrain/FuelSystem.yaml Powertrain/FuelSystem.yaml
ln -s ${OBJECTPATH}Powertrain/TractionBattery.yaml Powertrain/TractionBattery.yaml
ln -s ${OBJECTPATH}Powertrain/Transmission.yaml Powertrain/Transmission.yaml

echo "In Chassis directory..."

rm Chassis/Accelerator.yaml
rm Chassis/Brake.yaml
rm Chassis/ParkingBrake.yaml
rm Chassis/SteeringWheel.yaml
rm Chassis/Axle.vspec
rm Chassis/Wheel.yaml

ln -s ${OBJECTPATH}Chassis/Accelerator.yaml Chassis/Accelerator.yaml
ln -s ${OBJECTPATH}Chassis/Brake.yaml Chassis/Brake.yaml
ln -s ${OBJECTPATH}Chassis/ParkingBrake.yaml Chassis/ParkingBrake.yaml
ln -s ${OBJECTPATH}Chassis/SteeringWheel.yaml Chassis/SteeringWheel.yaml
ln -s ${OBJECTPATH}Chassis/Axle.vspec Chassis/Axle.vspec
ln -s ${OBJECTPATH}Chassis/Wheel.yaml Chassis/Wheel.yaml

echo "In include directory..."

rm include/ItemHeatingCooling.yaml
rm include/LockableMovableItem.vspec
rm include/MovableItem.yaml
rm include/PowerOptimize.yaml

ln -s ${OBJECTPATH}include/ItemHeatingCooling.yaml include/ItemHeatingCooling.yaml
ln -s ${OBJECTPATH}include/LockableMovableItem.vspec include/LockableMovableItem.vspec
ln -s ${OBJECTPATH}include/MovableItem.yaml include/MovableItem.yaml
ln -s ${OBJECTPATH}include/PowerOptimize.yaml include/PowerOptimize.yaml

echo "Trator symlinks update done."

