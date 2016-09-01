package geojson

import (
	"errors"
	"fmt"
	rtreego "github.com/dhconnelly/rtreego"
	gabs "github.com/jeffail/gabs"
	geo "github.com/kellydunn/golang-geo"
	ioutil "io/ioutil"
	"strconv"
	"sync"
)

/*

- gabs is what handles marshaling a random bag of GeoJSON
- rtreego is imported to convert a WOFFeature in to a handy rtreego.Spatial object for indexing by go-whosonfirst-pip
- geo is imported to convert a WOFFeature geometry into a list of geo.Polygon objects for doing containment checks in go-whosonfirst-pip
  (only Polygons and MultiPolygons are supported at the moment)

*/

// See also
// https://github.com/dhconnelly/rtreego#storing-updating-and-deleting-objects

// sudo make me an interface
// (201251207/thisisaaronland)

type WOFSpatial struct {
	bounds     *rtreego.Rect
	Id         int
	Name       string
	Placetype  string
	Offset     int // used when calling EnSpatializeGeom in order to know which polygon we care about
	Deprecated bool
	Superseded bool
}

// sudo make me an interface
// (201251207/thisisaaronland)

type WOFPolygon struct {
	OuterRing     geo.Polygon
	InteriorRings []geo.Polygon
}

func (p *WOFPolygon) CountPoints() int {

	count := len(p.OuterRing.Points())

	for _, r := range p.InteriorRings {
		count += len(r.Points())
	}

	return count
}

func (p *WOFPolygon) Contains(latitude float64, longitude float64) bool {

	pt := geo.NewPoint(latitude, longitude)
	contains := false

	if p.OuterRing.Contains(pt) {
		contains = true
	}

	if contains && len(p.InteriorRings) > 0 {

		wg := new(sync.WaitGroup)

		for _, r := range p.InteriorRings {

			wg.Add(1)

			go func(poly geo.Polygon, point *geo.Point) {

				defer wg.Done()

				/*

					File under yak-shaving: Some way to send an intercept to poly.Contains
					to stop the raycasting if any one of these goroutines gets the answer
					it needs independent the results of the others. Like I said... yaks.
					(20151028/thisisaaronland)
				*/

				if poly.Contains(point) {
					contains = false
				}

			}(r, pt)
		}

		wg.Wait()
	}

	return contains
}

func (sp WOFSpatial) Bounds() *rtreego.Rect {
	return sp.bounds
}

// sudo make me an interface
// (201251207/thisisaaronland)

type WOFFeature struct {
	Parsed *gabs.Container
}

func (wof WOFFeature) Body() *gabs.Container {
	return wof.Parsed
}

func (wof WOFFeature) Dumps() string {
	return wof.Parsed.String()
}

func (wof WOFFeature) Id() int {

	id, ok := wof.id("properties.wof:id")

	if ok {
		return id
	}

	id, ok = wof.id("properties.id")

	if ok {
		return id
	}

	id, ok = wof.id("id")

	if ok {
		return id
	}

	return -1
}

func (wof WOFFeature) id(path string) (int, bool) {

	body := wof.Body()

	var id_float float64
	var id_str string
	var id int

	var ok bool

	// what follows shouldn't be necessary but appears to be
	// for... uh, reasons (20151013/thisisaaronland)

	id_float, ok = body.Path(path).Data().(float64)

	if ok {
		id = int(id_float)
	} else {
		id, ok = body.Path(path).Data().(int)
	}

	// But wait... there's more (20151028/thisisaaronland)

	if !ok {

		id_str, ok = body.Path(path).Data().(string)

		if ok {

			id_int, err := strconv.Atoi(id_str)

			if err != nil {
				ok = false
			} else {
				id = id_int
			}
		}
	}

	if !ok {
		id = -1
	}

	return id, ok
}

func (wof WOFFeature) Name() string {

	name, ok := wof.name("properties.wof:name")

	if ok {
		return name
	}

	name, ok = wof.name("properties.name")

	if ok {
		return name
	}

	return "a place with no name"
}

func (wof WOFFeature) name(path string) (string, bool) {

	return wof.StringValue(path)
}

func (wof WOFFeature) Placetype() string {

	pt, ok := wof.placetype("properties.wof:placetype")

	if ok {
		return pt
	}

	pt, ok = wof.placetype("properties.placetype")

	if ok {
		return pt
	}

	return "here be dragons"
}

func (wof WOFFeature) Deprecated() bool {

	path := "edtf:deprecated"

	d, ok := wof.StringProperty(path)

	if ok && d != "" && d != "u" && d != "uuuu" {
		return true
	}

	return false
}

func (wof WOFFeature) Superseded() bool {

	path := "edtf:superseded"

	d, ok := wof.StringProperty(path)

	if ok && d != "" && d != "u" && d != "uuuu" {
		return true
	}

     	body := wof.Body()

	pointers := body.Path("properties.wof:superseded_by").Data()

	if len(pointers.([]interface{})) != 0 {
		return true
	}

	return false
}

func (wof WOFFeature) placetype(path string) (string, bool) {

	return wof.StringValue(path)
}

func (wof WOFFeature) StringProperty(prop string) (string, bool) {

	path := fmt.Sprintf("properties.%s", prop)
	return wof.StringValue(path)
}

func (wof WOFFeature) StringValue(path string) (string, bool) {

	body := wof.Body()

	var value string
	var ok bool

	value, ok = body.Path(path).Data().(string)
	return value, ok
}

// sudo make me a package function and accept an interface
// (20151207/thisisaaronland)

func (wof WOFFeature) EnSpatialize() (*WOFSpatial, error) {

	id := wof.Id()
	name := wof.Name()
	placetype := wof.Placetype()
	deprecated := wof.Deprecated()
	superseded := wof.Superseded()

	body := wof.Body()

	var swlon float64
	var swlat float64
	var nelon float64
	var nelat float64

	children, _ := body.S("bbox").Children()

	if len(children) != 4 {
		return nil, errors.New("weird and freaky bounding box")
	}

	swlon = children[0].Data().(float64)
	swlat = children[1].Data().(float64)
	nelon = children[2].Data().(float64)
	nelat = children[3].Data().(float64)

	llat := nelat - swlat
	llon := nelon - swlon

	/*
		fmt.Printf("%f - %f = %f\n", nelat, swlat, llat)
		fmt.Printf("%f - %f = %f\n", nelon, swlon, llon)
	*/

	pt := rtreego.Point{swlon, swlat}
	rect, err := rtreego.NewRect(pt, []float64{llon, llat})

	if err != nil {
		return nil, err
	}

	return &WOFSpatial{rect, id, name, placetype, -1, deprecated, superseded}, nil
}

// sudo make me a package function and accept an interface
// (20151207/thisisaaronland)

func (wof WOFFeature) EnSpatializeGeom() ([]*WOFSpatial, error) {

	id := wof.Id()
	name := wof.Name()
	placetype := wof.Placetype()
	deprecated := wof.Deprecated()
	superseded := wof.Superseded()

	spatial := make([]*WOFSpatial, 0)
	polygons := wof.GeomToPolygons()

	for offset, poly := range polygons {

		swlat := 0.0
		swlon := 0.0
		nelat := 0.0
		nelon := 0.0

		for _, pt := range poly.OuterRing.Points() {

			lat := pt.Lat()
			lon := pt.Lng()

			if swlat == 0.0 || swlat > lat {
				swlat = lat
			}

			if swlon == 0.0 || swlon > lon {
				swlon = lon
			}

			if nelat == 0.0 || nelat < lat {
				nelat = lat
			}

			if nelon == 0.0 || nelon < lon {
				nelon = lon
			}

			// fmt.Println(lat, lon, swlat, swlon, nelat, nelon)
		}

		llat := nelat - swlat
		llon := nelon - swlon

		pt := rtreego.Point{swlon, swlat}
		rect, err := rtreego.NewRect(pt, []float64{llon, llat})

		if err != nil {
			return nil, err
		}

		sp := WOFSpatial{rect, id, name, placetype, offset, deprecated, superseded}
		spatial = append(spatial, &sp)
	}

	return spatial, nil
}

func (wof WOFFeature) Contains(latitude float64, longitude float64) bool {

	polygons := wof.GeomToPolygons()
	contains := false

	wg := new(sync.WaitGroup)

	for _, p := range polygons {

		wg.Add(1)

		go func(poly *WOFPolygon, lat float64, lon float64) {

			defer wg.Done()

			if poly.Contains(lat, lon) {
				contains = true
			}

		}(p, latitude, longitude)
	}

	wg.Wait()

	return contains
}

// sudo make me a package function and accept an interface... maybe?
// (20151207/thisisaaronland)}

func (wof WOFFeature) GeomToPolygons() []*WOFPolygon {

	body := wof.Body()

	var geom_type string

	geom_type, _ = body.Path("geometry.type").Data().(string)
	children, _ := body.S("geometry").ChildrenMap()

	polygons := make([]*WOFPolygon, 0)

	for key, child := range children {

		if key != "coordinates" {
			continue
		}

		var coordinates []interface{}
		coordinates, _ = child.Data().([]interface{})

		if geom_type == "Polygon" {
			polygons = append(polygons, wof.DumpPolygon(coordinates))
		} else if geom_type == "MultiPolygon" {
			polygons = wof.DumpMultiPolygon(coordinates)
		} else {
			// pass
		}
	}

	return polygons
}

// sudo these don't need to be public methods
// (20151207/thisisaaronland)

func (wof WOFFeature) DumpMultiPolygon(coordinates []interface{}) []*WOFPolygon {

	polygons := make([]*WOFPolygon, 0)

	for _, ipolys := range coordinates {

		polys := ipolys.([]interface{})

		polygon := wof.DumpPolygon(polys)
		polygons = append(polygons, polygon)
	}

	return polygons
}

// sudo these don't need to be public methods
// (20151207/thisisaaronland)

func (wof WOFFeature) DumpPolygon(coordinates []interface{}) *WOFPolygon {

	polygons := make([]geo.Polygon, 0)

	for _, ipoly := range coordinates {

		poly := ipoly.([]interface{})
		polygon := wof.DumpCoords(poly)
		polygons = append(polygons, polygon)
	}

	return &WOFPolygon{
		OuterRing:     polygons[0],
		InteriorRings: polygons[1:],
	}
}

// sudo these don't need to be public methods
// (20151207/thisisaaronland)

func (wof WOFFeature) DumpCoords(poly []interface{}) geo.Polygon {

	polygon := geo.Polygon{}

	for _, icoords := range poly {

		coords := icoords.([]interface{})

		lon := coords[0].(float64)
		lat := coords[1].(float64)

		pt := geo.NewPoint(lat, lon)
		polygon.Add(pt)
	}

	return polygon
}

// see below (20151207/thisisaaronland)

func UnmarshalFile(path string) (*WOFFeature, error) {

	body, read_err := ioutil.ReadFile(path)

	if read_err != nil {
		return nil, read_err
	}

	return UnmarshalFeature(body)
}

// this is disabled for now even though it (or something like it) is
// probably the new new; note the way we end up parsing the JSON body
// twice... we should not do that (20151207/thisisaaronland)

/*
func UnmarshalFileMulti(path string) ([]*WOFFeature, error) {

	body, read_err := ioutil.ReadFile(path)

	if read_err != nil {
		return nil, read_err
	}

	parsed, parse_err := gabs.ParseJSON(body)

	if parse_err != nil {
		return nil, parse_err
	}

	isa, ok := parsed.Path("type").Data().(string)

	if !ok {
		return nil, errors.New("failed to determine type")
	}

	features := make([]*WOFFeature, 0)

	if isa == "Feature" {

		f, err := UnmarshalFeature(body)

		if err != nil {
			return nil, err
		}

		features = append(features, f)

	} else if isa == "FeatureCollection" {

		collection, err := UnmarshalFeatureCollection(body)

		if err != nil {
			return nil, err
		}

		features = collection
	} else {
		return nil, errors.New("unknown type")
	}

	return features, nil
}
*/

// see above inre passing bytes or an already parsed thing-y
// (20151207/thisisaaronland)

func UnmarshalFeatureCollection(raw []byte) ([]*WOFFeature, error) {

	parsed, parse_err := gabs.ParseJSON(raw)

	if parse_err != nil {
		return nil, parse_err
	}

	children, _ := parsed.S("features").Children()

	collection := make([]*WOFFeature, 0)

	for _, child := range children {

		f := WOFFeature{
			Parsed: child,
		}

		collection = append(collection, &f)
	}

	return collection, nil
}

// see above inre passing bytes or an already parsed thing-y
// (20151207/thisisaaronland)

func UnmarshalFeature(raw []byte) (*WOFFeature, error) {

	parsed, parse_err := gabs.ParseJSON(raw)

	if parse_err != nil {
		return nil, parse_err
	}

	rsp := WOFFeature{
		Parsed: parsed,
	}

	return &rsp, nil
}
