package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	txtFilePath = "go-mod-graph.txt"
	//txtFilePath  = "go-mod-graph_short.txt"
	htmlFilePath = "tree.html"

	defaultSymbol     = "circle" // circle | rect | roundRect | triangle | diamond | pin | arrow | none
	defaultSymbolSize = 12

	// -1 (all) | 0 | 1 | .. | N
	initialTreeDepth = 1
	limitDepth       = true
	maxDepth         = 10
)

func main() {
	log.Printf("Load file '%s'", txtFilePath)
	records, loadErr := loadTxt(txtFilePath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}
	log.Printf("%d records in file", len(records))

	log.Println("Prepare map")
	mapRecords, root := prepareMapData(records)
	log.Printf("%d records in map", len(mapRecords))
	//log.Println()
	//for k, v := range mapRecords {
	//	log.Printf("%s : [ %s ]", k, strings.Join(v, " | "))
	//}
	//log.Println()

	//log.Println()
	//log.Println()

	log.Println("Build tree")
	//treeRecords := buildTree(mapRecords, root)
	treeRecords := buildTreeResursive(mapRecords, root)

	log.Println("Plot tree")
	treeChart := plotTree(treeRecords)

	log.Println("Create html")
	pageErr := createHtml(htmlFilePath, treeChart)
	if pageErr != nil {
		log.Fatal(pageErr)
	}
}

func createHtml(filePath string, charts ...components.Charter) error {
	page := components.NewPage()
	//log.Println("Create html - Add charts")
	page.AddCharts(charts...)

	//log.Println("Create html - Create file")
	file, createErr := os.Create(filePath)
	if createErr != nil {
		return createErr
	}

	//log.Println("Create html - Render")
	return page.Render(io.MultiWriter(file))
}

func plotTree(treeData []opts.TreeData) *charts.Tree {
	tree := charts.NewTree()
	tree.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "95vh"}),
		charts.WithTitleOpts(opts.Title{Title: "Golang mod graph example"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: false}),
	)
	tree.AddSeries("tree", treeData).
		SetSeriesOptions(
			charts.WithTreeOpts(
				opts.TreeChart{
					Layout:           "orthogonal", // orthogonal | radial
					Orient:           "LR",         // LR | RL
					Roam:             true,
					InitialTreeDepth: initialTreeDepth,
					Leaves: &opts.TreeLeaves{
						Label: &opts.Label{Show: true, Position: "right"},
						LineStyle: &opts.LineStyle{
							Width: 2,
							Type:  "solid", // "solid" | "dashed" | "dotted"
						},
						//Emphasis: &opts.Emphasis{Label: &opts.Label{Show: true}},
					},
					Left: "10%", Right: "30%", Top: "5%", Bottom: "5%",
				},
			),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "left"}),
		)
	return tree
}

func buildTreeResursive(mapRecords map[string][]string, root string) []opts.TreeData {
	value := 0
	childrenStr := mapRecords[root]
	children := make([]*opts.TreeData, 0)
	for _, child := range childrenStr {
		//log.Printf("Name: %s - Value: %d - Child: %s", root, value, child)
		children = append(children, buildNode(mapRecords, child, value+1, false))
	}
	return []opts.TreeData{ // WARN: this cannot be an array of pointers
		{
			Name:       root,
			Value:      value,
			Collapsed:  false,
			Children:   children,
			Symbol:     defaultSymbol,
			SymbolSize: defaultSymbolSize,
		},
	}
}

func buildNode(mapRecords map[string][]string, name string, value int, collapsed bool) *opts.TreeData {
	//log.Printf("%d - %s", value, name)

	childrenStr := mapRecords[name]
	if childrenStr == nil {
		return &opts.TreeData{Name: name}
	}

	// WARN: tree rendering with too many nested children takes forever
	if limitDepth {
		if value > maxDepth {
			log.Printf("'%s' node is out of max depth (%d)", name, value)
			return &opts.TreeData{Name: name}
		}
	}

	children := make([]*opts.TreeData, 0)
	for _, child := range childrenStr {
		//log.Printf("Name: %s - Value: %d - Child: %s", name, value, child)
		children = append(children, buildNode(mapRecords, child, value+1, collapsed))
	}
	return &opts.TreeData{
		Name:       name,
		Value:      value,
		Collapsed:  collapsed,
		Children:   children,
		Symbol:     defaultSymbol,
		SymbolSize: defaultSymbolSize,
	}
}

func prepareMapData(records []string) (map[string][]string, string) {
	mapRecords := make(map[string][]string)

	root := ""
	for i, rec := range records {
		recSplit := strings.Split(rec, " ")
		if len(recSplit) == 2 {
			if i == 0 {
				root = recSplit[0]
			}
			if mapRecords[recSplit[0]] == nil {
				mapRecords[recSplit[0]] = make([]string, 0)
			}
			mapRecords[recSplit[0]] = append(mapRecords[recSplit[0]], recSplit[1])
		} else {
			log.Printf("Skipped line %s", rec)
		}
	}

	return mapRecords, root
}

func loadTxt(filePath string) ([]string, error) {
	file, openErr := os.Open(filePath)
	if openErr != nil {
		return nil, openErr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	records := make([]string, 0)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}
	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}
	return records, nil
}
