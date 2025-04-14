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
	"sort"
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
var configFileName string
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

var overlayToolsParameter string   // VSS-ools overlay parameter

var enumData []PropertyData

var configData ConfigValues
type ConfigValues struct {
	Default []ConfigValueDefinition  `json:"default"`
	Description []ConfigValueDefinition  `json:"description"`
}
type ConfigValueDefinition struct {
	Path string  `json:"path"`
	Value string `json:"value"`
}

type VspecContext struct {
	BasePath string
	Path string
	FName string
	KeyValue string
}

type StructData struct {
	Name string
	Property []PropertyData
}

var structData []StructData

func variantProcess(sourceFile string) error {  // sourceFile input is always vspec2 file, but if earlier processing has created a vspec file that is used
	extensionIndex := strings.Index(sourceFile, ".vspec2")
	if fileExists(sourceFile[:extensionIndex] + ".vspec") {
		sourceFile = sourceFile[:extensionIndex] + ".vspec"
	}
	sourceFp, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("variantProcess:Error reading %s: %s\n", sourceFile, err)
		return err
	}
	scanner := bufio.NewScanner(sourceFp)
	scanner.Split(bufio.ScanLines)
	var text string
	continueScan := true
	var newFileLines []string  // this is the content to save at the end of this configuration
	var savedLines []string
	for continueScan {
		continueScan = scanner.Scan()
		text = scanner.Text()
		if strings.Contains(text, "VariationPoint:") {
			fmt.Printf("Variation point line found in file=%s:%s\n", sourceFile, text)
			commentIndex := strings.Index(text, "#")
			if commentIndex == -1 {
				fmt.Printf("Error no comment found in line=%s\n", text)
				continue
			}
			variationType := text[commentIndex+1:]
			var variationLines []string
			var saveLine string
			variationLines, saveLine, scanner = readVariations(scanner)
			newFileLines = updateVariation(newFileLines, sourceFile, savedLines, variationLines, variationType, variabilityList, variantList)
			savedLines = nil
			savedLines = append(savedLines,saveLine)
		} else {
			savedLines = append(savedLines,text)
		}
	}
	newFileLines = copyRemainingLines(newFileLines, savedLines)
	sourceFp.Close()
	err = updateVspec(sourceFile, newFileLines)

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

func updateVariation(newFileLines []string, sourcefile string, savedLines []string, variationLines []string, variationType string, variabilityList []Variability, variantList []Variant) []string {
	for i := 0; i < len(savedLines); i++ {
		newFileLines = append(newFileLines,savedLines[i] + "\n")
	}
	newFileLines = addVariation(newFileLines, variationLines, variationType, variabilityList, variantList)
	return newFileLines
}

func addVariation(newFileLines []string, variationLines []string, variationType string, variabilityList []Variability, variantList []Variant) []string {
	var selectedVariations []string
	for i := 0; i < len(variantList); i++ {
		if variationType == variantList[i].VariantType {
			for j := 0; j < len(variabilityList); j++ {
				if getRootSegment(variationType) == variabilityList[j].VariabilityType {
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
					newFileLines = append(newFileLines,variationLines[j][commentIndex:] + "\n")
					fmt.Printf("Variant %s: Inserted:%s\n", selectedVariations[i], variationLines[j][commentIndex:])
				}
			}
		}
	}
	return newFileLines
}

func getRootSegment(dotLimitedName string) string {
	dotIndex := strings.Index(dotLimitedName, ".")
	if dotIndex == -1 {
		return dotLimitedName
	}
	return dotLimitedName[:dotIndex]
}

func copyRemainingLines(newFileLines []string, savedLines []string) []string {
	for i := 0; i < len(savedLines); i++ {
		newFileLines = append(newFileLines,savedLines[i] + "\n")
	}
	return newFileLines
}

// sourceFile input is always vspec2 file. 
//As the asymmetric instance config is the first to run it hall remove any existing "old" vspec
func instanceProcess(sourceFile string) error {
	extensionIndex := strings.Index(sourceFile, ".vspec2")
	if fileExists(sourceFile[:extensionIndex] + ".vspec") {
		err := os.Remove(sourceFile[:extensionIndex] + ".vspec")
		if err != nil {
			fmt.Printf("Deleting %s failed. Err=%s\n", sourceFile, err)
			return err
		}
	}
	file, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("instanceProcess:Error reading %s: %s\n", sourceFile, err)
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
			fmt.Printf("Instance line found in file=%s:%s\n", sourceFile, text)
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
							nodeNameIndex := strings.Index(subTree[j], nodeName) + len(nodeName)
							savedLines = append(savedLines, subTree[j][:nodeNameIndex] + "." + getRowInstance(instanceTag,i) + subTree[j][nodeNameIndex:])
						} else if strings.Contains(subTree[j], "LocalVP:") {
							sharpIndex := strings.Index(subTree[j], "#")
							localVpTag := subTree[j][sharpIndex+1:]
							savedLines = append(savedLines,"VariationPoint: #" + localVpTag + "." + getRowInstance(instanceTag,i))
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
					instSubTree := expandSubTree(subTree, filepath.Dir(sourceFile)+"/", nodeName, i, instanceTag)
					savedLines = addInstanceBranch(savedLines, nodeName, i, instanceTag, "")
					for j := 0; j < len(instSubTree); j++ {  // followed by subtree with configured instance
						instExpTree[0] = instSubTree[j]
						instanceExp := getInstanceExpression(instExpTree, 1, instanceTag)
						if len(instanceExp) > 0 {
							savedLines = append(savedLines,"  " + createConfiguredInstance(instanceExp, i, instanceTag))
						}  else if strings.Contains(instSubTree[j], "LocalVP") {
							savedLines = append(savedLines,"\n")
							sharpIndex := strings.Index(instSubTree[j], "#")
							localVpTag := instSubTree[j][sharpIndex+1:]
							savedLines = append(savedLines,"VariationPoint: #" + localVpTag + "." + getRowInstance(instanceTag,i))
						} else {
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
	if isConfigDone {
		err = updateVspec(sourceFile, savedLines)
	}
	return err
}

func updateVspec(sourceFile string, savedLines []string) error {   // if sourceFile is vspec2 file, create vspec file, else delete vspec file and rewrite with savedLines
	var err error
	extensionIndex := strings.Index(sourceFile, ".vspec2")
	if extensionIndex == -1 {
		extensionIndex = strings.Index(sourceFile, ".vspec")
		err = os.Remove(sourceFile)
		if err != nil {
			fmt.Printf("Deleting %s failed. Err=%s\n", sourceFile, err)
			return err
		}
	}
	var vspecFp *os.File
	vspecFp, err = os.OpenFile(sourceFile[:extensionIndex] + ".vspec", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("Could not create %s\n", sourceFile)
		return err
	}
	for i := 0; i < len(savedLines); i++ {
//fmt.Printf("SavedLine: %s\n", savedLines[i])
		linebreak := "\n"
		if len(savedLines[i]) > 0 && savedLines[i][len(savedLines[i])-1] == '\n' {
			linebreak = ""
		}
		vspecFp.Write([]byte(savedLines[i] + linebreak))
	}
	vspecFp.Close()
	return nil
}


func expandSubTree(subTree []string, path string, rootNodeName string, index int, instanceTag string) []string {
	var expandedTree []string
	for i := 0; i < len(subTree); i++ {
		if strings.Contains(subTree[i], "#include") && strings.Contains(subTree[i], rootNodeName) {
			nodeNameIndex := strings.Index(subTree[i], rootNodeName) + len(rootNodeName)
//fmt.Printf("rootNodeName: %s line=%s\n", rootNodeName, subTree[i])
			updatedIncludeExpression := subTree[i][:nodeNameIndex] + "." + getRowInstance(instanceTag,index) + subTree[i][nodeNameIndex:]
			includeExpansion := readIncludefile(updatedIncludeExpression, path, index, instanceTag)
			for j := 0; j < len(includeExpansion); j++ {
				expandedTree = append(expandedTree, includeExpansion[j])
			}
		} else {
			newNodeName := ""
			getNodeName(subTree[i], &newNodeName)
			if len(newNodeName) > 0 {
				dotIndex := strings.Index(subTree[i], ".")
				expandedTree = append(expandedTree, subTree[i][:dotIndex] + "." + getRowInstance(instanceTag,index) + subTree[i][dotIndex:])
			} else {
				expandedTree = append(expandedTree, subTree[i])
			}
		}
	} 
	return expandedTree
}

func readIncludefile(includeDirective string, path string, index int, instanceTag string) []string { 
// if config instance directive found, update it with config data, update node names and #include directives with rootnode data
	vspecFile, nodeNamePrefix, doInclude := decodeIncludeDirective(includeDirective)
	var includeLines []string
	if !doInclude {
		includeLines = append(includeLines, includeDirective)
//		includeLines = append(includeLines, includeDirective + "." + getRowInstance(instanceTag,index))  //Fixas???
		return includeLines
	}
	file, err := os.Open(path + vspecFile)
	if err != nil {
		fmt.Printf("readIncludefile:Error reading %s: %s\n", path + vspecFile, err)
		return nil
	}
//fmt.Printf("vspecfile: %s\n", path + vspecFile)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var line string
	continueScan := true
	newNodeName := ""
	for continueScan {
		continueScan = scanner.Scan()
		line = scanner.Text()
		getNodeName(line, &newNodeName)
		if len(line) == 0 || line[0] != '#' || strings.Contains(line, "#include") {
			if len(newNodeName) > 0 {
				line = nodeNamePrefix + "." + line
			} else if strings.Contains(line, "#include") {
				incFields := strings.Fields(line)
				line = incFields[0] + " " + incFields[1] + " " + nodeNamePrefix + "." + incFields[2]
			}
			includeLines = append(includeLines, line)
		}
		newNodeName = ""
	}
	return includeLines
}

func decodeIncludeDirective(includeDirective string) (string, string, bool) {  // Syntax: #include Axle.vspec Chassis.Axle, or - vpTag #include Axle.vspec Chassis.Axle
	fields := strings.Fields(includeDirective)
	if len(fields) == 3 {
		return fields[1], fields[2], true
	}
	if len(fields) == 5 {
		return fields[3], fields[4], false
	}
	return "", "", false
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
//	savedLines = append(savedLines, "  description: " + nodeName  + "." + getRowInstance(instanceTag,index))
	savedLines = append(savedLines, "  description: " + instanceTag + ", " + getRowInstance(instanceTag,index) + " instance")
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
		if len(tree[i]) > 0 && (strings.Contains(tree[i], "#include") || strings.Contains(tree[i], "LocalVP") || len(newNodeName) > 0) {
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
		if filepath.Ext(s) == ".vspec2" {
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
		if filepath.Ext(s) == ".vspec2" {
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
		if filepath.Ext(s) == ".vspec2" {
//			fmt.Printf("enumProcess:Vspec path=%s\n", s)
			err = enumProcess(s)
		}
	}
	return err
}

func enumProcess(sourceFile string) error {  // sourceFile input is always vspec2 file, but if earlier processing has created a vspec file that is used
	extensionIndex := strings.Index(sourceFile, ".vspec2")
	if fileExists(sourceFile[:extensionIndex] + ".vspec") {
		sourceFile = sourceFile[:extensionIndex] + ".vspec"
	}
	file, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("enumProcess:Error reading %s: %s\n", sourceFile, err)
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
	if isConfigDone {   // if sourceFile is vspec2 file, create vspec file, else delete vspec file and rewrite with savedLines
		extensionIndex = strings.Index(sourceFile, ".vspec2")
		if extensionIndex == -1 {
			extensionIndex = strings.Index(sourceFile, ".vspec")
			err = os.Remove(sourceFile)
			if err != nil {
				fmt.Printf("Deleting %s failed. Err=%s\n", sourceFile, err)
				return err
			}
		}
		var vspecFp *os.File
		vspecFp, err = os.OpenFile(sourceFile[:extensionIndex] + ".vspec", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("Could not create %s\n", sourceFile)
			return err
		}
		for i := 0; i < len(savedLines); i++ {
//fmt.Printf("SavedLine: %s\n", savedLines[i])
			vspecFp.Write([]byte(savedLines[i] + "\n"))
		}
		vspecFp.Close()
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
		if filepath.Ext(s) == ".vspec2" && !saveConf {
			extensionIndex := strings.Index(s, ".vspec2")
			err := os.Remove(s[:extensionIndex] + ".vspec")
			if err != nil {
				fmt.Printf("Failed to remove %s" + ".vspec\n", s[:extensionIndex])
			}
		} else if filepath.Ext(s) == ".orig" && !saveConf {
			extensionIndex := strings.Index(s, ".orig")
			err := os.Remove(s[:extensionIndex])
			if err != nil {
				fmt.Printf("Failed to remove %s\n", s[:extensionIndex])
			}
			err = os.Rename(s, s[:len(s) - 5])  // truncate ".orig"
			if err != nil {
				fmt.Printf("Failed to rename %s, error=%s\n", s, err)
			}
		}
	}
	return nil
}

func decodeOverlayConfigs(overlayConfigs string, vspecDir string) string {
	var overlays []string
	err := json.Unmarshal([]byte(overlayConfigs), &overlays)
	if err != nil {
		fmt.Printf("decodeOverlayConfigs():unmarshal %s, error=%s\n", overlayConfigs, err)
		return ""
	}
	toolsParameter := ""
	for i := 0; i < len(overlays); i++ {
		toolsParameter = toolsParameter + "-l " + vspecDir + overlays[i] + " "
	}
	return toolsParameter
}

func decodeVariantConfigs(variantConfigs string) []Variant { //JSON object:{"var-type1":"var-name1", ., "var-typeN":"var-nameN"}
	var variantList []Variant
	var variant Variant
	var variantMap = make(map[string]interface{})
	err := json.Unmarshal([]byte(variantConfigs), &variantMap)
	if err != nil {
		fmt.Printf("decodeVariantConfigs():unmarshal %s, error=%s\n", variantConfigs, err)
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

func readConfigFile(vspecDir string) []string {
	configFileName := vspecDir + configFileName
	data, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Printf("readConfigFile: Could not read %s\n", configFileName)
		return nil
	}
	variantsIndex := strings.Index(string(data), `"variants":`)
	if variantsIndex == -1 {
		fmt.Printf("readConfigFile: Could not find 'variants' key in %s\n", configFileName)
	}
	instancesIndex := strings.Index(string(data), `"instances":`)
	if instancesIndex == -1 {
		fmt.Printf("readConfigFile: Could not find 'instances' key in %s\n", configFileName)
	}
	overlaysIndex := strings.Index(string(data), `"overlays":`)
	if overlaysIndex == -1 {
		fmt.Printf("readConfigFile: Could not find 'overlays' key in %s\n", configFileName)
	}
	keyIndex := []int{variantsIndex, instancesIndex, overlaysIndex}
	sort.Ints(keyIndex)
	returnString := []string{"", "", ""}
	for i := 0; i < len(keyIndex); i++ {
		if keyIndex[i] != -1 {
			if i+1 < len(keyIndex) {
				returnString[i] = string(data[keyIndex[i]:keyIndex[i+1]-1])
				returnString[i] = strings.TrimSpace(returnString[i])
				returnString[i] = returnString[i][:len(returnString[i])-1]
			} else {
				returnString[i] = string(data[keyIndex[i]:])
				returnString[i] = strings.TrimSpace(returnString[i])
				returnString[i] = returnString[i][:len(returnString[i])-1]
			}
		}
	}
	return returnString
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

func readConfigValues(fileName string) (ConfigValues, error) {
	var configData ConfigValues
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("readConfigValues: Could not read %s\n", fileName)
		return configData, err
	}
	err = json.Unmarshal(data, &configData)
	if err != nil {
		fmt.Printf("readConfigValues:error=%s\n", err)
		return configData, err
	}
	return configData, nil
}

func defaultConfigIndex(path string, arrayPtr *[]ConfigValueDefinition) int {
	for i := 0; i < len(*arrayPtr); i++ {
		if (*arrayPtr)[i].Path == path {
			return i
		}
	}
	return -1
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
			case "sensor": fallthrough  //sensor used as vss-tools reject property in combination with allowed...
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
		Help: "Make command parameter must be either: all, yaml, json, csv, or binary", Default: "yaml"})
	confFName := parser.String("c", "configfile", &argparse.Options{Required: false, Help: "configuration file name", Default: "himConfig-truck.json"})
	vspecDir := parser.String("r", "rootdir", &argparse.Options{Required: false, Help: "path to vspec root directory", Default: "Vehicle/VSS-core/"})
	sConf := parser.Flag("s", "vspecsave", &argparse.Options{Required: false, Help: "Saves the configured .vspec2 files with extension .vspec"})
	preProcessOnly := parser.Flag("p", "preprocess", &argparse.Options{Required: false, Help: "Pre-process only, save configured vspec files. Do not run VSS-tools."})
	enumSubst := parser.Flag("n", "noEnumSubst", &argparse.Options{Required: false, Help: "No substitution of enum links to Datatype tree with actual datatypes"})
	overlayDisable := parser.Flag("d", "disableOverlays", &argparse.Options{Required: false, Help: "Disables VSS-tools overlay configurations"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	saveConf = *sConf
	makeCmd = *makeCommand
	configFileName = *confFName
	if !*enumSubst {
		if !fileExists(*vspecDir + "Datatypes.yaml") {
			cmd := exec.Command("/usr/bin/bash", "vspecExec.sh", "yaml", "../objects/Datatype/Datatype.vspec", "")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Executing make failed with error=%s\n", err)
			} else {
				err = os.Rename("exporterData/cvis.yaml", *vspecDir+"Datatypes.yaml")
				if err != nil {
					fmt.Printf("Failed to rename and move %s error=%s\n", *vspecDir+"Datatypes.yaml", err)
				}
			}
		}
	}
	enumSubstitute = !*enumSubst && fileExists(*vspecDir + "Datatypes.yaml")

	configurations := readConfigFile(*vspecDir)
	for i := 0; i < len(configurations); i++ {
		if strings.Contains(configurations[i], `"instances":`) {
			instanceList = decodeInstanceConfigs(configurations[i][12:])
		} else if strings.Contains(configurations[i], `"variants":`) {
			variantList = decodeVariantConfigs(configurations[i][11:])
		} else if strings.Contains(configurations[i], `"overlays":`) {
			overlayToolsParameter = decodeOverlayConfigs(configurations[i][11:], *vspecDir)
		}
	}

	variabilityList = readVariabilityFile(*vspecDir + "Variability.json")

	// configure instances <- Should be the first config pass
	err = filepath.WalkDir(*vspecDir, walkInstancePass)
	if err != nil {
		fmt.Printf("Instance preprocessing failed. Terminating.\n")
		os.Exit(1)
	}

	// configure global variation points
	err = filepath.WalkDir(*vspecDir, walkVariantPass)
	if err != nil {
		fmt.Printf("Global variant preprocessing failed. Terminating.\n")
		os.Exit(1)
	}

	// configure enums
	if enumSubstitute {
		enumData = readEnumDefinitions(*vspecDir + "Datatypes.yaml")
		err = filepath.WalkDir(*vspecDir, walkEnumSubstitute)
		if err != nil {
			fmt.Printf("Enum substitute preprocessing failed. Terminating.\n")
			os.Exit(1)
		}
	}

	if *preProcessOnly {
		fmt.Printf("VSS-tools is not called.\nConfigured vspec files are saved\n")
		os.Exit(0)
	}
	// run VSS-tools
	rootVspecFileName := getRootVspecFileName(*vspecDir)
	if makeCmd == "all" {
//		makeCmd = ""
	}
	var cmd *exec.Cmd
	if *overlayDisable {
		cmd = exec.Command("/usr/bin/bash", "vspecExec.sh", makeCmd, *vspecDir+rootVspecFileName, "")
	} else {
		cmd = exec.Command("/usr/bin/bash", "vspecExec.sh", makeCmd, *vspecDir+rootVspecFileName, overlayToolsParameter)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Executing make failed with error=%s. Terminating.\n", err)
		os.Exit(1)
	}

	// file clean up  
	filepath.WalkDir(*vspecDir, walkPostmake)
	if err == nil {
		fmt.Printf("\nMake command output from configured vspec file in %s is available in the exporterData directory.\n", *vspecDir)
	}
	if saveConf {
		fmt.Printf("Configured vspec files are not deleted\n")
	}
}
