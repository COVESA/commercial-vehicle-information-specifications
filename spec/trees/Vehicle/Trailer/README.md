# Trailer tree
This directory contains two different vspec root files, the SeparateTrailerSignalSpecification.vspec and the MultipleTrailersSignalSpecification.vspec.

The former is used to generate a tee that represents one trailer.

The latter is used to generate a tree that contains multiple trailer representations, the complete tree is then logically representing a train of connected trailers being towed by a truck.

For the former alternative to logically represent a train of connected trailers being towed by a truck, assuming that separate trees have been generated for the trailers, the logical combination into a train can be done by referencing to them in a [HIM configuration file](https://covesa.github.io/hierarchical_information_model/configuration_rule_set/).
An example of this is described [here](https://covesa.github.io/commercial-vehicle-information-specifications/examples/truck_trailer_viss/).

It is also possible to connect these trees to the truck tree to create a single tree for both truck and trailers by adding an #include directive at some suitable place in any of the the Truck vspec files.
