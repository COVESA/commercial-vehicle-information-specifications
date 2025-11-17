![Status - Incubating](https://img.shields.io/static/v1?label=Status&message=Incubating&color=FEFF3A&style=for-the-badge)
# commercial-vehicle-information-specifications  

### Maintainers
Ulf Björkengren - Ford Motor Company

### Working group meetings
The CVIS working group meets every other Monday at 15.00 CET.
The schedule and link to meeting can be found [here](https://wiki.covesa.global/display/WIK4/COVESA+Common+Meeting+Schedule).<br>
Minutes of meetings is found [here](https://covesa.atlassian.net/wiki/spaces/WIK4/pages/39068025/CVIS+Meeting+Topics+and+Meeting+Notes).

## Overview
The Commercial Vehicle Information Specifications (CVIS) project is aiming at developing signal and service catalogues for commercial vehicles using the 
[HIM rule set](https://covesa.github.io/hierarchical_information_model/).

Using the HIM resource profile rule set enables the development of separate trees that each is tailored for the needs for specific vehicles types and models,
and it can be extended to separate vehicle individuals.

HIM also provides support for its implementation in an interface where it simplifies the server management of multiple trees.

Catalogues will be developed for different commercial vehicle types, e. g. trucks, trailers, buses, etc.

The HIM rule set for signals (the resource profile) is identical to the [VSS rule set](https://covesa.github.io/vehicle_signal_specification/rule_set/), 
and this project will for the signal parts try to align and reuse as much as possible with the VSS project.

A tool that makes it possible to modify vspec trees is being developed. It is currently on its third generation, starting of with using symlinks to assemble a vspec tree from smaller "vspec objects".

This was followed by the HIM configurator tool where an extension of the vspec format, called vspec2, was developed to provide features such as generation of ragged matrices (e. g. needed when axles have different number of wheels).

Building on a proposal from Volvo Group, a third generation of the tool has been developed, now named the vspecPreprocessor.
It supports the same features as the HIM configurator but now without requiring usage of the vspec2 format.

### Documentation
The CVIS documentation is [here](https://covesa.github.io/commercial-vehicle-information-specifications/).

## Installing the repo
The reo can be cloned by issung the command

git clone https://github.com/COVESA/commercial-vehicle-information-specifications.git

As the repo contains the VSS-tools as a submodule you might need to issue the command

git submodule update --init --recursive

to get it installed too.

After this installment the make file needs to be called to make sure that all packages used by the VSS-tools gets installed

sudo make install

This command involves access to files that requires super user priviledges, hence the sudo prefix.

## Contributors
CVIS is an open standard and we invite anybody to contribute. Currently CVIS contains - among others - significant  contributions from

 - [Ford Motor Company](https://www.ford.com/)
 - [Volvo Group](https://www.volvogroup.com/en/)
