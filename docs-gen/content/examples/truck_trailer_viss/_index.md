---
title: Truck and trailer trees used in VISS access scenario
weight: 20
chapter: true
---

## Creation of truck and trailer trees that are used in a client-server scenario using VISS as the interface
The following describes how a demo can be set up in which a client can read and write signals from a truck and a trailer that is connected to it using the VISS interface.\
The two repos that are needed for this are the [CVIS](https://github.com/COVESA/commercial-vehicle-information-specifications)
and the [VISSR](https://github.com/COVESA/vissr) repos, so start with cloning these if not done earlier.\
$ git clone https://github.com/COVESA/commercial-vehicle-information-specifications.git\
$ git clone https://github.com/COVESA/vissr.git\
The database that VISSR is configured to use must be installed on the computer. Redis is the default configuration.\
Please see the docmentation of respective database for how to install if that is needed.

Next step is to use the HIM configurator to create the truck and trailer trees.
The HIM configurator must then first be built which the following command ,issued in the CVIS/spec/trees directory.\
$ go build -o himConfigurator\
For the above to succeed the Golang build system must have been installed, which is described [here](https://go.dev/doc/install).

The HIM configurator calls the VSS-tools exporter after its initial preprocessing of vspec2 files.
To set up the environment required by the VSS-tools exporter the Python virtual environment, venv, must be activated.
This is done by the following command.\
$ source venv.sh startme\
Venv can when not needed anymore be deactivated by the command:\
(.venv)$ deactivate

The existing vspec trees that are found in the CVIS/spec/trees/Vehicle/VSS-core2 and CVIS/spec/trees/Vehicle/Trailer1 directories will be used.
These trees are not yet complete signal reprensentations for trucks and trailers but as it is what we have available at this point let us use them.
The HIM configurator tool will generate binary format representations of the two trees by invoking it as shown below.
If a human readable YAML version is desired, remove the '-m binary' from the commands or replace it with '-m yaml'.\
To generate the trailer tree:\
$ ./himConfigurator -c himConfiguration.json -r Vehicle/Trailer1/ -m binary\
The binary tree is then found in CVIS/spec/trees/exporterData with the file name cvis.binary.
Rename this file to Trailer1.binary.\
To generate the truck tree:\
$ ./himConfigurator -c himConfig-truck.json -r Vehicle/VSS-core2/ -m binary\
The binary tree is then found in CVIS/spec/trees/exporterData with the file name cvis.binary.
Rename this file to Truck.binary.

The binary format is the format that the VISSR server expects the trees to have.
The directory VISSR/server/vissv2server/forest is the local repository for trees so copy the two files Truck.binary and Trailer1.binary there.
$ cp  <your-local-path>/CVIS/spec/trees/exporterData/Truck.binary <your-local-path>/VISSR/server/vissv2server/forest/Truck.binary\
$ cp  <your-local-path>/CVIS/spec/trees/exporterData/Trailer1.binary <your-local-path>/VISSR/server/vissv2server/forest/Trailer1.binary\
<your-local-path> is to be replaced with the path on your device to the VISSR and CVIS root directories.

If you for some reason failed in creating the binary trees using the HIM configurator there are copies available in the CVIS/demodata directory.

For the VISSR server to register the binary trees at startup the file VISSR/server/vissv2server/viss.him must be modified to contain the information for
these trees. The file CVIS/demodata/viss.him is configured or this, so it can be copied to the VISSR/server/vissv2server directory.
If you do not want the viss.him file that is in this directory to be overwritten it needs to be renamed before the copy operation.

To enable a demo where a client can access dynamically changing signal values the feeder needs to be configured for this.
The feder implementation in vissr/feeder/feeder-template/feederv3 can be started with a CLI parameter "-i truck-trailer-sim.json"
where the file truck-trailer-sim.json contains the simulated signal data.
More about configuring the simulation can be found [here](https://covesa.github.io/vissr/feeder/#simulated-vehicle-data-sources).
A minimal simulation file truck-trailer-sim.json is available in the CVIS/demodata directory.
This file can easily be extended with more signals and more samples per signal.
Copy that file to vissr/feeder/feeder-template/feederv3.

To start the VISSR server and feeder copy the bash shell file CVIS/demodata/runtrucktrailer.sh to the VISSR directory.
The command below will then start them using the mentioned simulator input to the feeder.\
$ ./runtrucktrailer.sh startme

What is needed now is a client that can be used to issue requests for truck or trailer signals to the server.
Teh VISSR repo has a number of clients in the VISSR/client/client-1.0 directories that can be used,
either as is or as a code template to be modified.\
The VISSR/client/client-1.0/Javascript/wsclient.html that uses the Websocket protocol can be used directly, by clicking on the file it will open in the browser.
Start to enter the IP address of the server into the uppermost field and click on the Server IP button.
If the server has the same IP address as this client then "localhost" can be written into the field instead.\
Now request messages can be copied into the field below followed by clicking on the Send button and the responses from the server will be shown below in the browser tab.
Below is shown a few requests that can be used.
```
{"action":"get","path":"Vehicle.TraveledDistance","requestId":"232"}
{"action":"get","path":"Vehicle.CurrentLocation","filter":{"variant":"paths","parameter":["Longitude","Latitude"]},"requestId":"237"}
{"action":"subscribe","path":"Trailer1.Chassis.Axle.Row1.Wheel.Pos7.Tire.Temperature","filter":{"variant":"timebased","parameter":{"period":"1000"}},"requestId":"246"}
{"action":"unsubscribe","subscriptionId":"1","requestId":"240"}
```
Other websocket request examples can be found in VISSR/client/client-1.0/Javascript/appclient_commands.txt that can be modified to signal paths used in this demo.

To simulate that one more trailer is towed by the truck the steps below can be followed.\
Issue a command to stop the server and feeder:\
$ ./runtrucktrailer.sh stopme

### Add one more trailer tree.
This can be done by first creating another trailer binary tree using the HIM configurator,
but a shortcut is possible by using the same binary tree as for the first trailer.
To do this add the following rows to the viss.him file:
```
HIM.Trailer2:
  type: direct
  domain: Automotive.Trailer.Data
  version: 0.1.0
  local: forest/Trailer.binary
  public: https://github.com/COVESA/commercial-vehicle-information-specifications/tree/main/spec/trees/Vehicle/Trailer1
  description: A tree for the second trailer in the train.
```
The snippet above can be inserted at the top or at the bottom, or in between any two tree declarations.\
To get back any values from get requests to Trailer2 the truck-trailer-sim.json file needs to be updated with that.\
Then restart the server and feeder:\
$ ./runtrucktrailer.sh startme\
Signal data from the new trailer can now be accessed using requests with paths having Trailer2 as the root node name (the first segment name in the path).

The VISSR server supports 'forest inquiry' requests as shown below.
This is not part of the VISSv3.0 spec but may become part of a next version.
```
{"action":"get", "path":"HIM", "filter":{"variant":"metadata", "parameter":"0"}, "requestId":"1957"}
```
The response will contain the metadata from the viss.him file in JSON format,
except the 'local' property that is excluded as it is server internal data.

The browser tab rendered by wsclient.html may look something like below, depending on which requests that have been issued to the server.
![VISSR tech stack](/commercial-vehicle-information-specifications/images/ws-client-screenshot.png)
