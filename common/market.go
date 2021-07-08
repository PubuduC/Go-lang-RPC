package common

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/subchen/go-trylock"
	"log"
	"os"
	"time"
)

// Vegetable struct represent a vegetable
type Vegetable struct {
	Name string `csv:"Name"`
	PricePerKg float32 `csv:"PricePerKg"`
	AvailableAmountOfKg float32 `csv:"AvailableAmountOfKg"`
}

// Market struct represent a market
type Market struct {
	database []*Vegetable
}

var mu = trylock.New()

func writeCsvFile(newVegs []*Vegetable) []*Vegetable {
	if ok := mu.TryLockTimeout(3 * time.Second); !ok {
		return nil
	}
	err := os.Remove("common/data.csv")

	if err != nil {
		fmt.Println(err)
		return nil
	}
	f, err := os.OpenFile("common/data.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to open data file ", err)
	}
	writer := csv.NewWriter(f)
	writer.Write([]string{"Name", "PricePerKg", "AvailableAmountOfKg"})

	for _, veg := range newVegs {
		row := []string{veg.Name, fmt.Sprintf("%.2f", veg.PricePerKg), fmt.Sprintf("%.2f", veg.AvailableAmountOfKg)}
		_ = writer.Write(row)
	}
	writer.Flush()

	f.Close()
	mu.Unlock()
	return newVegs
}

func readCsvFile(filePath string) []*Vegetable {
	if ok := mu.RTryLockTimeout(1 * time.Second); !ok {
		return nil
	}
	defer mu.RUnlock()
	f, err := os.OpenFile("common/data.csv", os.O_CREATE|os.O_RDONLY, 0644)

	//f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file " + filePath, err)
	}
	defer f.Close()

	var vegetables []*Vegetable

	if err := gocsv.UnmarshalFile(f, &vegetables); err != nil {
		log.Print(err)
		return nil
	}

	return vegetables
}

// NewMarket function returns a new instance of Market (pointer).
func NewMarket() *Market {
	return &Market{
		database : readCsvFile("common/data.csv"),
	}
}

// Get method gets a vegetable by name(procedure)
func (market *Market) Get(payload string, reply *Vegetable) error {

	for _, v := range market.database {
		if v.Name == payload {
			*reply = *v
			return nil
		}
	}

	return fmt.Errorf("vegetable with name '%s' not exists", payload)

}

// GetAvailableAmount method gets a vegetable by name(procedure)
func (market *Market) GetAvailableAmount(payload string, reply *float32) error {

	for _, v := range market.database {
		if v.Name == payload {
			*reply = v.AvailableAmountOfKg
			return nil
		}
	}

	return fmt.Errorf("vegetable with name '%s' not exists", payload)

}

// GetPricePerKg method gets a vegetable by name(procedure)
func (market *Market) GetPricePerKg(payload string, reply *float32) error {

	for _, v := range market.database {
		if v.Name == payload {
			*reply = v.PricePerKg
			return nil
		}
	}

	return fmt.Errorf("vegetable with name '%s' not exists", payload)

}

// GetAll method gets all vegetables(procedure)
func (market *Market) GetAll(payLoad string, reply *string) error {

	vegNames := ""

	if len(market.database) == 0 {
		return fmt.Errorf("no vegetables available")
	}
	for _, vegetable := range market.database {
		vegNames += vegetable.Name +", "
	}
	*reply = vegNames
	return nil

}

// Add method adds a vegetables to the market (procedure)
func (market *Market) Add(payload Vegetable, reply *Vegetable) error {

	// check if vegetable already exists in the database
	if _, ok := vegetableAlreadyExists(payload.Name, market.database); ok {
		return fmt.Errorf("vegetable with name '%s' already exists", payload.Name)
	}

	// add vegetable to the database
	vegetable := &Vegetable{
		Name: payload.Name,
		PricePerKg: payload.PricePerKg,
		AvailableAmountOfKg: payload.AvailableAmountOfKg,
	}
	market.database = append(market.database, vegetable)
	appendVegetableToCsv(vegetable)
	// set reply value
	*reply = payload

	// return `nil` error
	return nil
}

// Update method updates a vegetables to the market (procedure)
func (market *Market) Update(payload Vegetable, reply *Vegetable) error {

	// check if vegetable already exists in the database
	alreadyExist := false

	for _, v := range market.database {
		if v.Name == payload.Name {
			alreadyExist = true
			v.PricePerKg = payload.PricePerKg
			v.AvailableAmountOfKg = payload.AvailableAmountOfKg
		}
	}
	// if vegetable not exists in the database
	if !alreadyExist {
		vegetable := &Vegetable{
			Name: payload.Name,
			PricePerKg: payload.PricePerKg,
			AvailableAmountOfKg: payload.AvailableAmountOfKg,
		}
		market.database = append(market.database, vegetable)
		appendVegetableToCsv(vegetable)
		// set reply value
		*reply = payload

		// return `nil` error
		return nil
	}

	// update csv file if already exist
	writeCsvFile(market.database)

	// set reply value
	*reply = payload

	// return `nil` error
	return nil
}

func vegetableAlreadyExists(vegName string, list []*Vegetable) (bool, bool) {
	for _, b := range list {
		if b.Name == vegName {
			return true, true
		}
	}
	return false, false
}

func appendVegetableToCsv(vegetable *Vegetable)  {
	f, err := os.OpenFile("common/data.csv", os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	writer := csv.NewWriter(f)

	if len(records)==0 {
		writer.Write([]string{"Name", "PricePerKg", "AvailableAmountOfKg"})
	}

	var row []string
	row = append(row, vegetable.Name)
	row = append(row, fmt.Sprintf("%f", vegetable.PricePerKg))
	row = append(row, fmt.Sprintf("%f", vegetable.AvailableAmountOfKg))
	writer.Write(row)
	writer.Flush()
	defer f.Close()
}




