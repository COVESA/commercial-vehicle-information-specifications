# Setup of the HIM configurator in Windows, alternative 2
Steps to build and run HimConfigurator in Windows.

Pre-requisites in Windows:
* Latest Python version to be installed 
* Golang to be installed
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

  go.mod and go.sum files should be available to build a .go file. These files can be downloaded using below commands.

   go mod init cvis

   go mod tidy

Build HimConfigurator code with below command

    g. go build himConfigurator.go

Run the generated himConfigurator executable to generate vspec file based on required himConfiguration.
This also exports the vspec file to yaml file considering vspecExec.ps1 file available at spec/trees folder by taking into account of necessary overlay files.

    h. .\himConfigurator.exe  -m yaml

By default himConfig-truck.json will be selected.

But if we need to use himConfig-car, himConfigurator can be run with below command

    i. .\himConfigurator.exe  -m yaml -c himConfig-car.json

3. To generate yaml with vspec in VSS-Core tree. 

Run the command below from within commercial-vehicle-information-specifications\spec\trees\Vehicle\VSS-core directory 

    j. vspec export yaml  --vspec VssCoreSignalSpecification.vspec  --output cvis_car.yaml -u ../../../units.yaml -q ../../../quantities.yaml
         
4. To run himConfigurator in other trees  for example Vehicle\Truck  tree, steps in Win-setup\Setup-alt-1 can be followed as it is. 

In short,

Create the symlinks by running cv-truck-symlinks.ps1 script from the commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory
cv-truck-symlinks.ps1can be found in commercial-vehicle-information-specifications\ Win-setup\Setup-alt-1 directory

 HIM Configurator needs Datatypes.yaml, create it manually by running below command in the spec/objects/Datatype directory

    k. vspec export yaml -s ./Datatype.vspec -o Datatypes.yaml

Copy the DataTypes.yaml to the commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory

Build HimConfigurator code with below command

    l. go build himConfigurator.go
    
 Run the generated himConfigurator executable to generate vspec file based on required  himConfiguration.
 This also exports the vspec file to yaml file considering vspecExec.ps1 file available at spec/trees folder by taking into account of necessary overlay files.
 
Provide the tree folder in -r option

    m. .\himConfigurator.exe -m yaml -r  .\Vehicle\Truck\  -c himConfiguration.json


5. To generate yaml with vspec in Vehicle\Truck tree.
 
Run below command from within commercial-vehicle-information-specifications\spec\trees\Vehicle\Truck directory 

    n. vspec export yaml  --vspec TruckSignalSpecification.vspec  --output cvis_car.yaml -u ../../../units.yaml -q ../../../quantities.yaml
