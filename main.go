package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const filePath string = "dat.txt"

// const filePath string = "dat2.txt"

// const filePath string = "dat3.txt"

const resultsFilePath string = "rez.txt"

func main() {
	data := Read(filePath)
	n := len(data)
	mainchannel := make(chan Student)
	datachannel := make(chan Student)
	filteredchannel := make(chan Student)
	resultschannel := make(chan []Student)
	writeFlag := make(chan int)
	readFlag := make(chan int)

	threadCount := 6
	for i := 0; i < threadCount; i++ {
		go WorkProcess(datachannel, readFlag, filteredchannel)
	}
	go DataProcess(n, mainchannel, writeFlag, datachannel, readFlag)
	go ResultProcess(n, filteredchannel, resultschannel)
	go func() {
		for _, item := range data {
			writeFlag <- 1
			mainchannel <- item
		}
	}()
	results := <-resultschannel
	WriteResultsToFile(data, results, resultsFilePath)

}

func WorkProcess(dataChannel <-chan Student, readSignal chan<- int, filteredChannel chan Student) {
	for {
		readSignal <- 1
		student := <-dataChannel
		if student.Year == -1 {
			break
		}
		if student.getCalculation() >= 50.0 {
			filteredChannel <- student
		}
	}
	filteredChannel <- Student{Year: -1}
}

func DataProcess(n int, mainchannel <-chan Student, readFlag <-chan int, dataChannel chan<- Student, writeFlag <-chan int) {
	data := make([]Student, n/2)
	index := 0
	for {
		if index > len(data) {
			<-writeFlag
			index--
			if data[index].Year == -1 {
				<-writeFlag
				dataChannel <- Student{Year: -1}
				return
			}
			dataChannel <- data[index]
		} else if index == 0 {
			<-readFlag
			data[index] = <-mainchannel
			index++
		} else {
			select {
			case <-readFlag:
				data[index] = <-mainchannel
				index++
			case <-writeFlag:
				index--
				if data[index].Year == -1 {
					<-writeFlag
					dataChannel <- Student{Year: -1}
					return
				}
				dataChannel <- data[index]
			}
		}
	}
}

func ResultProcess(size int, filteredchannel <-chan Student, resultsChannel chan<- []Student) {
	data := make([]Student, size)
	index := 0
	for {
		item := <-filteredchannel
		if item.Year == -1 {
			resultContainer := make([]Student, index)
			for i := 0; i < index; i++ {
				resultContainer[i] = data[i]
			}
			resultsChannel <- resultContainer
			return
		} else {
			index++
			if index == 1 {
				data[0] = item
			} else if item.Name > data[index-2].Name {
				data[index-1] = item
			} else {
				for i := 0; i < index-1; i++ {
					if item.Name < data[i].Name {
						for u := index - 1; u > i; u-- {
							data[u] = data[u-1]
						}
						data[i] = item
						break
					}
				}
			}
		}
	}
}

func Read(fileName string) []Student {
	var studs []Student
	file, e := os.Open(fileName)
	if e != nil {
		fmt.Println("Error is = ", e.Error())
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		var stud Student
		stud.Name = s[0]
		val, _ := strconv.Atoi(s[1])
		stud.Year = int(val)
		value, _ := strconv.ParseFloat(s[2], 32)
		stud.Grade = float32(value)
		studs = append(studs, stud)
	}
	file.Close()
	studs = append(studs, Student{Year: -1})
	return studs
}

func WriteResultsToFile(originalData []Student, result []Student, path string) {
	os.Remove(path)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed writing to file %s", err)
	}

	dataWriter := bufio.NewWriter(file)
	// PRINT INITIAL DATA
	_, _ = dataWriter.WriteString(fmt.Sprintf("%45v\n", "Initial"))
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")
	_, _ = dataWriter.WriteString(fmt.Sprintf("|%20v|%20v|%20v|%20v|\n", "Name", "Year", "Grade", "Percentage"))
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")
	for _, Student := range originalData {
		if Student.Name != " " && Student.Year != 0 {
			_, _ = dataWriter.WriteString(Student.toString() + "\n")
		}
	}
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")

	// PRINT RESULTS
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	_, _ = dataWriter.WriteString(fmt.Sprintf("\n\n%45v\n", "Result"))
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")
	_, _ = dataWriter.WriteString(fmt.Sprintf("|%20v|%20v|%20v|%20v|\n", "Name", "Year", "Grade", "Percentage"))
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")
	for _, Student := range result {
		if Student.Name != " " && Student.Year != 0 {
			_, _ = dataWriter.WriteString(Student.toString() + "\n")
		}
	}
	_, _ = dataWriter.WriteString(strings.Repeat("-", 85) + "\n")

	_ = dataWriter.Flush()
	_ = file.Close()
}

type Student struct {
	Name  string
	Year  int
	Grade float32
}

func (c *Student) toString() string {
	return fmt.Sprintf("|%20v|%20v|%20v|%20.1f|", c.Name, c.Year, c.Grade, c.getCalculation())
}

func (c *Student) getCalculation() float32 {
	return c.Grade * 10
}
