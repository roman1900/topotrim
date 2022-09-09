package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
)

type TransformType struct {
	Scale     []float64 `json:"scale"`
	Translate []float64 `json:"translate"`
}
type Geometry struct {
	Arcs       [][]interface{}        `json:"arcs"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}
type Object struct {
	Type     string     `json:"type"`
	Geometry []Geometry `json:"geometries"`
}
type TopoJSON struct {
	Type      string            `json:"type"`
	Arcs      [][][]int32       `json:"arcs"`
	Transform TransformType     `json:"transform"`
	Objects   map[string]Object `json:"objects"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var topojson TopoJSON
	var topoRecon TopoJSON
	file := flag.String("i", "", "input file")
	field := flag.String("m", "POA_CODE16", "geometry property to match on")
	state := flag.String("s", "7", "string to match on the property")
	output := flag.String("o", "output.json", "file to write the output to")
	flag.Parse()
	if len(*file) > 0 {
		topo, err := os.ReadFile(*file)
		check(err)
		err = json.Unmarshal(topo, &topojson)
		check(err)
		fmt.Println("Objects Found:")
		for k, v := range topojson.Objects {
			fmt.Println(k)
			topoRecon.Objects = make(map[string]Object)
			topoRecon.Objects[k] = Object{}
			if en, ok := topoRecon.Objects[k]; ok {
				en.Type = v.Type
				topoRecon.Objects[k] = en
			}
			for _, geo := range v.Geometry {
				for ke, va := range geo.Properties {
					if ke == *field {
						if va.(string)[:1] == *state {
							fmt.Println("TASSIE!!!!!!!")
							if en, ok := topoRecon.Objects[k]; ok {
								en.Geometry = append(en.Geometry, geo)
								topoRecon.Objects[k] = en
							}
						}
					}
					//fmt.Println(ke, v)
				}
				for _, f := range geo.Arcs {
					for _, s := range f {
						switch s := s.(type) {
						case float64:
							fmt.Println(s)
						case []interface{}:
							fmt.Println(s)
						default:
							fmt.Println(reflect.TypeOf(s))
						}
					}
				}
			}
		}

		topoRecon.Type = topojson.Type
		topoRecon.Transform = topojson.Transform
		topoRecon.Arcs = topojson.Arcs
		b, err := json.Marshal(topoRecon)
		check(err)
		err = os.WriteFile(*output, b, 0666)
		check(err)
	} else {
		fmt.Println("A json file must be provided using the -i flag")
	}

}
