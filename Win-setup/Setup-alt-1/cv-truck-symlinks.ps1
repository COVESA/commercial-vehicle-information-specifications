# This PowerShell script is executed in the directory of the root vspecfile for the domain.

# The set of hard link commands below are for the trees/cv/truck directory.

 

$OBJECTPATH = "../../../objects/"

 

Write-Output "Updating Truck hard links."

Write-Output "In Vehicle directory..."

 

Remove-Item Vehicle/CurrentLocation.vspec -Force

Remove-Item Vehicle/GenericVehicle.vspec -Force

Remove-Item Vehicle/LowVoltageBattery.vspec -Force

Remove-Item Vehicle/Trailer.vspec -Force

Remove-Item Vehicle/VehicleIdentification.vspec -Force

Remove-Item Vehicle/Version.vspec -Force

 

New-Item -ItemType HardLink -Path Vehicle/CurrentLocation.vspec -Target "$OBJECTPATH/Vehicle/CurrentLocation.vspec"

New-Item -ItemType HardLink -Path Vehicle/GenericVehicle.vspec -Target "$OBJECTPATH/Vehicle/GenericVehicle.vspec"

New-Item -ItemType HardLink -Path Vehicle/LowVoltageBattery.vspec -Target "$OBJECTPATH/Vehicle/LowVoltageBattery.vspec"

New-Item -ItemType HardLink -Path Vehicle/Trailer.vspec -Target "$OBJECTPATH/Vehicle/Trailer.vspec"

New-Item -ItemType HardLink -Path Vehicle/VehicleIdentification.vspec -Target "$OBJECTPATH/Vehicle/VehicleIdentification.vspec"

New-Item -ItemType HardLink -Path Vehicle/Version.vspec -Target "$OBJECTPATH/Vehicle/Version.vspec"

 

Write-Output "In Powertrain directory..."

 

Remove-Item Powertrain/BatteryConditioning.vspec -Force

Remove-Item Powertrain/CombustionEngine.vspec -Force

Remove-Item Powertrain/ElectricMotor.vspec -Force

Remove-Item Powertrain/EnergyManagement.vspec -Force

Remove-Item Powertrain/EngineCoolant.vspec -Force

Remove-Item Powertrain/FuelSystem.vspec -Force

Remove-Item Powertrain/TractionBattery.vspec -Force

Remove-Item Powertrain/Transmission.vspec -Force

 

New-Item -ItemType HardLink -Path Powertrain/BatteryConditioning.vspec -Target "$OBJECTPATH/Powertrain/BatteryConditioning.vspec"

New-Item -ItemType HardLink -Path Powertrain/CombustionEngine.vspec -Target "$OBJECTPATH/Powertrain/CombustionEngine.vspec"

New-Item -ItemType HardLink -Path Powertrain/ElectricMotor.vspec -Target "$OBJECTPATH/Powertrain/ElectricMotor.vspec"

New-Item -ItemType HardLink -Path Powertrain/EnergyManagement.vspec -Target "$OBJECTPATH/Powertrain/EnergyManagement.vspec"

New-Item -ItemType HardLink -Path Powertrain/EngineCoolant.vspec -Target "$OBJECTPATH/Powertrain/EngineCoolant.vspec"

New-Item -ItemType HardLink -Path Powertrain/FuelSystem.vspec -Target "$OBJECTPATH/Powertrain/FuelSystem.vspec"

New-Item -ItemType HardLink -Path Powertrain/TractionBattery.vspec -Target "$OBJECTPATH/Powertrain/TractionBattery.vspec"

New-Item -ItemType HardLink -Path Powertrain/Transmission.vspec -Target "$OBJECTPATH/Powertrain/Transmission.vspec"

 

Write-Output "In Chassis directory..."

 

Remove-Item Chassis/Accelerator.vspec -Force

Remove-Item Chassis/Brake.vspec -Force

Remove-Item Chassis/ParkingBrake.vspec -Force

Remove-Item Chassis/SteeringWheel.vspec -Force

Remove-Item Chassis/Axle.vspec2 -Force

Remove-Item Chassis/Wheel.vspec -Force

 

New-Item -ItemType HardLink -Path Chassis/Accelerator.vspec -Target "$OBJECTPATH/Chassis/Accelerator.vspec"

New-Item -ItemType HardLink -Path Chassis/Brake.vspec -Target "$OBJECTPATH/Chassis/Brake.vspec"

New-Item -ItemType HardLink -Path Chassis/ParkingBrake.vspec -Target "$OBJECTPATH/Chassis/ParkingBrake.vspec"

New-Item -ItemType HardLink -Path Chassis/SteeringWheel.vspec -Target "$OBJECTPATH/Chassis/SteeringWheel.vspec"

New-Item -ItemType HardLink -Path Chassis/Axle.vspec2 -Target "$OBJECTPATH/Chassis/Axle.vspec2"

New-Item -ItemType HardLink -Path Chassis/Wheel.vspec -Target "$OBJECTPATH/Chassis/Wheel.vspec"

 

Write-Output "In include directory..."

 

Remove-Item include/ItemHeatingCooling.vspec -Force

Remove-Item include/LockableMovableItem.vspec -Force

Remove-Item include/MovableItem.vspec -Force

Remove-Item include/PowerOptimize.vspec -Force

 

New-Item -ItemType HardLink -Path include/ItemHeatingCooling.vspec -Target "$OBJECTPATH/include/ItemHeatingCooling.vspec"

New-Item -ItemType HardLink -Path include/LockableMovableItem.vspec -Target "$OBJECTPATH/include/LockableMovableItem.vspec"

New-Item -ItemType HardLink -Path include/MovableItem.vspec -Target "$OBJECTPATH/include/MovableItem.vspec"

New-Item -ItemType HardLink -Path include/PowerOptimize.vspec -Target "$OBJECTPATH/include/PowerOptimize.vspec"

 

Write-Output "Truck hard links update done."
