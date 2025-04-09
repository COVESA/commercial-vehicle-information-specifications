# Setup of the HIM configurator in Windows

 On Windows, before running this code:

Prerequisities: Visual Studio Code, Golang, the latest Python version (mark "Add to PATH" during the installation. pip will be installed as well).

1. Install the official COVESA tools by running "pip install vss-tools". If you already have it installed to upgrade to the latest version run "pip install --upgrade vss-tools".
Check that it works by running "vspec --help".

2. Clone the CVIS repository from https://github.com/COVESA/commercial-vehicle-information-specifications

3. Create the symlinks by running cv-truck-symlinks.ps1 script from the commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory

4. HIM Configurator needs Datatypes.yaml (see in the code below), create it manually by running in the spec/objects/Datatype directory:

vspec export yaml -s ./Datatype.vspec -o Datatypes.yaml

5. Copy the DataTypes.yaml to the commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory

6. To run the himConfiguratorWindows.go:

Rename the himConfigurator.go to himConfigurator.bak to prevent the package name conflict.

Place the provided .vscode directory into the commercial-vehicle-information-specifications directory.

The .vscode directory holds the launch.json file with correct HIM Configurator options to be run in VS Code.

Place the provided himConfiguratorWindows.go, go.mod and go.sum files in the commercial-vehicle-information-specifications\spec\trees directory.

Open the himConfiguratorWindows.go in VS Code and press F5. It will run the HIM Configurator with options "-m yaml -v Vehicle/Truck/" from
the commercial-vehicle-information-specifications\spec\trees directory.

7. Now you have all the standard .vspec files generated in all the subdirectories. You can run the vspec command manually
to generate the TruckSignalSpecificatiom.yaml file.

Run: vspec export yaml -u ../../../units.yaml -q ../../../quantities.yaml -s .\TruckSignalSpecification.vspec -o .\TruckSignalSpecification.yaml

8. Use other export options if you want to generate TruckSignalSpecification in other file formats.

## File locations
The launch.json file shall be saved in the .vscode directory.
The location of the other files are mentioned in the instructions above.

