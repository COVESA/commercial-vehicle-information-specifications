---
title: HIM configurator
weight: 20
chapter: true
---

# HIM configurator

## Build instructions

To build the HIM configurator it is necessary to first instal the Golang build system. Searching on "install golang" will lead to many matches of which one is:

https://go.dev/doc/install

To build the HIM configurator, open a trminal and go to the spec/trees directory where the himConfigurator.go is stored, and issue the command

$ go build -o himConfigurator

## Using the HIM configurator

Starting the HIM configurator with the command:

$ ./himConfigurator -h
will show the command line options that is possible to apply at startup.

```
$ ./himConfigurator -h
usage: print [-h|--help] [-m|--makecommand (all|yaml|json|csv|binary)]
             [-v|--vspecdir "<value>"] [-c|--saveconf] [-e|--enumSubstitute]

             HIM preprocessor

Arguments:

  -h  --help            Print help information
  -m  --makecommand     Make command parameter must be either: all, yaml, csv,
                        or binary. Default: all
  -v  --vspecdir        path to vspec root directory. Default:
                        Heavyduty/Tractor/
  -c  --saveconf        Saves the configured vspec file with extension .conf.
                        Default: false
  -e  --enumSubstitute  Substitute enum links to Datatype tree with actual
                        datatypes. Default: false
```
The -m command line option is used to set which VSS-tools exporter the configured tree should have.
The value 'all' leads to that all the exporters that are compatible with the HIM configurator are executed.
If not used the default is 'all'.

The -v command line option is used to set the path to the tree that should be used as input for the configuration.
The path is relative to the tree directory, and should have the slash '/' character at the end of the path.
If not used the default is 'Heavyduty/Tractor/'.

The -c command line option is used to save a copy of the vspec files that are changed during the HIM configurator session.
These files are saved with their original name post-fixed by '.conf'.
If not set these files are deleted after being used as input to the call of the VSS-tools exporter(s).
If not used the default is false, i. e. not to save the files.

The -e command line option is used to substitute the external datatype references with the actual enum definitions
from the common Datatypes tree in the tree(s) that the VSS-tools generate.

