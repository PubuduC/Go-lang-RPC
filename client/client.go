package main

import (
	"DSLab1-209319K/common"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

func main() {

	// get RPC client by dialing at `rpc.DefaultRPCPath` endpoint
	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:9000") // or `localhost:9000`

	/*--------------*/

	//create veg variable of type `common.Vegetable`
	var veg common.Vegetable
	var vegNames string
	var vegAvailableAmount float32
	var vegPricePerKg float32

	// client server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello this is client server!")
	})
	// Add vegetable POST method
	mux.HandleFunc("/vegetables/add", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			fmt.Fprint(writer, "Unable to parse form: ", err)
			log.Print("Unable to parse form: ", err)
			return
		}

		pricePerKgValue, err := strconv.ParseFloat(request.FormValue("PricePerKg"), 32)
		if err != nil {
			fmt.Fprint(writer, "Unable to read input PricePerKg: Please enter valid value (float/int)")
			log.Print("Unable to read input PricePerKg: ", err)
			return
		}
		pricePerKgFloat32 := float32(pricePerKgValue)

		availableAmountOfKgValue, err := strconv.ParseFloat(request.FormValue("AvailableAmountOfKg"), 32)
		if err != nil {
			fmt.Fprint(writer, "Unable to read input AvailableAmountOfKg: Please enter valid value (float/int)")
			log.Print("Unable to read input AvailableAmountOfKg: ", err)
			return
		}
		availableAmountOfKgFloat32 := float32(availableAmountOfKgValue)

		// add vegetable
		if err := client.Call("Market.Add", common.Vegetable{
			Name: request.FormValue("Name"),
			PricePerKg:  pricePerKgFloat32,
			AvailableAmountOfKg: availableAmountOfKgFloat32,
		}, &veg); err != nil {
			fmt.Println(err)
			fmt.Fprint(writer, err)
			return
		} else {
			fmt.Printf("Vegetable '%s' created \n", veg.Name)
			fmt.Fprintf(writer, "Vegetable '%s' created \n", veg.Name)
		}
	})

	mux.HandleFunc("/vegetables/update", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			fmt.Fprint(writer, "Unable to parse form: ", err)
			log.Print("Unable to parse form: ", err)
			return
		}

		pricePerKgValue, err := strconv.ParseFloat(request.FormValue("PricePerKg"), 32)
		if err != nil {
			fmt.Fprint(writer, "Unable to read input PricePerKg: Please enter valid value (float/int)")
			log.Print("Unable to read input PricePerKg: ", err)
			return
		}
		pricePerKgFloat32 := float32(pricePerKgValue)

		availableAmountOfKgValue, err := strconv.ParseFloat(request.FormValue("AvailableAmountOfKg"), 32)
		if err != nil {
			fmt.Fprint(writer, "Unable to read input AvailableAmountOfKg: Please enter valid value (float/int)")
			log.Print("Unable to read input AvailableAmountOfKg: ", err)
			return
		}
		availableAmountOfKgFloat32 := float32(availableAmountOfKgValue)

		// add vegetable
		if err := client.Call("Market.Update", common.Vegetable{
			Name: request.FormValue("Name"),
			PricePerKg:  pricePerKgFloat32,
			AvailableAmountOfKg: availableAmountOfKgFloat32,
		}, &veg); err != nil {
			fmt.Println(err)
			fmt.Fprint(writer, err)
			return
		} else {
			fmt.Printf("Vegetable '%s' updated \n", veg.Name)
			fmt.Fprintf(writer, "Vegetable '%s' updated \n", veg.Name)
			fmt.Fprintf(writer, "Available amount in kg: %f \n", veg.AvailableAmountOfKg)
			fmt.Fprintf(writer, "Price per kg: %f \n", veg.PricePerKg)
		}
	})


	// Get vegetable GET method
	mux.HandleFunc("/vegetables/get", func(writer http.ResponseWriter, request *http.Request) {
		if err := client.Call("Market.Get", request.FormValue("Name"), &veg); err != nil {
			fmt.Fprint(writer, err)
		} else {
			fmt.Fprintf(writer, "Vegetable '%s' found \n", veg.Name)
			fmt.Fprintf(writer, "Vegetable '%.3f' kg found \n", veg.AvailableAmountOfKg)
			fmt.Fprintf(writer, "Price per kg is '%.2f' \n", veg.PricePerKg)
		}
	})

	// Get vegetable pricePerKg GET method
	mux.HandleFunc("/vegetables/get/pricePerKg", func(writer http.ResponseWriter, request *http.Request) {
		if err := client.Call("Market.GetPricePerKg", request.FormValue("Name"), &vegPricePerKg); err != nil {
			fmt.Fprint(writer, err)
		} else {
			fmt.Fprintf(writer, "Vegetable '%s' found \n", request.FormValue("Name"))
			fmt.Fprintf(writer, "Price per kg is '%.2f' \n", vegPricePerKg)
		}
	})

	// Get vegetable pricePerKg GET method
	mux.HandleFunc("/vegetables/get/availableAmount", func(writer http.ResponseWriter, request *http.Request) {
		if err := client.Call("Market.GetAvailableAmount", request.FormValue("Name"), &vegAvailableAmount); err != nil {
			fmt.Fprint(writer, err)
		} else {
			fmt.Fprintf(writer, "Vegetable '%s' found \n", request.FormValue("Name"))
			fmt.Fprintf(writer, "Price per kg is '%.3f' \n", vegAvailableAmount)
		}
	})

	// Get all vegetables GET method
	mux.HandleFunc("/vegetables/get/all", func(writer http.ResponseWriter, request *http.Request) {
		if err := client.Call("Market.GetAll", "just a string", &vegNames); err != nil {
			fmt.Fprint(writer, err)
		} else {
			fmt.Fprintf(writer, "%d No. of Vegetables found \n", len(strings.Split(vegNames, ",")))
			fmt.Fprint(writer, vegNames)
		}
	})

	http.ListenAndServe(":7000", mux)

}

