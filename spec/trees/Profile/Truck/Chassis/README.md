# Axle features
The axle features defined in the AxleFeature directory are typically related to trucks and not directly to passenger cars.
In order to not 'bloat' the passenger car centric VSS vspec tree a separate truck vspec tree is created.
It is identical to the VSS tree except for the Chassis vspec file that is extended with include directives for the axle features in the AleFeature directory (which is not found in the VSS vspec directory structure).

If the Axlefeature is configured to NONE the Truck vspec tree can be used to generate an output tree that is identical to the standard passenger car (VSS) tree with the only difernce being the root node name Truck instead of Vehicle. This can obviously be fixed by updating the root node name in the root vspec file.

The Truck vspec tree can be used to generate Car configurations by selecting the AxleFeature configuration NONE which removes 
