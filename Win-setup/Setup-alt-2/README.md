# Setup of the HIM configurator in Windows, alternative 2
Steps to build and run HimConfigurator in Windows.

Pre-requisites in Windows:
* Latest Python version to be installed 
* Golang 
* Update system environment variables with PATH for Python and Golang.
* Git is available
* Run the commands from cmd window or VS-Code.

1. VSS tools setup

This setup must be made in python virtual environment.

Python virtual environment can be setup with below commands

    a. python -m venv .venv   - This creates .venv folder (if this command doesn’t work try with “pip install virtualenv” before this command) 

    b. cd .venv 

    c. .\Scripts\activate

    d. pip install --pre vss-tools

       vss-tools can be cloned using command below

    e.  git clone https://github.com/COVESA/vss-tools.git

Required repositories like VSS , CVIS can be cloned from github

Clone CVIS from github using 

    f. git clone https://github.com/COVESA/commercial-vehicle-information-specifications.git

2. To run HimConfigurator on VSS-Core tree

Copy go.mod and go.sum  from commercial-vehicle-information-specifications\Win-setup directory to commercial-vehicle-information-specifications\spec\trees directory 

Build HimConfigurator code with 

    g. go build himConfigurator.go

Run the generated himConfigurator executable to generate vspec file based on required himConfiguration.

    h. .\himConfigurator.exe -p  -m yaml

By default himConfig-truck.json will be selected.

Ex : If we need to use himConfig-Car

    i. .\himConfigurator.exe -p  -m yaml -c himConfig-car.json

3. To generate yaml with vspec in VSS-Core tree . Run the command from within commercial-vehicle-information-specifications\spec\trees\Vehicle\VSS-core directory 

    j. vspec export yaml  --vspec VssCoreSignalSpecification.vspec  --output cvis_car.yaml -u ../../../units.yaml -q ../../../quantities.yaml
         
4. To run himConfigurator in other trees  for example commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck  directory ,  HIM Configurator needs Datatypes.yaml), create it manually by running in the spec/objects/Datatype directory

    k. vspec export yaml -s ./Datatype.vspec -o Datatypes.yaml

Copy the DataTypes.yaml to the commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory

 Run the generated himConfigurator executable to generate vspec file based on required  himConfiguration.
 
Provide the tree folder in -r option

    l. .\himConfigurator.exe -p -m yaml -r  .\Vehicle\Truck\  -c himConfiguration.json
