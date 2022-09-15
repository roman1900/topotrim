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

func arcRange(objects map[string]Object, max int) (int, int) {
	start := max
	end := 0
	for _, v := range objects {
		for _, v := range v.Geometry {
			for _, f := range v.Arcs {
				for _, s := range f {
					switch s := s.(type) {
					case float64:
						i := 0
						if s < 0 {
							i = ^int(s)
						} else {
							i = int(s)
						}
						if i < start {
							start = i
						}
						if i > end {
							end = i
						}
					case []interface{}:
						for _, s := range s {
							switch s := s.(type) {
							case float64:
								i := 0
								if s < 0 {
									i = ^int(s)
								} else {
									i = int(s)
								}
								if i < start {
									start = i
								}
								if i > end {
									end = i
								}
							}
						}

					default:
						fmt.Println(reflect.TypeOf(s))
					}
				}
			}
		}
	}
	return start, end
}

func truncateArcs(arcs [][][]int32, start int, end int) [][][]int32 {
	res := [][][]int32{}

	for i := start; i <= end; i++ {
		res = append(res, arcs[i])
	}
	return res
}

func refactorArcs(objects *map[string]Object, offset int) {
	for k, o := range *objects {
		for l, g := range o.Geometry {
			for i, f := range g.Arcs {
				for ia, s := range f {
					switch s := s.(type) {
					case float64:
						if s < 0 {
							(*objects)[k].Geometry[l].Arcs[i][ia] = s + float64(offset)
						} else {
							(*objects)[k].Geometry[l].Arcs[i][ia] = s - float64(offset)
						}

					case []interface{}:
						temp := []float64{}
						for _, s := range s {
							switch s := s.(type) {
							case float64:
								if s < 0 {
									temp = append(temp, s+float64(offset))
								} else {
									temp = append(temp, s-float64(offset))
								}
							}
						}
						(*objects)[k].Geometry[l].Arcs[i][ia] = temp

					default:
						fmt.Println(reflect.TypeOf(s))
					}
				}
			}
		}
	}
}
func main() {
	var topojson TopoJSON
	var topoRecon TopoJSON
	file := flag.String("i", "", "input file")
	field := flag.String("m", "POA_CODE16", "geometry property to match on")
	search := flag.String("s", "7", "string to match on the property")
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
						if va.(string)[:len(*search)] == *search {
							if en, ok := topoRecon.Objects[k]; ok {
								en.Geometry = append(en.Geometry, geo)
								topoRecon.Objects[k] = en
							}
						}
					}
				}
			}
		}

		topoRecon.Type = topojson.Type
		topoRecon.Transform = topojson.Transform
		start, end := arcRange(topoRecon.Objects, len(topojson.Arcs)-1)
		fmt.Println("Arcs to be extracted START:", start, "END:", end)
		topoRecon.Arcs = truncateArcs(topojson.Arcs, start, end)
		fmt.Println("Refactoring Object Arcs")
		refactorArcs(&topoRecon.Objects, start)
		fmt.Println("Converting back to json")
		b, err := json.Marshal(topoRecon)
		check(err)
		fmt.Println("Writing new file to ", *output)
		err = os.WriteFile(*output, b, 0666)
		check(err)
	} else {
		fmt.Println("A json file must be provided using the -i flag")
	}
}
