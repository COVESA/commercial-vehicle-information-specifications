/**
* (C) 2024 Ford Motor Company
*
* All files and artifacts in the repository at https://github.com/covesa/commercial-vehicle-information-specifications
* are licensed under the provisions of the license provided by the LICENSE file in this repository.
*
**/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"encoding/json"
	"bufio"
	"io/fs"
	"path/filepath"
	"github.com/akamensky/argparse"
)

var  saveConf bool
var makeCmd string
var enumSubstitute bool

type VariationPoint struct {
	VariantName string
	VariabilityName []string
}

type Variability struct {
	VariabilityType string
	VariationPointList []VariationPoint
}

var variabilityList []Variability

type Variant struct {
	VariantType string
	VariantName string	
}

var variantList []Variant

type RowDef struct {
	RowName string
}

type RowColumnDef struct {
	Column []ColumnDef
}

type ColumnDef struct {
	ColumnName string
}

type Instance struct {
	InstanceName string
	Row []RowDef
	RowColumn []RowColumnDef
}

var instanceList []Instance

type PropertyData struct {
	Name string
	NodeType string
	Datatype string
	Allowed []string
	Min string
	Max string
	Unit string
}

var enumData []PropertyData

type StructData struct {
	Name string
	Property []PropertyData
}

var structData []StructData

func variantProcess(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", fileName, err)
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text string
	continueScan := true
	var variationFp *os.File
	var savedLines []string
	for continueScan {
		continueScan = scanner.Scan()
		text = scanner.Text()
		if strings.Contains(text, "VariationPoint:") {
			fmt.Printf("Variation point line found in file=%s:%s\n", fileName, text)
			commentIndex := strings.Index(text, "#")
			if commentIndex == -1 {
				fmt.Printf("Error no comment found in line=%s\n", text)
				continue
			}
			variationType := text[commentIndex+1:]
			var variationLines []string
			var saveLine string
			variationLines, saveLine, scanner = readVariations(scanner)
			variationFp = updateVariationFile(variationFp, fileName, savedLines, variationLines, variationType, variabilityList, variantList)
			savedLines = nil
			savedLines = append(savedLines,saveLine)
		} else {
			savedLines = append(savedLines,text)
		}
	}
	if variationFp != nil {
		copyRemainingLines(variationFp, savedLines)
		variationFp.Close()
	}
	file.Close()
	if variationFp != nil {   // rename files: fileName -> fileName + ".orig" and fileName + ".tmp" -> fileName
		err = os.Rename(fileName, fileName + ".orig")
		if err != nil {
			fmt.Printf("Renaming %s failed. Err=%s\n", fileName, err)
		} else {
			err = os.Rename(fileName+".tmp", fileName)
			if err != nil {
				fmt.Printf("Renaming %s failed. Err=%s\n", fileName+".tmp", err)
			}
		}
	}
	return err
}

func readVariations(scanner *bufio.Scanner) ([]string, string, *bufio.Scanner) {
	var variations []string
	var text string
	continueScan := true
	for continueScan {
		continueScan = scanner.Scan()
		text = scanner.Text()
		if strings.Contains(text, "-") && strings.Contains(text, "#include") {
			variations = append(variations, text)
		} else {
			continueScan = false
		}
	}
	return variations, text, scanner
}

func updateVariationFile(variationFp *os.File, fileName string, savedLines []string, variationLines []string, variationType string, variabilityList []Variability, variantList []Variant) *os.File {
	if variationFp == nil {
		tmpFileName := fileName+".tmp"
		var err error
		variationFp, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("Could not create %s\n", tmpFileName)
			return nil
		}
	}
	for i := 0; i < len(savedLines); i++ {
		variationFp.Write([]byte(savedLines[i] + "\n"))
	}
	addVariation(variationFp, variationLines, variationType, variabilityList, variantList)
	return variationFp
}

func addVariation(variationFp *os.File, variationLines []string, variationType string, variabilityList []Variability, variantList []Variant) {
	var selectedVariations []string
	for i := 0; i < len(variantList); i++ {
		if variationType == variantList[i].VariantType {
			for j := 0; j < len(variabilityList); j++ {
				if variationType == variabilityList[j].VariabilityType {
					for k := 0; k < len(variabilityList[j].VariationPointList); k++ {
						if variantList[i].VariantName == variabilityList[j].VariationPointList[k].VariantName {
							selectedVariations = variabilityList[j].VariationPointList[k].VariabilityName
						}
					}
				}
			}
			
		}
	}
		for i := 0; i < len(selectedVariations); i++ {
			for j := 0; j < len(variationLines); j++ {
				if strings.Contains(variationLines[j], "- " + selectedVariations[i]) {
					commentIndex := strings.Index(variationLines[j], "#")
					if commentIndex != -1 {
						variationFp.Write([]byte(variationLines[j][commentIndex:] + "\n"))
						fmt.Printf("Variant %s: Inserted:%s\n", selectedVariations[i], variationLines[j][commentIndex:])
					}
				}
			}
		}
}


func copyRemainingLines(variationFp *os.File, savedLines []string) {
	for i := 0; i < len(savedLines); i++ {
		variationFp.Write([]byte(savedLines[i] + "\n"))
	}
}

func instanceProcess(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", fileName, err)
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text string
	var nodeName string
	var thisNode []string
	var subTree []string
	var nextNode []string
	continueScan := true
	var savedLines []string
	isConfigDone := false
	for continueScan {
		continueScan = scanner.Scan()
		text = scanner.Text()
		getNodeName(text, &nodeName)
		isConfigInstance, instanceTag := checkConfigInstance(text, 0)
		if isConfigInstance {
			isConfigDone = true
			fmt.Printf("Instance line found in file=%s:%s\n", fileName, text)
			if instanceRows(instanceTag) == -1 {
				fmt.Printf("Instance configuration=%s not found. Line skipped.\n", instanceTag)
				continue
			}
			thisNode, subTree, nextNode, scanner = readSubtree(scanner, nodeName)
			instanceExpression := getInstanceExpression(thisNode, 1, instanceTag)
			if len(instanceExpression) > 0 { // finish this node, create row branch nodes with column instances added to it, followed by subtree, and next node
				for i := 0; i < len(thisNode); i++ {// finish this node
					if !strings.Contains(thisNode[i], "instances1:") {
						savedLines = append(savedLines,thisNode[i])
					}
				}
				for i := 0; i < instanceRows(instanceTag); i++ {  // create row branch nodes with column instances added to it
					savedLines = addInstanceBranch(savedLines, nodeName, i, instanceTag, instanceExpression)
					for j := 0; j < len(subTree); j++ {  // followed by subtree
						if strings.Contains(subTree[j], "#include") {
							savedLines = append(savedLines,subTree[j] + "." + getRowInstance(instanceTag,i))
						} else {
							savedLines = append(savedLines,subTree[j])
						}
					}
				}
				for i := 0; i < len(nextNode); i++ {  // and next node
					savedLines = append(savedLines,nextNode[i])
				}	
			} else {  // finish this node, create row branch nodes, followed by subtree with configured instance, and next node
				for i := 0; i < len(thisNode); i++ {// finish this node
					savedLines = append(savedLines,thisNode[i])
				}
				instExpTree := make([]string, 1)  // needed in calls to getInstanceExpression
				for i := 0; i < instanceRows(instanceTag); i++ {  // create row branch nodes
//fmt.Printf("instanceRow no=%d\n", i)
					instSubTree := expandSubTree(subTree, filepath.Dir(fileName)+"/", nodeName, i, instanceTag)
					savedLines = addInstanceBranch(savedLines, nodeName, i, instanceTag, "")
					for j := 0; j < len(instSubTree); j++ {  // followed by subtree with configured instance
						instExpTree[0] = instSubTree[j]
						instanceExp := getInstanceExpression(instExpTree, 1, instanceTag)
						if len(instanceExp) > 0 {
							savedLines = append(savedLines,"  " + createConfiguredInstance(instanceExp, i, instanceTag))
						}  else {
							savedLines = append(savedLines,instSubTree[j])
						}
					}
				}
				for i := 0; i < len(nextNode); i++ {  // and next node
					savedLines = append(savedLines,nextNode[i])
				}	
			}
		} else {
			savedLines = append(savedLines,text)
		}
	}
	file.Close()
	if isConfigDone {   // if fileName.orig does not exist: rename fileName to fileName + ".orig" and rewrite fileName with savedLines
//fmt.Printf("isConfigDone=true\n")
		if !fileExists(fileName + ".orig") {
			err = os.Rename(fileName, fileName + ".orig")
			if err != nil {
				fmt.Printf("Renaming %s failed. Err=%s\n", fileName, err)
			}
		} else {
			err = os.Remove(fileName)
			if err != nil {
				fmt.Printf("Deleting %s failed. Err=%s\n", fileName, err)
			}
		}
		if err == nil {
			var instanceFp *os.File
			instanceFp, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				fmt.Printf("Could not create %s\n", fileName)
			} else {
				for i := 0; i < len(savedLines); i++ {
//fmt.Printf("SavedLine: %s\n", savedLines[i])
					instanceFp.Write([]byte(savedLines[i] + "\n"))
				}
				instanceFp.Close()
			}
		}
	}
	return err
}

func expandSubTree(subTree []string, path string, rootNodeName string, index int, instanceTag string) []string {
	var expandedTree []string
	for i := 0; i < len(subTree); i++ {
		if strings.Contains(subTree[i], "#include") && strings.Contains(subTree[i], rootNodeName) {
//fmt.Printf("expandSubTree: Include to be expanded found for root node:%s=%s\n", rootNodeName, subTree[i])
			includeExpansion := readIncludefile(subTree[i], path, index, instanceTag)
			for j := 0; j < len(includeExpansion); j++ {
				expandedTree = append(expandedTree, includeExpansion[j])
			}
		} else {
			expandedTree = append(expandedTree, subTree[i])
		}
	} 
	return expandedTree
}

func readIncludefile(includeDirective string, path string, index int, instanceTag string) []string { 
// if config instance directive found, update it with config data, update node names and #include directives with rootnode data
	vspecFile, nodeNamePrefix := decodeIncludeDirective(includeDirective)
	file, err := os.Open(path + vspecFile)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", path + vspecFile, err)
		return nil
	}
//fmt.Printf("vspecfile: %s\n", path + vspecFile)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var line string
	continueScan := true
	var includeLines []string
	newNodeName := ""
	for continueScan {
		continueScan = scanner.Scan()
		line = scanner.Text()
		getNodeName(line, &newNodeName)
		if len(line) == 0 || line[0] != '#' || strings.Contains(line, "#include") {
			if len(newNodeName) > 0 {
				line = nodeNamePrefix + "." + getRowInstance(instanceTag,index) + "." + line
			} else if strings.Contains(line, "#include") {
//fmt.Printf("includeLine: %s\n", line)
				incFields := strings.Fields(line)
				line = incFields[0] + " " + incFields[1] + " " + nodeNamePrefix + "." + getRowInstance(instanceTag,index) + "." + incFields[2]
			}
			includeLines = append(includeLines, line)
		}
		newNodeName = ""
	}
	return includeLines
}

func decodeIncludeDirective(includeDirective string) (string, string) {  // e.g. #include Axle.vspec Chassis.Axle
	fields := strings.Fields(includeDirective)
	if len(fields) > 2{
		return fields[1], fields[2]
	}
	return fields[1], ""
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}

func getNodeName(line string, nodeName *string) {
	if len(line) > 0 && line[0] != ' ' && line[0] != '#' && line[len(line)-1] == ':' {
		*nodeName = line[:len(line)-1]
	}
}

func getInstanceExpression(tree []string, instLevel int, instanceTag string) string {
	for i:= 0; i < len(tree); i++ {
		isConfigInstance, instanceTag2 := checkConfigInstance(tree[i], instLevel)
		if isConfigInstance && instanceTag == instanceTag2 {
			return tree[i]
		}
	}
	return ""
}

func checkConfigInstance(line string, instLevel int) (bool, string) {  // configured instances is expressed on one line as instancesn: x #tag
	instKey := "instances" + strconv.Itoa(instLevel)
	expressionFields := strings.Fields(line)
	if len(expressionFields) == 3 && strings.Contains(expressionFields[0], instKey) && strings.Contains(expressionFields[2], "#") {
		sharpIndex := strings.Index(line, "#") + 1
		return true, line[sharpIndex:]
	}
	return false, ""
}

func addInstanceBranch(savedLines []string, nodeName string, index int, instanceTag string, instanceExpression string) []string {
	savedLines = append(savedLines, "")
	savedLines = append(savedLines, nodeName + "." + getRowInstance(instanceTag,index) + ":")
	savedLines = append(savedLines, "  type: branch")
	if len(instanceExpression) > 0 {
		savedLines = append(savedLines, "  " + createConfiguredInstance(instanceExpression, index, instanceTag))
	}
	savedLines = append(savedLines, "  description: " + nodeName  + "." + getRowInstance(instanceTag,index))
	return savedLines
}

func createConfiguredInstance(instanceExpression string, index int, instanceTag string) string {
	return "instances: " + getRowColumnInstance(instanceTag,index)
}

func getRowColumnInstance(instanceConfigName string,index int) string {
	for i := 0; i < len(instanceList); i++ {
		if instanceList[i].InstanceName == instanceConfigName {
			instanceExpr := `["`
			for j := 0; j < len(instanceList[i].RowColumn[index].Column); j++ {
				instanceExpr += instanceList[i].RowColumn[index].Column[j].ColumnName + `", "`
			}
			instanceExpr = instanceExpr[:len(instanceExpr)-3] + "]"
			return instanceExpr
		}
	}
	return ""
}

func getRowInstance(instanceConfigName string,index int) string {
	for i := 0; i < len(instanceList); i++ {
		if instanceList[i].InstanceName == instanceConfigName {
			return instanceList[i].Row[index].RowName
		}
	}
	return ""
}

func instanceRows(instanceConfigName string) int {
	for i := 0; i < len(instanceList); i++ {
		if instanceList[i].InstanceName == instanceConfigName {
			return len(instanceList[i].Row)
		}
	}
	return -1
}

func readSubtree(scanner *bufio.Scanner, rootNodeName string) ([]string, []string, []string, *bufio.Scanner) {
	var tree []string  // will contain parts of root node, subtree, and the following node
	var text string
	newNodeName := ""
	continueScan := true
	for continueScan {  // read lines until a new node that is not part of the subtree
		continueScan = scanner.Scan()
		text = scanner.Text()
		getNodeName(text, &newNodeName)
		if len(newNodeName) > 0 && !strings.Contains(newNodeName, rootNodeName) {
			continueScan = false
		}
		tree = append(tree, text)
	}
	//find boundary between root node and subtree
	splitIndex1 := 0
	newNodeName = ""
	for i := 0; i < len(tree); i++ {
//fmt.Printf("1stpass:line[%d]=%s\n", i, tree[i])
		getNodeName(tree[i], &newNodeName)
		if len(tree[i]) > 0 && (strings.Contains(tree[i], "#include") || len(newNodeName) > 0) {
			splitIndex1 = i
			break
		}
	}
	//find boundary between subtree and following node
	splitIndex2 := len(tree)-2
	for i := len(tree)-2; i > splitIndex1-1; i-- {
		if len(tree[i]) > 0 && ((strings.Contains(tree[i], "#include") && strings.Contains(tree[i], rootNodeName)) || (tree[i][0] != '#' && len(strings.TrimSpace(tree[i])) > 0)) {
//fmt.Printf("2ndpass:line[%d]=%s\n", i, tree[i])
			splitIndex2 = i + 1
			break
		}
	}
//fmt.Printf("readSubtree:splitIndex1=%d, splitIndex2=%d\n", splitIndex1, splitIndex2)
	return tree[:splitIndex1], tree[splitIndex1:splitIndex2], tree[splitIndex2:], scanner
}

func walkVariantPass(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
//		fmt.Printf("Enter dir=%s\n", s)
	} else {
		if filepath.Ext(s) == ".vspec" {
//			fmt.Printf("Vspec path=%s\n", s)
			err = variantProcess(s)
		}
	}
	return err
}

func walkInstancePass(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
//		fmt.Printf("Enter dir=%s\n", s)
	} else {
		if filepath.Ext(s) == ".vspec" {  // it is assumed that variantProcess() is run and created .orig file copies.
//			fmt.Printf("instanceProcess:Vspec path=%s\n", s)
			err = instanceProcess(s)
		}
	}
	return err
}

func walkEnumSubstitute(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
//		fmt.Printf("Enter dir=%s\n", s)
	} else {
		if filepath.Ext(s) == ".vspec" {  // it is assumed that variantProcess() and instanceProcess() re run and created .orig file copies.
//			fmt.Printf("enumProcess:Vspec path=%s\n", s)
			err = enumProcess(s)
		}
	}
	return err
}

func enumProcess(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", fileName, err)
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text string
	continueScan := true
	var savedLines []string
	isConfigDone := false
	for continueScan {
		continueScan = scanner.Scan()
		text = scanner.Text()
		if isDataTypedEnum(text) {
			enumLines := getExpandedEnumData(text)
			for i := 0; i<len(enumLines); i++ {
				savedLines = append(savedLines,enumLines[i])
			}
			isConfigDone = true
		} else {
			savedLines = append(savedLines,text)
		}
	}
	file.Close()
	if isConfigDone {   // if fileName.orig does not exist: rename fileName to fileName + ".orig" and rewrite fileName with savedLines
//fmt.Printf("isConfigDone=true\n")
		if !fileExists(fileName + ".orig") {
			err = os.Rename(fileName, fileName + ".orig")
			if err != nil {
				fmt.Printf("Renaming %s failed. Err=%s\n", fileName, err)
			}
		} else {
			err = os.Remove(fileName)
			if err != nil {
				fmt.Printf("Deleting %s failed. Err=%s\n", fileName, err)
			}
		}
		if err == nil {
			var instanceFp *os.File
			instanceFp, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				fmt.Printf("Could not create %s\n", fileName)
			} else {
				for i := 0; i < len(savedLines); i++ {
//fmt.Printf("SavedLine: %s\n", savedLines[i])
					instanceFp.Write([]byte(savedLines[i] + "\n"))
				}
				instanceFp.Close()
			}
		}
	}
	return err
}

func isDataTypedEnum(text string) bool {
	if getExtEnumRef(text) != "" {
			return true
	}
	return false
}

func getExtEnumRef(text string) string {
	dtIndex := strings.Index(text, "datatype:")
	if dtIndex != -1 && text[0] != '#' {
		if strings.Contains(text[dtIndex+8+1:], ".") {
			return strings.TrimSpace(text[dtIndex+8+1:])
		}
	}
	return ""
}

func getExpandedEnumData(text string) []string {
	var expandedData []string
	enumRef := getExtEnumRef(text)
	for i := 0; i < len(enumData); i++ {
		if enumData[i].Name == enumRef {
fmt.Printf("%s expanded\n", enumData[i].Name)
			expandedData = append(expandedData, "  datatype: string")
			expandedData = append(expandedData, "  allowed:")
			for j := 0; j < len(enumData[i].Allowed); j++ {
				expandedData = append(expandedData, enumData[i].Allowed[j])
			}
		}
	}
	if len(expandedData) == 0 {
		expandedData = append(expandedData, text)
	}
	return expandedData
}

func walkPostmake(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
//		fmt.Printf("Enter dir=%s\n", s)
	} else {
		if filepath.Ext(s) == ".orig" {
			origIndex := strings.Index(s, ".orig")
			postProcess(s[:origIndex])
		}
	}
	return nil
}

func postProcess(fileName string) {  // rename or delete fileName, fileName+".orig" -> fileName
	var err error
	if saveConf {
		err = os.Rename(fileName, fileName+".conf")
	} else {
		err = os.Remove(fileName)
	}
	if err != nil {
		if saveConf {
			fmt.Printf("Renaming %s failed. Err=%s\n", fileName, err)
		} else {
			fmt.Printf("Deleting %s failed. Err=%s\n", fileName, err)
		}
	} else {
		err = os.Rename(fileName+".orig", fileName)
		if err != nil {
			fmt.Printf("Renaming %s failed. Err=%s\n", fileName+".orig", err)
		} else {
		}
	}
}

func decodeVariantConfigs(variantConfigs string) []Variant { //JSON object:{"var-type1":"var-name1", ., "var-typeN":"var-nameN"}
	var variantList []Variant
	var variant Variant
	var variantMap = make(map[string]interface{})
	err := json.Unmarshal([]byte(variantConfigs), &variantMap)
	if err != nil {
		fmt.Printf("decodeVariantConfigs():unmarshal error=%s\n", err)
		return nil
	}
	i := 0
	for k, v := range variantMap {
		switch vv := v.(type) {
		case string:
//			fmt.Println(vv, "is string")
			variantList = append(variantList, variant)
			variantList[i].VariantType = k
			variantList[i].VariantName = v.(string)
			i++
		default:
			fmt.Println(vv, "is of an unknown type")
		}
	}
	return variantList
}

func readVariabilityFile(variabilityFile string) []Variability {
	var variabilityList []Variability
	data, err := os.ReadFile(variabilityFile)
	if err != nil {
		fmt.Printf("Could not open %s\n", variabilityFile)
		return nil
	}
	var variabilityListMap = make(map[string]interface{})
	err = json.Unmarshal(data, &variabilityListMap)
	if err != nil {
		fmt.Printf("readVariabilityFile():unmarshal error=%s", err)
		return nil
	}
	variabilityList = unpackVarMapLevel1(variabilityListMap, variabilityList)
	return variabilityList
}

func unpackVarMapLevel1(variabilityMap map[string]interface{}, variabilityList []Variability) []Variability {
	i := 0
	for k, v := range variabilityMap {
		variabilityList = append(variabilityList, Variability{})
		switch vv := v.(type) {
		case []interface{}:
//			fmt.Println(vv, "is an array:, len=", len(vv))
			variabilityList[i].VariabilityType = k
			variabilityList[i].VariationPointList = make([]VariationPoint, len(vv))
			variabilityList = unpackVarMapLevel2(i, vv, variabilityList)
		case map[string]interface{}:
//			fmt.Println(vv, "is a map:")
			variabilityList[i].VariabilityType = k
			variabilityList[i].VariationPointList = make([]VariationPoint, 1)
			variabilityList[i].VariationPointList[0].VariantName = k
		default:
			fmt.Println(vv, "is of an unknown type")
		}
		i++
	}
	return variabilityList
}

func unpackVarMapLevel2(index1 int, variantDefMap []interface{}, variabilityList []Variability) []Variability {
	index2 := 0
	for _, v := range variantDefMap {
		switch vv := v.(type) {
		case map[string]interface{}:
//			fmt.Println(vv, "is a map:")
			variabilityList = unpackVarMapLevel3(index1, index2, vv, variabilityList)
		default:
			fmt.Println(vv, "is of an unknown type")
		}
		index2++
	}
	return variabilityList
}

func unpackVarMapLevel3(index1 int, index2 int, variantDefMap map[string]interface{}, variabilityList []Variability) []Variability {
	for k, v := range variantDefMap {
		variabilityList[index1].VariationPointList[index2].VariantName = k
		switch vv := v.(type) {
		case interface{}:
//			fmt.Println(vv, "is an interface")
			variabilityList = unpackVarMapLevel4(index1, index2, vv, variabilityList)
		default:
			fmt.Println(vv, "is of an unknown type")
		}
		index2++
	}
	return variabilityList
}

func unpackVarMapLevel4(index1 int, index2 int, variabilityNameArrayMap interface{}, variabilityList []Variability) []Variability {
		switch vv := variabilityNameArrayMap.(type) {
		case []interface{}:
//			fmt.Println(vv, "is []interface, len=", len(vv))
			variabilityList[index1].VariationPointList[index2].VariabilityName = make([]string, len(vv))
			index3 := 0
			for _, v := range vv {
//				fmt.Println(v, "is string")
				variabilityList[index1].VariationPointList[index2].VariabilityName[index3] = v.(string)
				index3++
			}
		case string:
//			fmt.Println(vv, "is a string")
			variabilityList[index1].VariationPointList[index2].VariabilityName = make([]string, 1)
			variabilityList[index1].VariationPointList[index2].VariabilityName[0] = variabilityNameArrayMap.(string)
		default:
			fmt.Println(vv, "is of an unknown type")
		}
	return variabilityList
}

func decodeInstanceConfigs(instanceConfigs string) []Instance {
//fmt.Printf("instanceConfigs=|%s|\n", instanceConfigs)
	var instanceList []Instance
	var instanceMap = make(map[string]interface{})
	err := json.Unmarshal([]byte(instanceConfigs), &instanceMap)
	if err != nil {
		fmt.Printf("decodeInstanceConfigs():unmarshal error=%s\n", err)
		return nil
	}
	instanceList = unpackInstMapLevel1(instanceMap, instanceList)
	return instanceList
}

func unpackInstMapLevel1(instanceMap map[string]interface{}, instanceList []Instance) []Instance {
	i := 0
	for k, v := range instanceMap {
		instanceList = append(instanceList, Instance{})
		instanceList[i].InstanceName = k
//fmt.Printf("k1=%s\n", k)
		switch vv := v.(type) {
			case []interface{}:
//				fmt.Println(vv, "is an array:, len=", len(vv))
				instanceList = unpackInstMapLevel2(i, vv, instanceList)
			default:
				fmt.Println(vv, "is of an unknown type")
		}
		i++
	}
	return instanceList
}

func unpackInstMapLevel2(index1 int, instDefMap []interface{}, instanceList []Instance) []Instance {
	for k, v := range instDefMap {  // range should always be 2; k==0->RowDef, k==1->RowColumnDef
		switch vv := v.(type) {
		case string:
//			fmt.Println(vv, "is a string:")
			if k == 0 {  // RowDef
				rowName := expandRowColumnName(vv)
				instanceList[index1].Row = make([]RowDef, len(rowName))
				for i := 0; i < len(rowName); i++ {
					instanceList[index1].Row[i].RowName = rowName[i]
				}
			} else { //RowColumnDef, corner case with one row
//				fmt.Println(vv, "is a row-column def string:")
				columnName := expandRowColumnName(vv)
				instanceList[index1].RowColumn = make([]RowColumnDef, 1)
				instanceList[index1].RowColumn[0].Column = make([]ColumnDef, len(columnName))
				for i := 0; i < len(columnName); i++ {
					instanceList[index1].RowColumn[0].Column[i].ColumnName = columnName[i]
				}
			}
		case []interface{}:
//				fmt.Println(vv, "is an []interface:, len=", len(vv))
				if k == 0 {
					instanceList[index1].Row = make([]RowDef, len(vv))
					for k2, v2 := range vv {
//						fmt.Println(v2, "is string")
						instanceList[index1].Row[k2].RowName = v2.(string)
					}
				} else {
//					fmt.Println(vv, "is an row-column []interface:, len=", len(vv))
					instanceList[index1].RowColumn = make([]RowColumnDef, len(vv))
					for i := 0; i < len(vv); i++ {
						instanceList = unpackInstMapLevel3(index1, i, vv[i], instanceList)
					}
				}
		default:
			fmt.Println(vv, "is of an unknown type")
		}
	}
	return instanceList
}

func expandRowColumnName(codedName string) []string {   // either an xxx[a,b] or single "xyz"
	var name []string
	frontBracketIndex := strings.Index(codedName, "[")
	if frontBracketIndex != -1 {
		baseName := codedName[:frontBracketIndex]
		a, b := extractNameSuffixBoundaries(codedName[frontBracketIndex+1:len(codedName)-1])
		if a < 0 || b < a {
			fmt.Printf("expandRowColumnName: Invalid name suffix boundaries a=%d, b=%d\n", a, b)
			return nil
		}
		name = make([]string, b-a+1)
		for i := 0; i < b-a+1; i++ {
			name[i] = baseName + strconv.Itoa(a+i)
//fmt.Printf("rowColumnName[%d]=%s\n", i, name[i])
		}
	} else {
		name = make([]string, 1)
		name[0] = codedName
	}
	return name
}

func extractNameSuffixBoundaries(codedBoundaries string) (int, int) {  // a,b
	commaIndex := strings.Index(codedBoundaries, ",")
	if commaIndex == -1 {
		fmt.Printf("Decoding row name boundaries failed, encoding=%s\n", codedBoundaries)
		return -1, -1
	}
	aStr := strings.TrimSpace(codedBoundaries[:commaIndex])
	bStr := strings.TrimSpace(codedBoundaries[commaIndex+1:])
	a, err := strconv.Atoi(aStr)
	if err != nil {
		fmt.Printf("Converting row name boundary index a failed, a=%s\n", aStr)
		return -1, -1
	}
	b, err := strconv.Atoi(bStr)
	if err != nil {
		fmt.Printf("Converting row name boundary index b failed, a=%s\n", bStr)
		return -1, -1
	}
	return a, b
}

func unpackInstMapLevel3(index1 int, index2 int, columnDefMap interface{}, instanceList []Instance) []Instance {
//	fmt.Println(columnDefMap, "is a column def")
	switch vv := columnDefMap.(type) {
	case string:
//		fmt.Println(vv, "is a string")
		columnName := expandRowColumnName(vv)
		instanceList[index1].RowColumn[index2].Column = make([]ColumnDef, len(columnName))
		for i := 0; i < len(columnName); i++ {
			instanceList[index1].RowColumn[index2].Column[i].ColumnName = columnName[i]
		}
	case []interface{}:
//		fmt.Println(vv, "is an []interface:, len=", len(vv))
		instanceList[index1].RowColumn[index2].Column = make([]ColumnDef, len(vv))
		for i := 0; i < len(vv); i++ {
//			fmt.Println(vv[i], "is a string")
			instanceList[index1].RowColumn[index2].Column[i].ColumnName = vv[i].(string)
		}
	default:
		fmt.Println(vv, "is of an unknown type")
	}
	return instanceList
}

func readConfigFile(vspecDir string) (string, string) {
	variantReadError := false
	instanceReadError := false
	configFileName := vspecDir + "himConfiguration.json"
	data, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Printf("readConfigFile: Could not read %s\n", configFileName)
		return "", ""
	}
	variantsIndex := strings.Index(string(data), `"variants":`)
	if variantsIndex == -1 {
		fmt.Printf("readConfigFile: Could not find 'variants' key in %s\n", configFileName)
		variantReadError = true
	}
	instancesIndex := strings.Index(string(data), `"instances":`)
	if instancesIndex == -1 {
		fmt.Printf("readConfigFile: Could not find 'instances' key in %s\n", configFileName)
		instanceReadError = true
	}
	variantReturn := ""
	if !variantReadError {
		variantsStr := string(data[variantsIndex+11:instancesIndex])
		variantsStr = strings.TrimSpace(variantsStr)
		variantReturn = variantsStr[:len(variantsStr)-1]
	}
	instancesReturn := ""
	if !instanceReadError {
		instancesStr := string(data[instancesIndex+12:])
		instancesStr = strings.TrimSpace(instancesStr)
		instancesReturn = strings.TrimSpace(instancesStr[:len(instancesStr)-1])
	}
	return variantReturn, instancesReturn
}

func getRootVspecFileName(vspecRootDir string) string {
	entries, err := os.ReadDir(vspecRootDir)
	if err != nil {
		fmt.Printf("getRootVspecFileName:error=%s\n", err)
		return ""
	}
	for _, e := range entries {
		if !e.IsDir() && strings.Contains(e.Name(), ".vspec") { //Should only be one vspec file in root dir
			return e.Name()
		}
	}
	return ""
}

func readEnumDefinitions(fileName string) []PropertyData {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("readEnumDefinitions: Could not read %s\n", fileName)
		return nil
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	continueScan := true
	var enumDefs []PropertyData
	enumIndex := 0
	var structDefs []StructData
	structIndex := -1
	fieldIndex := 0
	var thisNode PropertyData
	nextNodeName := ""
	populateStruct := false
	for continueScan {
		nextNodeName, thisNode, continueScan = getNode(scanner, nextNodeName)
		switch thisNode.NodeType {
			case "branch":
				if populateStruct {
					populateStruct = false
				}
			case "struct":
				populateStruct = true
				structIndex++
				structDefs = append(structDefs, StructData{})
				structDefs[structIndex].Name = thisNode.Name
			case "property":
				if populateStruct {
					structDefs[structIndex].Property = append(structDefs[structIndex].Property, PropertyData{})
					structDefs[structIndex].Property[fieldIndex] = thisNode
					fieldIndex++					
				} else {
					enumDefs = append(enumDefs, PropertyData{})
					enumDefs[enumIndex] = thisNode
					enumIndex++					
				}
			default:
				fmt.Printf("readEnumDefinitions: Invalid nodetype=%s\n", thisNode.NodeType)
		}
	}
	return enumDefs
}

func getNode(scanner *bufio.Scanner, nextNodeName string) (string, PropertyData, bool) {
	var line string
	continueScan := true
	thisNode := clearPropertyNode(nextNodeName)
	nextLine := ""
	nodeComplete := false
	for continueScan && !nodeComplete {
		if len(nextLine) == 0 {
			continueScan = scanner.Scan()
			line = scanner.Text()
		} else {
			line = nextLine
			nextLine = ""
		}
		key, value := analyzeLine(line)
		switch key {
			case "name":
				if len(thisNode.Name) == 0 {
					thisNode.Name = value
				} else{
					nextNodeName = value
					nodeComplete = true
				}
			case "type":
				thisNode.NodeType = value
			case "datatype":
				thisNode.Datatype = value
			case "allowed":
				thisNode.Allowed, nextLine, continueScan = getAllowedValues(scanner)
			case "min":
				thisNode.Min = value
			case "max":
				thisNode.Max = value
			case "unit":
				thisNode.Unit = value
			case "skipline": continue
		}
	}
	return nextNodeName, thisNode, continueScan
}

func analyzeLine(line string) (string, string) {
	if len(line) > 0 && line[len(line)-1] == ':' && line[0] != ' ' {
		return "name", line[:len(line)-1]
	}
	if strings.Contains(line, "datatype:") {
		return "datatype", extractValue(line)
	}
	if strings.Contains(line, "type:") {
		return "type", extractValue(line)
	}
	if strings.Contains(line, "allowed:") {
		return "allowed", ""
	}
	if strings.Contains(line, "min:") {
		return "min", extractValue(line)
	}
	if strings.Contains(line, "max:") {
		return "max", extractValue(line)
	}
	if strings.Contains(line, "unit:") {
		return "unit", extractValue(line)
	}
	return "skipline", ""
}

func getAllowedValues(scanner *bufio.Scanner) ([]string, string, bool) {
	var line string
	continueScan := true
	var allowedValues []string
	for continueScan {
		continueScan = scanner.Scan()
		line = scanner.Text()
		if strings.Contains(line, "- ") {
			allowedValues = append(allowedValues, line)
		} else {
			return allowedValues, line, continueScan
		}
	}
	return allowedValues, "", continueScan
}

func extractValue(line string) string {
	colonIndex := strings.Index(line, ":")
	return strings.TrimSpace(line[colonIndex+1:])
}

func clearPropertyNode(nextNodeName string) PropertyData {
	var propertyNode PropertyData
	propertyNode.Name = nextNodeName
	propertyNode.NodeType = ""
	propertyNode.Datatype = ""
	propertyNode.Allowed = nil
	propertyNode.Min = ""
	propertyNode.Max = ""
	propertyNode.Unit = ""
	return propertyNode
}

func main() {
	parser := argparse.NewParser("print", "HIM configurator")
	makeCommand := parser.Selector("m", "makecommand", []string{"all", "yaml", "json", "csv", "binary"}, &argparse.Options{Required: false,
		Help: "Make command parameter must be either: all, yaml, csv, or binary", Default: "all"})
//	configFileName := parser.String("p", "pathconfigfile", &argparse.Options{Required: false, Help: "path to configuration file", Default: "himConfiguration.json"})
	vspecDir := parser.String("v", "vspecdir", &argparse.Options{Required: false, Help: "path to vspec root directory", Default: "Vehicle/CargoVehicle/"})
	sConf := parser.Flag("c", "saveconf", &argparse.Options{Required: false, Help: "Saves the configured vspec file with extension .conf", Default: false})
	enumSubst := parser.Flag("e", "enumSubstitute", &argparse.Options{Required: false, Help: "Substitute enum links to Datatype tree with actual datatypes", Default: false})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	saveConf = *sConf
	makeCmd = *makeCommand
	if *enumSubst {
		if !fileExists(*vspecDir + "Datatypes.yaml") {
			cmd := exec.Command("/usr/bin/bash", "make.sh", "yaml", "./spec/objects/Datatype/Datatype.vspec")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Executing make failed with error=%s\n", err)
			} else {
				err = os.Rename("../../vss_rel_0.1-dev.yaml", *vspecDir+"Datatypes.yaml")  //TODO: VERSION=0.1-dev no will not be persistent
				if err != nil {
					fmt.Printf("Failed to rename and move %s error=%s\n", *vspecDir+"Datatypes.yaml", err)
				}
			}
		}
	}
	enumSubstitute = *enumSubst && fileExists(*vspecDir + "Datatypes.yaml")

	variantConfigs, instanceConfigs := readConfigFile(*vspecDir)
	if variantConfigs != "" {
		variantList = decodeVariantConfigs(variantConfigs)
	}
	if instanceConfigs != "" {
		instanceList = decodeInstanceConfigs(instanceConfigs)
	}
	variabilityList = readVariabilityFile(*vspecDir + "Variability.json")

	err = filepath.WalkDir(*vspecDir, walkVariantPass)
	if err != nil {
		fmt.Printf("Variant preprocessing failed. Terminating.\n")
		os.Exit(1)
	}

	err = filepath.WalkDir(*vspecDir, walkInstancePass)
	if err != nil {
		fmt.Printf("Instance preprocessing failed. Terminating.\n")
		os.Exit(1)
	}

	if enumSubstitute {
		enumData = readEnumDefinitions(*vspecDir + "Datatypes.yaml")
/*fmt.Printf("len(enumData)=%d\n", len(enumData))
for i := 0; i < len(enumData); i++ {
	fmt.Printf("enumData[%d].Name=%s\n", i, enumData[i].Name)
	fmt.Printf("enumData[%d].Datatype=%s\n", i, enumData[i].Datatype)
	fmt.Printf("len(enumData[%d].Allowed)=%d\n", i, len(enumData[i].Allowed))
}*/
		err = filepath.WalkDir(*vspecDir, walkEnumSubstitute)
		if err != nil {
			fmt.Printf("Enum substitute preprocessing failed. Terminating.\n")
			os.Exit(1)
		}
	}

	rootVspecFileName := getRootVspecFileName(*vspecDir)
	if makeCmd == "all" {
		makeCmd = ""
	}
	cmd := exec.Command("/usr/bin/bash", "make.sh", makeCmd, "./spec/trees/" + *vspecDir+rootVspecFileName)
//	cmd := exec.Command("/usr/bin/bash", "make.sh", makeCmd, *vspecDir+rootVspecFileName)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Executing make failed with error=%s. Terminating.\n", err)
		os.Exit(1)
	}

	filepath.WalkDir(*vspecDir, walkPostmake)
	if err == nil {
		fmt.Printf("\nMake command output from configured vspec file in %s is available in the root directory.\n", *vspecDir)
	}
}
