// test program to parse dwml xml format from NOAA API

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime/pprof"
	"time"
)

var inputFile = flag.String("infile", "1~.xml", "Input file path")

func dump_httpresp(resp *http.Response) {
	fmt.Println("resp: ", resp)

	st := reflect.TypeOf(resp)
	fmt.Println(st)

	val := reflect.ValueOf(resp).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		typeName := valueField.Type().Name()

		fmt.Printf("Field Name: %s(%s),\t\t\t Field Value: %v\n", typeField.Name, typeName, valueField.Interface())
	}

	fmt.Println("=== headers ===")
	for k, v := range resp.Header {
		fmt.Printf("%20s = %20s\n", k, v)
	}
}

func dump_value(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		typeName := valueField.Type().Name()
		fmt.Printf("%s(%s): %v\n", typeField.Name, typeName, valueField.Interface())
	}
}

// ===========================================================================
// NOAA API dwml structures (top to bottom order)

// root element

type Dwml struct {
	XMLName xml.Name `xml:"dwml"`
	Header  Header   `xml:"head"`
	Data    Data     `xml:"data"`
}

// 3 second level elements

type Header struct {
	XMLName xml.Name `xml:"head"`
	Product Product  `xml:"product"`
}

type Data struct {
	XMLName     xml.Name     `xml:"data"`
	TimeLayouts []TimeLayout `xml:"time-layout"`
	Parameters  Parameters   `xml:"parameters"`
}

// head element children

type Product struct {
	XMLName xml.Name `xml:"product"`
	Src     string   `xml:"srsName,attr"`
	Name    string   `xml:"concise-name,attr"`
	Mode    string   `xml:"operational-mode,attr"`
}

// data element children

type TimeLayout struct {
	XMLName       xml.Name `xml:"time-layout"`
	Coordinate    string   `xml:"time-coordinate,attr"`
	Summarization string   `xml:"summarization,attr"`
	Key           string   `xml:"layout-key"`
	StartTime     []string `xml:"start-valid-time"`
	EndTime       []string `xml:"end-valid-time"`
}

type Parameters struct {
	XMLName       xml.Name  `xml:"parameters"`
	Temperature   []Valueset `xml:"temperature"`
	WindSpeed     []Valueset `xml:"wind-speed"`
	Direction     []Valueset `xml:"direction"`
	CloudAmount   []Valueset `xml:"cloud-amount"`
	Precipitation []Valueset `xml:"precipitation"`
	Humidity      []Valueset `xml:"humidity"`
}

type Visibility struct {
	XMLName       xml.Name
	Units string              `xml:"units,attr"`
}

type ConditionValue struct {
	Coverage string              `xml:"coverage,attr"`
	Intencity string              `xml:"intencity,attr"`
	WeatherType string              `xml:"weather-type,attr"`
	Qualifier string              `xml:"qualifier,attr"`
}

type WeatherConditions struct {
	Value []ConditionValue `xml:"value"`
}

type Weather struct {
	TimeLayout string              `xml:"time-layout,attr"`
	Conditions []WeatherConditions `xml:"weather-conditions"`
}

type Valueset struct {
	Type       string   `xml:"type,attr"`
	Units      string   `xml:"units,attr"`
	TimeLayout string   `xml:"time-layout,attr"`
	Name       string   `xml:"name"`
	Values     []string `xml:"value"`
}

// ===========================================================================

func decode_dwml(xmlFile *os.File) {
	decoder := xml.NewDecoder(xmlFile)

	p := &Dwml{}
	if err := decoder.Decode(p); err != nil {
		fmt.Println("ERROR: %v", err)
		return
	}

	fmt.Println("Time layouts:")

	for _, v := range p.Data.TimeLayouts {
		fmt.Println("    ", v.Key, len(v.StartTime))
	}

	fmt.Println("Air:")

	for _, v := range p.Data.Parameters.Temperature {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}

	for _, v := range p.Data.Parameters.Humidity {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}

	fmt.Println("Wind:")

	for _, v := range p.Data.Parameters.WindSpeed {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}

	for _, v := range p.Data.Parameters.Direction {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}

	fmt.Println("Sky:")

	for _, v := range p.Data.Parameters.CloudAmount {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}

	for _, v := range p.Data.Parameters.Precipitation {
		fmt.Println("    ", v.Name, v.TimeLayout, len(v.Values))
	}
}

func file_cached(fname *string) bool {

	fi, err := os.Stat(*fname)

	if err != nil || fi == nil {
		return false
	}

	if fi.ModTime().Unix()-time.Now().UTC().Unix() < 60*72 {
		return true
	}

	newname := *fname + ".old"

	err = os.Rename(*fname, newname)
	if err != nil {
		fmt.Println("ERROR: %v", err)
	}

	return false
}

func main() {

	flag.Parse()

	dump_resp := false

	if !file_cached(inputFile) {

		fmt.Println("loading from NOAA")

		url := "http://www.weather.gov/forecasts/xml/SOAP_server/ndfdXMLclient.php?whichClient=NDFDgen&zipCodeList=10001&product=time-series&maxt=maxt&mint=mint&temp=temp&wspd=wspd&wdir=wdir&wx=wx&rh=rh&snow=snow&wwa=wwa&sky=sky&appt=appt&Submit=Submit"

		resp, err := http.Get(url)

		fmt.Println("err: ", err)

		if dump_resp {
			dump_httpresp(resp)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println(body)

		err = ioutil.WriteFile(*inputFile, body, 0666)
		fmt.Println("write: ", err)
	}

	xmlFile, err := os.Open(*inputFile)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer xmlFile.Close()

	f, err := os.Create("1~.prof")
	if err != nil {
		fmt.Println(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	decode_dwml(xmlFile)
}
