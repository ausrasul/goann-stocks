package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strconv"
  "bytes"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, line := range lines {
    fmt.Fprintln(w, line)
  }
  return w.Flush()
}


func normalize(min float32, max float32, val string) (normal string){
	size := max - min
	var normalF float32
	valF, _ := strconv.ParseFloat(val, 32)
	normalF = ((float32(valF) - min) * (2/size) -1);
	normal = strconv.FormatFloat(float64(normalF), 'f', -1, 32)
	return
}

func denormalize(min float32, max float32, val float32) (denormal string){
	size := max - min
	var denormalF float32
	//valF, _ := strconv.ParseFloat(val, 32)
	valF := val
	denormalF = (1 + float32(valF)) * (size/2) + min
	//normalF = ((float32(valF) - min) * (2/size) -1)
	denormal = strconv.FormatFloat(float64(denormalF), 'f', -1, 32)
	return
}

func createStockData(conf Config){
	rowFile := conf.SqlRowFile
	outFile := conf.TrainOutFile
	layers := conf.Layers
	inputNumber := layers[0]
	outputNumber := layers[len(layers)-1]
	minPrice := conf.MinPrice
	maxPrice := conf.MaxPrice
	
	// read the row file
	lines, err := readLines(rowFile)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	//inDataRows := len(lines)
	//for inDataRows > (inputNumber + outputNumber){
	// parse line by line and print to out file
	count := uint32(0)
	totalTrainData := uint32(0)
	var inString bytes.Buffer
	var outString bytes.Buffer
	var buffer []string
	for _, line := range lines{
		normal := normalize(minPrice, maxPrice, line)
		if count < inputNumber {
			inString.WriteString(normal)
			inString.WriteString(" ")
			count++
		} else if count == inputNumber {
			outString.WriteString(normal)
			buffer = append(buffer, inString.String())
			buffer = append(buffer, outString.String())
			inString.Reset()
			outString.Reset()
			count = 0
			totalTrainData++
		}
	}
	var header []string
	var headerB bytes.Buffer
	totalTrainDataS := strconv.FormatInt(int64(totalTrainData), 10)
	inputNumberS := strconv.FormatInt(int64(inputNumber), 10)
	outputNumberS := strconv.FormatInt(int64(outputNumber), 10)
	headerB.WriteString(totalTrainDataS)
	headerB.WriteString(" ")
	headerB.WriteString(inputNumberS)
	headerB.WriteString(" ")
	headerB.WriteString(outputNumberS)
	header = append(header, headerB.String())
	header = append(header, buffer...)
	err = writeLines(header, outFile)
	if err != nil {
		log.Fatal("cannot write to out file")
	}
}