---
title: Driver tree
weight: 20
chapter: true
---

# The driver tree
The vspec driver tree is used to represent the truck driver data.
The source for this driver data is the truck dirver recording unit (in europe often named tachograph, in the US xxx).
As the driver recording unit may hold data not only for the current driver but also for previous drivers,
there might be several instances of this driver tree, one for each driver that the driver recording unit has a record on.
The data from these different instances should be accessible via the interface that is exposed by vehicle server managing the access to the driver recording device.

