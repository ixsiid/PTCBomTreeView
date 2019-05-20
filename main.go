package main

import (
	"fmt"
	"os"
	"io/ioutil"
//	"reflect"
	"strings"
	"encoding/json"
	"errors"
	"strconv"
	"regexp"
    "path/filepath"
)

func getFileNameWithoutExt(path string) string {
    return path[:len(path)-len(filepath.Ext(path))]
}

type Translator struct{
	wordA []string
	wordB []string
	size int
}

func (t *Translator) Read(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error: ", file)
		return
	}

	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	var text string = string(buffer)

	lines := strings.Split(strings.Replace(text, "\r\n", "\n", -1), "\n")
	t.wordA = make([]string, len(lines))
	t.wordB = make([]string, len(lines))
	t.size = 0
	for _, line := range lines {
		words := strings.Split(line, ":")
		if len(words) != 2 {
			continue
		}

		t.wordA[t.size] = words[0]
		t.wordB[t.size] = words[1]
		t.size++
	}
}

func (t Translator) Translate(text string, start string, end string) string {
	for i := 0; i < t.size; i++ {
		text = strings.Replace(text, start + t.wordA[i] + end, t.wordB[i], -1)
	}
	return text
}

func (t Translator) PrintAll() {
	for i := 0; i < t.size; i++ {
		fmt.Println(t.wordA[i], "->", t.wordB[i])
	}
}

func main() {
	var t Translator
	t.Read("translate.txt")
		
	f, err := os.Open("template.html")
	if err != nil {
		fmt.Println("Error: can not open template.html")
		os.Exit(1)
	}

	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	template := string(buffer)

	for i, arg := range os.Args{

		if i == 0 {
			fmt.Println("Generate bom tree")
			continue
		}
		fmt.Println(i, ": Start '", arg, "'")
		
		arg = filepath.Clean(arg)
		var configure Translator
		configure.Read(filepath.Join(filepath.Dir(arg), "template.conf"))

		
		tree, err := generate(arg, t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		json, _ := json.Marshal(tree)
		jsonString := string(json)

		htmlFile := getFileNameWithoutExt(arg) + ".html"
		fmt.Println(htmlFile)
		html := strings.Replace(template, `{"name":"$JSON_DATA$"}`, jsonString, -1)
		html = configure.Translate(html, "$", "$")


		ioutil.WriteFile(htmlFile, []byte(html), os.ModePerm)

		fmt.Println("Writed " + htmlFile)
	}
	fmt.Println("Done.")
}

type Element struct {
	Quantity int    `json:"quantity"`
	Type     string `json:"type"`
    Name     string `json:"name"`
	Project  string `json:"project"`
	Number   string `json:"number"`
    Class    string `json:"class"`
    Supplier string `json:"supplier"`
    Material string `json:"material"`
	Finish   string `json:"finish"`
	Child []*Element `json:"child"`
}

func generate(file string, t Translator) (Element, error) {
	f, err := os.Open(file)
	if err != nil {
		return Element{}, errors.New("Error: " + file)
	}

	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	var UTF8_BOM = string([]byte{239, 187, 191})
	var text string = strings.Replace(strings.TrimLeft(string(buffer), UTF8_BOM), "\r\n", "\n", -1)
	text = t.Translate(text, "", "")

	subAssembly := map[string]*Element{}
	var root *Element
	sections := strings.Split(text, "\n\n")
	space := regexp.MustCompile(`\s+`)
	for _, section := range sections {
		elements := strings.Split(section, "\n")
		attrs := strings.Split(elements[0], " ")
		parentType := attrs[0]
		title := attrs[1]
		elements = elements[1:]

		var p Element
		t := &p
		t.Name = title

		if parentType == "Assembly" {
			root = t
		}

		if parentType == "SubAssembly" {
			t = subAssembly[t.Name]
		}

		if parentType == "PartsList" {
			continue
		}

		for _, element := range elements {
			var p Element
			values := strings.Split(space.ReplaceAllString(strings.Trim(element, " "), " "), " ")
			if len(values) < 2 {
				continue
			}
			values = append(values, []string{"", "", "", "", "", "", "", "", ""}...)
			p.Quantity, _ = strconv.Atoi(values[0])
			p.Type = values[1]
			p.Name = values[2]
			p.Project = values[3]
			p.Number = values[4]
			p.Material = values[5]
			p.Class = values[6]
			p.Supplier = values[7]
			p.Finish = values[8]

			t.Child = append(t.Child, &p)
			if p.Type == "SubAssembly" {
				subAssembly[p.Name] = &p
			}
		}
	}

	return *root, nil
}

