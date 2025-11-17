---
title: "CVIS vspec trees"
---

## Tree development
The [HIM data rule set](https://covesa.github.io/hierarchical_information_model/data_rule_set/) is used to define signals in a tree.
This syntax can be directly used as input to any of the exporter tools provided by [VSS-tools](https://github.com/COVESA/vss-tools).

Vspec trees for different vehicle types, and also for other domains such e. g. Driver, are being developed.
The root vspec file and other files needed as input to the VSS-tools exporters are for the different vehicle types found in respective directory under the spec/trees/Vehicle directory.
Currently tree data for the following domains are represented:
* [Truck tree](https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Truck)
* [Trailer tree](https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Trailer)
* [Car tree](https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Car)
* [Bus tree](https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Bus)
* [Driver tree](https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Driver)

### Truck tree
The truck tree is identical to the Car tree except for the Chassis.vspec file that is extended to include AxleFeature branches.

### Trailer tree
The trailer tree directory includes two trees, one that represents a separate, standalone trailer tree, and another that represents a tree containing multiple trailers.
The latter makes it possible to have one tree representing a train of one or more trailers.\
The difference is found in the root vspec file of respective tree, the other vspec files in the subdirectories are the same.

### Car tree
The car tree is identical to the [VSS](https://github.com/COVESA/vehicle_signal_specification/tree/master/spec) tree.
It is included here to show how the vspecPreprocessor can be used to e. g. generate a ragged matrix
for the configuration of seats with e. g. two seats in the first row and three in the second row.

### Bus tree
The bus tree is identical to the [VSS](https://github.com/COVESA/vehicle_signal_specification/tree/master/spec) tree.
The vspecPreprocessor can e. g. be used to generate a bus specific seat layout.

### Driver tree
The driver tree is yet to be populated. Currently is only contains empty branches to enable different data for EU and US legislation, respectively.
