---
title: VSS-core tree
weight: 20
chapter: true
---

# The VSS-core tree
The VSS-core tree is a vehicle type agnostic tree that supports the generation of vehicle type specific trees by usage of the HIM configurator.
The HIM configurator is then instructed to read input from vehicle specific configuration files.
This concept is in theory applicable to any vehicle type that shares at least some basic concepts with other vehicle types.
Vehicle types that are likely to be able to use this concept are e. g.:
* Car
* Truck
* Bus
* Van
* Pick-up

It may also be useful for vehicle types like:
* Trailer
* Motorcycle

For vehicle types that only share a small amount of basic concepts with other vehicle types the model that is also present in the CVIS project,
a structure with common objects that can be shared may be a better fit.
