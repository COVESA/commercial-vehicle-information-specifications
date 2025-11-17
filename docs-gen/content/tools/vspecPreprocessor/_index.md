---
title: vspecPreprocessor
weight: 20
chapter: true
---

# vspecPreprocessor

## Usage

As the vspecPreprocessor is written in Python it does not have o be compiled in a first step but can be run in interpreter mode.
To get help on the different CLI parameters it supports, open a terminal, navigate to the spec/trees directory and issue the command:\
$ python3 vspecPreprocessor.py -h

This will lead to the following being displayed:
```
usage: python3 vspecPreprocessor.py [-h] -i INPUTFILE [-s SCOPEFILE]
                                    [-o OUTPUTFILE] [-v VSPECFILE]
                                    [-f {yaml,json,binary}]

The VspecPreprocessor tool configures the vspec files by creating overlays to
be submitted with the vspec files to VSS-tools

options:
  -h, --help            show this help message and exit
  -i INPUTFILE, --inputfile INPUTFILE
                        JSON file containing configuration data
  -s SCOPEFILE, --scopefile SCOPEFILE
                        JSON file defining the scope of the different config
                        features
  -o OUTPUTFILE, --outputfile OUTPUTFILE
                        Overlay vspec file to be used with VSS-tools
  -v VSPECFILE, --vspecfile VSPECFILE
                        Root vspec file of the tree
  -f {yaml,json,binary}, --format {yaml,json,binary}
                        Exporter output format
```
The vspecPreprocessor requires two files as input, the INPUTFILE and the SCOPEFILE.
These are described in respective sections below.\
Its output is found in the OUTPUTFILE which is a VSS-tools overlays file containing the directives to configure the vspec domain tree
that it is combined with as input to the VSS-tools exporter.\
Currently the vspecPreprocessor does not issue the command to the VSS-tools exporter, it has to be done manually in a following step.
However, if a VSPECFILE pointing to the root vspec file of the domain tree the vspecPreprocessor will display the command that is to be used.
It will then by default select the YAML exporter, but with the -f CLI parameter another exporter can be selected.

## Initializing the Python virtual environment
The vspecPreprocessor creates input to the VSS-tools exporters for the generation of the transformed vspec files.
According to the [instructions](https://github.com/COVESA/vss-tools/blob/master/README.md)
for VSS-tools tool a Python virtual environment should be set up in which it then is run when exporting the vspec files to other formats.
The work flow is to the first time initialize this environment, and then every time before running a VSS-tools exporter
activate the virtual environment, and when done terminate the virtual environment.

If the virtual environment has not been set up before on the computer, it could e. g. already has been used in a VSS context to transform vspec files,
then the first time it needs to be created and configured.
To create it issue the following command:
```
$ python3 -m venv ~/.venv
```
Then it needs to be configured for VSS-tools, which shall be done with the environment activated:
```
$ source ~/.venv/bin/activate
(.venv)$ pip install --pre vss-tools
(.venv)$ deactivate
```
The above is only needed to be done once.
It might be necessary to install both python and pip if that is not already installed on the computer.

To then use the vspecPreprocessor it is sufficient with activating the virtual environment before using the vspecPreprocessor.
```
$ source ~/.venv/bin/activate
```
Then the VSS-tools exporter can be run with commands like:
```
(.venv)vspec export yaml -u Vehicle/Truck/units.yaml -q Vehicle/Truck/quantities.yaml -l Config/Truck/truck.vspec -s Vehicle/Truck/TruckSignalSpecification.vspec -o cvis.yaml
```
When done with VSS-tools the environment can be deactivated with the command.
```
(.venv)$ deactivate
```
If the VSS-tools exporter complains with error logs it may help to, one time, try the command:
```
(.venv)$ pip install -e .
```
A bash script file in the spec/trees directory, venv.sh, can be used instead of typing the commands manually in the terminal window.
If this file is issued as shown below it responds with displaying the two supported usages.
```
$ source venv.sh 
usage: source venv.sh startme|installme
```
As shown it can either be used to install venv, which as described above is only needed to be run the very first time.
The other alternative can thereafter be used to activate venv.

## vspecPreprocessor example input files
The directory spec/trees/Config contains example input files for some of the vspec trees that can be found under the spec/trees/Vehicle directory.

### Truck input files
The vehicleConfig-truck.json and the configScope.json are examples of INPUTFILE and SCOPEFILE, respectively.\
An example of their usage as input to the vspecPreprocessor is found in the vspecPreprocessor usage examples section below.

### Trailer input files
This directory contains two sets of input files.\
One will generate an overlay to be used together with a vspec tree representing a separate, standalone tree.\
The other will generate an overlay to be used together with a vspec tree representing a tree containing multiple trailers.\
To generate an overlay file for a separate trailer tree the vehicleConfig-separateTrailer.json and the configScope-separateTrailer.json are examples of INPUTFILE and SCOPEFILE, respectively.\
To generate an overlay file for a tree containing multiple trailers the vehicleConfig-multipleTrailers.json and the configScope-multipleTrailers.json are examples of INPUTFILE and SCOPEFILE, respectively.\
Examples of their respective usage as input to the vspecPreprocessor is found in the vspecPreprocessor usage examples section below.

### Car input files
The vehicleConfig.json and the configScope.json are examples of INPUTFILE and SCOPEFILE, respectively.\
An example of their usage as input to the vspecPreprocessor is found in the vspecPreprocessor usage examples section below.

## vspecPreprocessor usage examples
Below follows a few usage examples together with what the vspecPreprocessor displays in the terminal.
A reference to the respective VSPECFILE is used to show the command that can be issue to the VSS-tools exporter.
The default yaml output format is used.
* Create overlay for a truck.
```
$ python3 vspecPreprocessor.py -i Config/Truck/vehicleConfig-truck.json -o Config/Truck/truck.vspec -s Config/Truck/configScope.json -v Vehicle/Truck/TruckSignalSpecification.vspec

Overlay configuration saved to Config/Truck/truck.vspec

Exporter command: vspec export yaml -u Vehicle/Truck/units.yaml -q Vehicle/Truck/quantities.yaml -l Config/Truck/truck.vspec -s Vehicle/Truck/TruckSignalSpecification.vspec -o cvis.yaml
```
* Create overlay for a separate trailer.
```
$ python3 vspecPreprocessor.py -i Config/Trailer/vehicleConfig-separateTrailer.json -o Config/Trailer/separateTrailer.vspec -s Config/Trailer/configScope-separateTrailer.json -v Vehicle/Trailer/SeparateTrailerSignalSpecification.vspec

Overlay configuration saved to Config/Trailer/separateTrailer.vspec

Exporter command: vspec export yaml -u Vehicle/Trailer/units.yaml -q Vehicle/Trailer/quantities.yaml -l Config/Trailer/separateTrailer.vspec -s Vehicle/Trailer/SeparateTrailerSignalSpecification.vspec -o cvis.yaml
```
* Create overlay for multiple trailers in one tree.
```
$ python3 vspecPreprocessor.py -i Config/Trailer/vehicleConfig-multipleTrailers.json -o Config/Trailer/multipleTrailers.vspec -s Config/Trailer/configScope-multipleTrailers.json -v Vehicle/Trailer/MultipleTrailersSignalSpecification.vspec

Overlay configuration saved to Config/Trailer/multipleTrailers.vspec

Exporter command: vspec export yaml -u Vehicle/Trailer/units.yaml -q Vehicle/Trailer/quantities.yaml -l Config/Trailer/multipleTrailers.vspec -s Vehicle/Trailer/MultipleTrailersSignalSpecification.vspec -o cvis.yaml
```
* Create overlay for a passenger car.
```
$ python3 vspecPreprocessor.py -i Config/Car/vehicleConfig.json -o Config/Car/car.vspec -s Config/Car/configScope.json -v Vehicle/Car/VehicleSignalSpecification.vspec

Overlay configuration saved to Config/Car/car.vspec

Exporter command: vspec export yaml -u Vehicle/Car/units.yaml -q Vehicle/Car/quantities.yaml -l Config/Car/car.vspec -s Vehicle/Car/VehicleSignalSpecification.vspec -o cvis.yaml
```

## vspecPreprocessor features
The features supported by the vspecPreprocessor comes from what features the HIM configurator supported.
Further features are in the planning, such as configuration of default values for attribute signals,
definition of allowed values in a type definition tree (to make them reusable in a robust way).

### Instantiation of ragged matrices
A ragged matrix contains rows with varying number of columns in the different rows.
This is needed for the configuration of e. g. truck axle/wheel layouts where axles typically have different number of wheels.
Another example is seat position layouts in passenger cars,  busses, etc.\
Instantiation is supported for both 2-dimensional and 3-dimensional matrices.
The two examples above are both 2-dimensional, a 3-dimensional example is the instantiation of multiple trailers in the same tree,
which each contains an axle/wheel configuration.
An example of how the configuration of axle/wheels can look in the INPUTFILE is shown below.
```
            "instances": [
                {
                    "Axle": [
                        [
                            "Axle10",
                            "Axle12"
                        ],
                        [
                            [
                                "Pos3",
                                "Pos13"
                            ],
                            [
                                "Pos4",
                                "Pos12"
                            ]
                        ]
                    ]
                }
            ]
```
Here two axles are configured with two wheels each but on different positions.
A snippet of a 3-dimensional instantiation example can look in the INPUTFILE is shown below.
```
                    "Trailers": [
                        {
                            "Trailer1": {
                                "Axle": [
                                    [
                                        "Axle1",
                                        ...
```
The instance configuration data in the INPUTFILE must be complementd with data in the SCOPEFILE as the INPUTFILE data does not
say anything of where in the tree this should be applied, i. e. the missing parts to be able to synthesize a complete path.
For the 2-dimensional example above the SCOPEFILE data could look like below.
```
            "instance-scope": [
                {
                    "Axle": [
                        "Trailer.Chassis.Axle",
                        ".Wheel"
                    ]
                }
            ]
```
The key-name "Axle" is the logical connection between the two files. The two elements of the Axle-array are used to synthesize a path with the pattern:\
```
Trailer.Chassis.Axle.<Data from Axle-element 1>.Wheel.<Data from Axle-element 2, 3, etc>\
```
Similarly for the 3-dimensional example the SCOPEFILE data could look like below.
```
                    "Trailers": [
                        "Trailers",
                        ".Chassis.Axle",
                        ".Wheel"
                    ]
```
This is then used to synthesize the following path pattern:\
```
Trailers.<Key-name from Trailers objects>.Chassis.Axle.<Data from Axle-element 1 from Trailers objects>.Wheel.<Data from Axle-element 2, 3, etc from Trailers objects>\
```

### Configuration of different variants
A vspec tree may contain branches defining e. g. both electric and internal combustion engines.
If a tree is then to be configured for an electric engine only, the internal combustion engine branch needs to be removed,
and possibly an attribute needs to be set to 'ELECTRIC', etc.\
An example of how variant configuration is expressed in the INPUTFILE is shown below.
```
            "variants": [
                {
                    "TrailerType": "SEMI_TRAILER"
                }
            ]
```
The symbolic name of the variant to be configured is in this case TrailerType.
The available different variant selections for TrailerType are defined in the SCOPEFILE, in this example the SEMI_TRAILER was selected.
Below is shown a snippet of the list available TrailerType selections in the SCOPEFILE.
```
            "variant-scope": [
                {
                    "TrailerType": {
                        "SEMI_TRAILER": [
                            {
                                "Path": "Trailer.TrailerType",
                                "Directive": "  default: 'SEMI_TRAILER'",
                                "Description": "Configures the trailer tree as a SEMI_TRAILER"
                            }
                        ],
```
A variant selection consists of an array of JSON objects each containing the key-value pairs Path, Directive abd Description.\
Path points to the node in the tree that is to be modified.\
Directive holds the operation that is to be applied to the node.\
Description is only informative, it is not used by the tool.\
The variant selection can contain any number of these JSON objects, they can be applied to any nodes in the tree, also multiple times on the same node.\
Variant selections may require large variant selection lists if it is to  be applied to nodes that was created by an instantiation in the same configuration.
For that the instance variation feature is a better choice.

### Configuration of different variants on instantiated branches
This is on a high level the same feature as the variant configuration described above but as it requires different input to the vspecPreprocessor
it is represented as a separate feature.
An example is the configuration of 'axle features' that could be found on different axles of a truck or trailer
where one axle might support steering capability, another axle supports capability to be lifted from the road surface, an yet another to supports driving power.\
The vspecPreprocessor does not currently support the  instance-variant feature when the instantiation feature is used to create a 3-dimensional matrix.\
An example of how instance variant configuration is expressed in the INPUTFILE is shown below.
```
            "instance-variants": [
                {
                    "AxleFeature": ["LIFT", "Axle10"]
                },
                {
                    "AxleFeature": ["STEER", "Axle12"]
                }
            ]
```
The symbolic name of the instance variant to be configured is in this case AxleFeature.
The available different variant selections for AxleFeature are defined in the SCOPEFILE,
in this example LIFT was selected for Axle10, and STEER was selected for Axle12.
Below is shown a snippet of the list available AxleFeature selections in the SCOPEFILE.
```
            "instance-variant-scope": [
                {
                    "AxleFeature": {
                        "STEER": [
                            {
                                "Path": "Trailer.Chassis.Axle.X2.Steerable",
                                "Directive": "  delete: no",
                                "Description": "Configures the axle feature for a specific instance to STEER."
                            },
                            {
                                "Path": "Trailer.Chassis.Axle.X2.Liftable",
                                "Directive": "  delete: yes"
                            },
                            {
                                "Path": "Trailer.Chassis.Axle.X2.Driving",
                                "Directive": "  delete: yes"
                            }
                        ],
```
The content of an instance variant selection has the same format as for the variant selections, see above for the description of it.
