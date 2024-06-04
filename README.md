![Status - Incubating](https://img.shields.io/static/v1?label=Status&message=Incubating&color=FEFF3A&style=for-the-badge)
# commercial-vehicle-information-specifications  

### Maintainers
Ulf Björkengren - Ford Motor Company


## Overview
The Commercial Vehicle Information Specifications (CVIS) project is aiming at developing signal and service catalogues for commercial vehicles using the 
[HIM rule set](https://github.com/COVESA/hierarchical_information_model).

Using the HIM rule set enables the development of separate trees that each is tailored for the needs for specific vehicles,
and it also enables the definition of trees representing services in the form of procedures with input and output parameters.

HIM also provides support for its implementation in an interface where it simplifies the server management of multiple trees.

Catalogues will be developed for different commercial vehicle types, e. g. Heavy duty tractors and trailers, buses, etc.

The HIM rule set for signals is identical to the [VSS](https://github.com/COVESA/vehicle_signal_specification), 
and this project will for the signal parts try to align and reuse as much as possible with the VSS project.

One realization of this ambition is the structure of how different trees are defined, where the parts of the tree that is domain specific
is represented in "private" directories, while the parts that can be shared with other trees are stored in "common" directories.

The use of symbolic links allow the private tree to access the common parts as if they were part of the private tree.
This design is compatible with the [VSS-tools](https://github.com/COVESA/vss-tools) exporter tool,
so it can be used to transform data in the vspec file structure to the different other formats supported by that tool.

The documentation is [here](https://covesa.github.io/commercial-vehicle-information-specifications/) (under construction).

## HIM configurator tool

To build the HIM configurator you need to have the Golang build system installed on your computer.
Then you go to the directory where the himConfigurator.go file is located, and issues the command:

go build -o himConfigurator

To see the command line options it can be started with, issue the command:

./himConfigurator -h

Please see the [documentation](https://covesa.github.io/commercial-vehicle-information-specifications/) for more details.

## Installing the repo
The reo can be cloned by issung the command

git clone https://github.com/COVESA/commercial-vehicle-information-specifications.git

As the repo contains the VSS-tools as a submodule you might need to issue the command

git submodule update --init --recursive

to get it installed too.

After this installment the make file needs to be called to make sure that all packages used by the VSS-tools gets installed

sudo make install

This command involves access to files that requires super user priviledges, hence the sudo prefix.
