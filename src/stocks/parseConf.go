package main

import "encoding/json"
import "os"
import "github.com/white-pony/go-fann"
import "flag"
import "log"

type Config struct {
	DataFile string
    NetFile string
    Layers []uint32
	MaxEpochs uint32
	EpochsBetweenReports uint32
    DesiredError float32
	SqlRowFile string
	TrainOutFile string
	MaxPrice float32
	MinPrice float32
	Training  bool
	CreateTrain bool
}

type JsonValues struct {
	Values []FannValues
}
type FannValues struct {
	Date string
	Close float32
	High float32
	Low float32
	Inputs []fann.FannType
	Outputs []fann.FannType
}

func getConf(fileName string) (conf Config, values JsonValues){
	var confFileName string
	var inputFileName string
	var Training bool
	var CreateTrain bool
	if len(fileName) == 0 {
		confFileNamePtr := flag.String("c", "", "absolute conf file name")
		inputFileNamePtr := flag.String("i", "", "absolute input file name")
		TrainingPtr := flag.Bool("t", false, "true if training, false if test")
		CreateTrainPtr := flag.Bool("g", false, "true if you want to create training data, false if not")
		flag.Parse()
		confFileName = *confFileNamePtr
		inputFileName = *inputFileNamePtr
		CreateTrain = *CreateTrainPtr
		Training = *TrainingPtr
		if confFileName == "" {
			log.Fatal("enter conf file name with absolute path. use -c param")
		}
		if Training != true && Training != false {
			log.Fatal("enter run type, training true or false. use -t param")
		}
		if Training == false && inputFileName == "" && CreateTrain == false{
			log.Fatal("This is a value test, specify an input file. use -i param or switch to train with -t or generate train -g")
		}
		if CreateTrain != true && CreateTrain != false {
			log.Fatal("enter run type, Generate training data. use -g param")
		}
	} else {
		confFileName = fileName
	}
	file, err := os.Open(confFileName)
	if err != nil {
		log.Fatal("Error opening config file", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Error parsing config file:", err)
	}
	if Training == false && CreateTrain == false{
		file, err = os.Open(inputFileName)
		if err != nil {
			log.Fatal("Error opening input file", err)
		}
		decoder = json.NewDecoder(file)
		err = decoder.Decode(&values)
		if err != nil {
			log.Fatal("Error parsing input file:", err)
		}
	}
	conf.Training = Training
	conf.CreateTrain = CreateTrain
	if len(conf.DataFile) < 6 || len(conf.NetFile) < 6 || len(conf.Layers) < 2 || conf.DesiredError < 0.000000001 {
		log.Fatal("Config file contains error", conf)
	}
	return
}
