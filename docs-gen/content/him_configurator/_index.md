---
title: HIM configurator
weight: 20
chapter: true
---

# HIM configurator

## Build instructions

To build the HIM configurator it is necessary to first instal the Golang build system. Searching on "install golang" will lead to many matches of which one is:

https://go.dev/doc/install

To build the HIM configurator, open a terminal and go to the spec/trees directory where the himConfigurator.go is stored, and issue the command

$ go build -o himConfigurator

## Initializing the Python virtual environment
The HIM configurator uses the VSS-tools exporters for the generation of the transformed vspec files.
According to the [instructions](https://github.com/COVESA/vss-tools/blob/master/README.md)
for this tool a Python virtual environment should be set up in which it then is run when exporting the vspec files to other formats.
The work flow is to the first time initialize this environment, and then every time before running the HIM configurator
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

To then use the HIM configurator it is sufficient with activating the virtual environment before using the HIM configurator.
```
$ source ~/.venv/bin/activate
```
Then the HIM configurator can be run in the spec/trees directory with commands like:
```
(.venv)$ ./himConfigurator -m binary -v Vehicle/Trailer/
```
Finally the environment shall be deactivated it after the use.
```
(.venv)$ deactivate
```
If the HIM configurator complains with error logs it may help to, one time, try the command:
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

## Using the HIM configurator

Starting the HIM configurator with the command:

(.venv)$ ./himConfigurator -h

will show the command line options that is possible to apply at startup.

```
usage: print [-h|--help] [-m|--makecommand (all|yaml|json|csv|binary)]
             [-c|--configfile "<value>"] [-r|--rootdir "<value>"] [-v|--vspec]
             [-p|--preprocess] [-n|--noEnumSubst]

             HIM configurator

Arguments:

  -h  --help         Print help information
  -m  --makecommand  Make command parameter must be either: all, yaml, json, csv, or
                     binary. Default: yaml
  -c  --configfile   configuration file name. Default: himConfig-truck.json
  -r  --rootdir      path to vspec root directory. Default: Vehicle/VSS-core/
  -v  --vspec        Saves the configured .vspec2 files with extension .vspec
  -p  --preprocess   Pre-process only. Do not run VSS-tools
  -n  --noEnumSubst  No substitution of enum links to Datatype tree with actual
                     datatypes
```
The -m command line option is used to set which VSS-tools exporter the configured tree should have.
The value 'all' leads to that all the exporters that are compatible with the HIM configurator are executed.
If not used the default is 'yaml'.

The -c command line option is used to select the HIM configurator file. This is useful when a tree has multiple configuration files.
Default is himConfig-truck.json.

The -r command line option is used to set the path to the tree that should be used as input for the configuration.
The path is relative to the tree directory, and should have the slash '/' character at the end of the path.
If not used the default is 'Vehicle/VSS-core/'.

If the -v command  is set the the HIM configurator does not delete the vspec files that it generates from the vspec2 files that are found in the tree structure.
If not set these vspec files are deleted after being used as input to the call of the VSS-tools exporter(s).
If not used the default is false, i. e. not to save the files.\
With the vspec files saved it becomes possible to run the VSS-tols exporters "manually" with the vspec root node as input.
This might be helpful if the development environment is not a Linux compatible environment,
or if there is a need to debug the VSS-tools processing, se below.

The -p command leads to that the HIM configurator does not issue a command to the VSS-tools exporters.
If the -v command also was set, then the vspec files generated by the HIM configurator is saved, see above,
and the VSS-tools exporters can be activated "manually" with the pre-processed tree as input.

The -n command line option is used if substitution is not desired of the external datatype references with the actual enum definitions
from the common Datatypes tree in the generated tree(s). However, VSS-tools currently does not accept the syntax using a reference.

## VSS-tools debug
'in the development of the trees it might happen that errors are introdued in the vspec2/vspec files.
VSS-tools have a compehensive error logging support, but this does not show up in the HIM configurator UI.
It is however possible to get this error logging support by manually issuing a make command in a terminal window
on the make file in the cvis root catalog, after first running the HIM configurator on the tree under development,
with the -v CLI parameter so that the vspec files are saved.
The make command shall then have the following general format:
```
make yaml VSPECROOT=./spec/trees/Vehicle/Truck/TruckSignalSpecification.vspec
```
where yaml can be replaced by any other supported exporter, and the path to the root vspec file could be to any tree in the CVIS "forest".
The VSS-tools error logging will then be shown in the terminal window.

## Creation of a vehicle variant specific signal tree
The HIM configurator enables a model where a "super tree" covering all variations that a vehicle may be equipped with,
e. g. propulsion technology like ICE, PHEV, EV, etc. can be defined in the vspec files,
then the HIM configurator can be used to create a "variant specific" vehicle signal tree from the variation point configuration in the himConfiguration.json file.


## Usage of a vehicle variant specific signal tree
The output from the HIM configurator is the selected parts of the "vspec super tree" in a file with one of the supported formats.
This file can then e. g. be used in a vehicle of the selected variant by a server that is managing the access to the signals.
The server can then use the tree to "vet" client request - checking that the signal is present, whether it is read-only or read-write,
that credentials are valid if access control is applied, etc.
An example of a server using it like this is the [VISS reference server](https://github.com/COVESA/vissr).

