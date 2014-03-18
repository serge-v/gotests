package main

import (
	"fmt"
	"net/http"
	"reflect"
)

func dump_httpresp(resp *http.Response) {
	fmt.Println("resp: ", resp)

	st := reflect.TypeOf(resp)
	fmt.Println(st)

	val := reflect.ValueOf(resp).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		typeName := valueField.Type().Name()
 
		fmt.Printf("%s(%s),\t\t\t Field Value: %v\n", typeField.Name, typeName, valueField.Interface())
	}
	
	fmt.Println("=== headers ===")
	for k, v := range resp.Header {
		fmt.Printf("%20s = %20s\n", k, v)
	}
}

func main() {

	for i := 0; i < 1000; i++ {
		resp, err := http.Get("http://localhost:4001/sadasdasd/")
		if err != nil {
			fmt.Println(err)
			break
		}

		if i == 999 {
			dump_httpresp(resp)
		}
	}

	fmt.Println("1000 reqs done")

}
