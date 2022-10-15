package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Dogancan94/app/helpers"
)

//Varaibles
func main() {
	fmt.Println("Hello, world.")

	var whatToSay string
	var age int

	whatToSay = "Goodbye, cruel world."
	helpers.InfoLog(whatToSay)

	age = 7
	fmt.Println("Age is ", age)

	word, otherWord := saySomething()
	fmt.Println("Said: "+word, otherWord)

	showMeSomePointer()

	person := Person{
		Firsname:    "Dogancan",
		Lastname:    "Kinik",
		PhoneNumber: "555 555 1212",
	}
	log.Println("User is:", person.Firsname, person.Lastname, person.PhoneNumber)

	createMap()
	createSlice()
	createDecisionStructure()
	createLoop()
	createInterface()
	createChannel()
	createFromJson()

	result, err := divide(100, -10)
	if err != nil {
		helpers.InfoLog(err)
		return
	}
	helpers.InfoLog(result)

}

func saySomething() (string, string) {
	return "something", "else"
}

//Pointers
func showMeSomePointer() {
	color := "Green"
	log.Println("string is set to", color)

	changeByPointer(&color)
	log.Println("string is set to", color)
}

func changeByPointer(s *string) {
	newColor := "Red"
	*s = newColor
}

//Types & Structs
type Person struct {
	Firsname    string
	Lastname    string
	PhoneNumber string
	Age         int
	BirthDate   time.Time
}

type Address struct {
	Address string
}

func (m *Address) printAddress() string {
	return m.Address
}

//Maps
type Animal struct {
	Age     int
	Type    string
	Species string
}

func createMap() {
	firstNames := make(map[string]string)
	firstNames["mom"] = "Damla"
	firstNames["parrot"] = "Alex"
	log.Println(firstNames)

	notes := make(map[string]float64)
	notes["Math"] = 90
	notes["Science"] = 80.5
	log.Println(notes)

	animals := make(map[string]Animal)
	animals["Alex"] = Animal{Age: 14, Type: "Parrot", Species: "Jako"}
	log.Println(animals)
}

//Slices
func createSlice() {
	var animals []string
	animals = append(animals, "Dog")
	animals = append(animals, "Cat")
	animals = append(animals, "Parrot")
	animals = append(animals, "Fish")
	log.Println(animals)

	var notes []int
	notes = append(notes, 100)
	notes = append(notes, 90)
	notes = append(notes, 80)
	sort.Ints(notes)
	log.Println(notes)

	var pets []Animal
	pets = append(pets, Animal{Age: 14, Type: "Parrot", Species: "Jako"})
	pets = append(pets, Animal{Age: 6, Type: "Cat", Species: "Siam"})
	pets = append(pets, Animal{Age: 3, Type: "Dog", Species: "Golden"})
	sort.Slice(pets, func(i, j int) bool {
		return pets[i].Age < pets[j].Age
	})

	log.Println(pets)
	log.Println(pets[0:2])
}

//Decision Structure
func createDecisionStructure() {
	isValid := false
	if isValid {
		log.Println("It is true")
	} else {
		log.Println("It is false")
	}

	animal := "parrot"
	if animal == "cat" {
		log.Println("I am Cat!!!!")
	}

	switch animal {
	case "cat":
		log.Println("I am Cat!!!!")
	case "dog":
		log.Println("I am Dog!!!!")
	case "parrot":
		log.Println("I am Parrot!!!!")
	case "fish":
		log.Println("I am Fish!!!!")
	default:
		log.Println("I am nothing!!!!!")
	}
}

//Loops
func createLoop() {
	for i := 0; i <= 5; i++ {
		log.Println(i)
	}

	humans := []string{"Dogancan", "Erdem", "Damla", "Balki"}
	for _, human := range humans {
		log.Println(human)
	}

	animals := make(map[string]string)
	animals["dog"] = "Mars"
	animals["cat"] = "Aristo"
	animals["parrot"] = "Alex"
	animals["budgie"] = "Boncuk"

	for animalType, animal := range animals {
		log.Println(animalType + " name is " + animal)
	}
}

//interface
type Mammals interface {
	Says() string
	NumberOfLegs() int
}

type Dog struct {
	Name  string
	Breed string
}

func (d *Dog) Says() string {
	return "Woof"
}

func (d *Dog) NumberOfLegs() int {
	return 4
}

type Gorilla struct {
	Name          string
	Color         string
	NumberOfTeeth int
}

func (g *Gorilla) Says() string {
	return "Hello"
}

func (g *Gorilla) NumberOfLegs() int {
	return 2
}

func createInterface() {
	dog := Dog{
		Name:  "Mars",
		Breed: "Golden Retriever",
	}

	gorilla := Gorilla{
		Name:          "Charlie",
		Color:         "Brown",
		NumberOfTeeth: 32,
	}

	PrintInfo(&dog)
	PrintInfo(&gorilla)
}

func PrintInfo(mammal Mammals) {
	log.Println("Mammal says", mammal.Says(), "and has", mammal.NumberOfLegs(), "legs")
}

//Channels
func createChannel() {
	intChan := make(chan int)
	defer close(intChan)
	go calculateValue(intChan)
	num := <-intChan
	helpers.InfoLog(num)
}

const numPool = 1000

func calculateValue(intChan chan int) {
	number := helpers.CreateRandomNumber(numPool)
	intChan <- number
}

//JSON
type Actor struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Haircolor string `json:"hair_color"`
	Hasdog    bool   `json:"has_dog"`
}

func createFromJson() {
	myJson := `[
					{
						"first_name": "Clark",
						"last_name": "Kent",
						"hair_color": "black",
						"has_dog": true
					},
					{
						"first_name": "Bruce",
						"last_name": "Wayne",
						"hair_color": "black",
						"has_dog": false
					}
				]`

	var actors []Actor
	err := json.Unmarshal([]byte(myJson), &actors)
	if err != nil {
		helpers.InfoLog(err)
	}

	//write json from struct
	var actorSlice []Actor

	var actor Actor
	actor.Firstname = "Keanu"
	actor.Lastname = "Reaves"
	actor.Haircolor = "black"
	actor.Hasdog = true
	actorSlice = append(actorSlice, actor)

	newJson, err := json.MarshalIndent(actorSlice, "", "  ")
	if err != nil {
		log.Println("Error Marshalling")
	}
	fmt.Println(string(newJson))

	helpers.InfoLog(actors)
}

//TESTING

func divide(x, y float32) (float32, error) {
	if y == 0 {
		return 0, errors.New("divider can't be 0")
	}
	result := x / y
	return result, nil
}
