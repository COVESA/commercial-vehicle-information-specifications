---
title: "Commercial Vehicle Information Specifications"
---
# Commercial Vehicle Information Specifications

## Overview
The HIM rule set for signals is used in this project, and since it is identical to the VSS rule set this framework inherits "patterns" from VSS

* The vspec file format.
* The usage of VSS-tools to transform from the vspec format to other exporter formats.
* The CVIS signal trees are defined in two directory structures -

objects directory structure: This is where the common "tree objects" are stored that may be shared between multiple trees.
trees directory structure: Thisis where the unique trees for different domains are stored.

The directory structure for a single tree follows the VSS pattern with "#include" links in the vspec files that logically links to other files of the tree. However, to link to a file in the common objects structure the corresponding file in the trees structure is realized as a symbolic link file. This means that when the content of the file is accessed the underlying file system follows the symbolic link to the file in the objects structure for the actual content of the file. This is transparent to the entity accessing the file, so e. g. the exporter tools from VSS-tools will when used for a transformation of a specific tree access file content from the common objects files tranparently.

## HIM configurator

The framework also contains a new tool, the HIM configurator. In its current version it provides support for two types of tree configuration:

* Variation point configuration: If the tree defined by the vspec files contains data structures that are not typically used together in a specific deployment of the tree, then these can be tagged as a variability point, and the HIM configurator can be used to pick the desired structure(s).
* Instance configuration: This is an extension of the existing VSS-tools instance support that provides the possibility for a two dimensional instantiation to have unique "column instantiation" for each "row instantiation".

If this new tool is found to be useful it is planned to add support for "default configuration" later.
