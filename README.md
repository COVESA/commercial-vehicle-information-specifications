![Status - Incubating](https://img.shields.io/static/v1?label=Status&message=Incubating&color=FEFF3A&style=for-the-badge)
# commercial-vehicle-information-specifications  

### Maintainers
Ulf Björkengren - Ford Motor Company

### Working group meetings
The CVIS working group meets every other Monday at 16.00 CET.
The schedule and link to meeting can be found [here](https://wiki.covesa.global/display/WIK4/COVESA+Common+Meeting+Schedule).<br>
Minutes of meetings is found [here](https://wiki.covesa.global/display/WIK4/CVIS+Meeting+Topics+and+Meeting+Notes).

## Overview
The Commercial Vehicle Information Specifications (CVIS) project is aiming at developing signal and service catalogues for commercial vehicles using the 
[HIM rule set](https://covesa.github.io/hierarchical_information_model/).

Using the HIM resource profile rule set enables the development of separate trees that each is tailored for the needs for specific vehicles types and models,
and it can be extended to separate vehicle individuals.

HIM also provides support for its implementation in an interface where it simplifies the server management of multiple trees.

Catalogues will be developed for different commercial vehicle types, e. g. trucks, trailers, buses, etc.

The HIM rule set for signals (the resource profile) is identical to the [VSS rule set](https://covesa.github.io/vehicle_signal_specification/rule_set/), 
and this project will for the signal parts try to align and reuse as much as possible with the VSS project.

One attempt to realize this ambition is to use a structure where the parts of the tree that is domain specific
is represented in "private" directories, while the parts that can be shared with other trees are stored in "common" directories.
The use of symbolic links allow the private tree to access the common parts as if they were part of the private tree.
This design is compatible with the [VSS-tools](https://github.com/COVESA/vss-tools) exporter tool,
so it can be used to transform data in the vspec file structure to the different other formats supported by that tool.

Another attempt being tried out is to create a common tree for "all" vehicle types and then use the HIM configurator to configure this
tree to become vehicle type specific.

The CVIS documentation is [here](https://covesa.github.io/commercial-vehicle-information-specifications/).

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
