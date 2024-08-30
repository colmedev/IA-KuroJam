package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Original struct for reading JSON
type OriginalCareer struct {
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Tasks         []string            `json:"tasks"`
	Knowledge     map[string]Category `json:"knowledge"`
	Capacities    map[string]Category `json:"capacities"`
	Skills        map[string]Category `json:"skills"`
	Technology    map[string]Category `json:"technology"`
	Personality   Personality         `json:"personality"`
	Education     string              `json:"education"`
	AverageSalary string              `json:"averageSalary"`
	LowerSalary   string              `json:"lowerSalary"`
	HighestSalary string              `json:"highestSalary"`
}

type Category struct {
	Name  string   `json:"name"`
	Areas []string `json:"areas"`
}

type Personality struct {
	Description string   `json:"description"`
	Attributes  []string `json:"attributes"`
}

// Final struct for writing JSON
type FinalCareer struct {
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Tasks         []string    `json:"tasks"`
	Knowledge     []Category  `json:"knowledge"`
	Capacities    []Category  `json:"capacities"`
	Skills        []Category  `json:"skills"`
	Technology    []Category  `json:"technology"`
	Personality   Personality `json:"personality"`
	Education     string      `json:"education"`
	AverageSalary string      `json:"averageSalary"`
	LowerSalary   string      `json:"lowerSalary"`
	HighestSalary string      `json:"highestSalary"`
}

func main() {
	// Open the file for reading
	file, err := os.Open("./careers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var finalCareers []FinalCareer

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()

		// Unmarshal the original JSON into the original struct
		var original OriginalCareer
		err := json.Unmarshal([]byte(line), &original)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		// Convert to the final structure
		finalCareers = append(finalCareers, FinalCareer{
			Title:         original.Title,
			Description:   original.Description,
			Tasks:         original.Tasks,
			Knowledge:     convertMapToArray(original.Knowledge),
			Capacities:    convertMapToArray(original.Capacities),
			Skills:        convertMapToArray(original.Skills),
			Technology:    convertMapToArray(original.Technology),
			Personality:   original.Personality,
			Education:     original.Education,
			AverageSalary: original.AverageSalary,
			LowerSalary:   original.LowerSalary,
			HighestSalary: original.HighestSalary,
		})
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Open a new file for writing the processed JSON data
	outputFile, err := os.Create("processed_careers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Write each final career struct as a JSON object on a new line
	for _, c := range finalCareers {
		data, err := json.Marshal(c)
		if err != nil {
			log.Fatal(err)
		}

		_, err = outputFile.WriteString(string(data) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Successfully converted and saved the JSON data.")
}

// Helper function to convert map to an array of categories
func convertMapToArray(inputMap map[string]Category) []Category {
	var outputArray []Category
	for _, value := range inputMap {
		outputArray = append(outputArray, value)
	}
	return outputArray
}
