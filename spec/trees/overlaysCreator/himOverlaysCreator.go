/**
* (C) 2025 Ford Motor Company
*
* All files and artifacts in the repository at https://github.com/covesa/commercial-vehicle-information-specifications
* are licensed under the provisions of the license provided by the LICENSE file in this repository.
*
**/

package main

import (
	"fmt"
	"reflect"
	"os"
	"strings"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/akamensky/argparse"
)

var pickFileName string
var outputLines []string
var ancestorList []string
var yamlTree map[string]interface{}


func saveNode(k string, v interface{}) {
	saveData(fmt.Sprintf("%s:", k))
	saveVal(v, 1)
	saveData("\n")
}

func saveVal(v interface{}, depth int) {
    typ := reflect.TypeOf(v).Kind()
    if typ == reflect.Int || typ == reflect.String || typ == reflect.Float64 {
    	value := fmt.Sprintf(" %v\n", v)
    	value= strings.Replace(value, ":", ";", -1)
    	saveData(value)
    } else if typ == reflect.Slice {
    	saveData("\n")
        saveSlice(v.([]interface{}), depth+1)
    } else if typ == reflect.Map {
    	saveData("\n")
        saveMap(v.(map[interface{}]interface{}), depth+1)
    } else {
        fmt.Printf("Type for %v is %s\n", v, typ)
    }
}

func saveMap(m map[interface{}]interface{}, depth int) {
    for k, v := range m {
    	saveData(fmt.Sprintf("%s%s:", strings.Repeat(" ", depth), k.(string)))
        saveVal(v, depth+1)
    }
}

func saveSlice(slc []interface{}, depth int) {
    for _, v := range slc {
    	saveData(fmt.Sprintf("%s", strings.Repeat(" ", depth)))
    	saveData("-")
        saveVal(v, 0)
    }
}

func saveData(line string) {
	outputLines = append(outputLines, line)
}

func readPathFile(vspecDir string) []string {
	pickFileName := "../" + vspecDir + pickFileName
	data, err := os.ReadFile(pickFileName)
	if err != nil {
		fmt.Printf("readConfigFile: Could not read %s\n", pickFileName)
		return nil
	}
	var pickFiles []string
	err = json.Unmarshal(data, &pickFiles)
	if err != nil {
		fmt.Printf("readPathFile:unmarshal %s, error=%s\n", string(data), err)
		return nil
	}	
	return pickFiles
}

func openFile(fName string) *os.File {
	outFp, err := os.Create(fName)
	if err != nil {
		fmt.Printf("Could not create %s\n", fName)
		return nil
	}
	return outFp
}

func writeLines(outputLines []string, outFp *os.File) {
	for i := 0; i < len(outputLines); i++ {
		outFp.Write([]byte(outputLines[i]))
	}
}

func isDot(c rune) bool {
	return c == '.'
}

func savedBefore(ancestorPath string) bool {
	for i := 0; i < len(ancestorList); i++ {
		if ancestorList[i] == ancestorPath {
			return true
		}
	}
	return false
}

func synthesizePath(nodeNames []string, index int) string {
	var path string
	for i := 0; i < index; i++ {
		path += nodeNames[i] + "."
	}
	return path[:len(path)-1]	
}

func getNodeData(path string) interface{} {
	for k, v := range yamlTree {
		if k == path {
			return v
		}
	}
	return nil
}

func checkAncestorNodes(path string) {
	nodeNames := strings.FieldsFunc(path, isDot)
	for i := 0; i < len(nodeNames); i++ {
		ancestorPath := synthesizePath(nodeNames, i+1)
		if !savedBefore(ancestorPath) {
			value := getNodeData(ancestorPath)
			saveNode(ancestorPath, value)
			ancestorList = append(ancestorList, ancestorPath)
		}
	}
}

func updateRootNodeName(outputLines []string, newRootName string) []string {
	dotIndex := strings.Index(outputLines[0], ".")
	if dotIndex == -1 {
		dotIndex = len(outputLines[0])
	}
	for i := 0; i < len(outputLines); i++ {
		if outputLines[i][0] != ' ' && outputLines[i][len(outputLines[i])-1] == ':' {
		if len(outputLines[i]) == dotIndex {
			outputLines[i] = newRootName + ":"
		} else {
			outputLines[i] = newRootName + outputLines[i][dotIndex-1:]
		}
		}
	}
	return outputLines
}

func main() {
	parser := argparse.NewParser("print", "HIM configurator")
	yamlFName := parser.String("y", "yamlfile", &argparse.Options{Required: false, Help: "File name of YAML tree", Default: "cvis.yaml"})
	pickFName := parser.String("p", "pickfile", &argparse.Options{Required: false, Help: "File name of pickpath array", Default: "himPickPaths.json"})
	outFName := parser.String("o", "outfile", &argparse.Options{Required: false, Help: "File name of overlay output", Default: "overlayPick.vspec"})
	vspecDir := parser.String("r", "rootdir", &argparse.Options{Required: false, Help: "Path to vspec root directory", Default: "Vehicle/VSS-core/"})
	nodeName := parser.String("n", "nodeName", &argparse.Options{Required: false, Help: "Root node name update"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	pickFileName = *pickFName

	data, err := os.ReadFile(*yamlFName)
	if err != nil {
		fmt.Printf("Could not open %s\n", *yamlFName)
		return
	}

	yamlTree = make(map[string]interface{})
	err = yaml.Unmarshal(data, &yamlTree)
	if err != nil {
		fmt.Printf("Could not parse YAML file %s\n", *yamlFName)
		os.Exit(-1)
	}

	pickFiles := readPathFile(*vspecDir)

	for i := 0; i < len(pickFiles); i++ {
		for k, v := range yamlTree {
			if len(k) >= len(pickFiles[i]) && pickFiles[i] == k[:len(pickFiles[i])] {
				checkAncestorNodes(k)
				saveNode(k, v)
			}
		}
	}
	if nodeName != nil && len(*nodeName) > 0 {
		outputLines = updateRootNodeName(outputLines, *nodeName)
	}
	outFp := openFile("../"+*vspecDir+*outFName)
	writeLines(outputLines,outFp)
	outFp.Close()
}
