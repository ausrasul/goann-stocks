package main

import (
	"fmt"
	"github.com/white-pony/go-fann"
	"log"
)

func main() {
	conf, values := getConf("")
	log.Print(conf)
	if conf.Training {
		numLayers := uint(len(conf.Layers))
		layers := conf.Layers
		desiredError := conf.DesiredError
		maxEpochs := conf.MaxEpochs
		epochsBetweenReports := conf.EpochsBetweenReports


		ann := fann.CreateStandard(numLayers, layers)
        ann.SetTrainingAlgorithm(fann.TRAIN_RPROP)

		ann.SetActivationFunctionHidden(fann.GAUSSIAN)
		ann.SetActivationFunctionOutput(fann.SIN_SYMMETRIC)
        //ann.SetLearningMomentum(0.5)
        //ann.SetLearningRate(0.2)
		ann.TrainOnFile(conf.DataFile, maxEpochs, epochsBetweenReports, desiredError)
		ann.Save(conf.NetFile)
		ann.Destroy()
	} else if conf.CreateTrain {
		createStockData(conf)
	} else {
		var cash float32 = 20000
		boughtDays := 0
		positions := 0
		var price float32 = 0
		for _, row := range values.Values{
			if positions > 0 && boughtDays > 0 {
				// count down to sell.
				boughtDays--
			}
			fmt.Println("Date: ", row.Date, " Close price: ", row.Close, " High: ", row.High)
			fmt.Println("Cash: ", cash, " Positions: " , positions, " Days left: ", boughtDays)
			inputs := row.Inputs
			
			ann := fann.CreateFromFile(conf.NetFile)
			output := ann.Run(inputs)
			for col, o := range output {
				//a := denormalize(conf.MinPrice, conf.MaxPrice, float32(o))
				fmt.Println(" Result: ", o, " Actual: ", row.Outputs[col])
				if col == 1 {
					if positions <= 0 {
						if o >= 0.7{
							// buy
							boughtDays = 5
							//positions = 1
							for cash >= row.Close {
								positions++
								cash = cash - row.Close
							}
							price = row.Close
							fmt.Println("Bought ", positions, " stocks for ", row.Close, " each, looking for: ", price * 1.02, " ========================== ")
						}
					}
				}
			}
			goalPrice := price * 1.02
			if positions > 0 && row.High >= goalPrice && boughtDays < 5 {
				// Sell :)
				fmt.Print("Selling ", positions, " stocks for profit ", goalPrice)
				for positions > 0 {
					cash = cash + goalPrice
					positions--
				}
				boughtDays = 0
				fmt.Println(" current cash: ", cash)
			}
			if positions > 0 && boughtDays <= 0 {
				//sell :(
				fmt.Print("Forced to sell ", positions, " stocks for ", row.Close)
				for positions > 0 {
					cash = cash + row.Close
					positions--
				}
                fmt.Println(" current cash: ", cash)
			}
		}
		/*input := values.Values//[]fann.FannType{values.Values}
		fmt.Println("Input: ", values)
		for _, j := range input{
			b := denormalize(conf.MinPrice, conf.MaxPrice, float32(j))
			fmt.Println("input:", b)
		}
		
		ann := fann.CreateFromFile(conf.NetFile)
		for round := 0; round < rounds; round++{
			output := ann.Run(input)
			fmt.Println("input:", input)
			for _, o := range output {
				a := denormalize(conf.MinPrice, conf.MaxPrice, float32(o))
				fmt.Println("Result: ", a)
			}
			input = shiftInputs(input, output)
		}
		ann.Destroy()*/
	}
}

func shiftInputs(input []fann.FannType, output []fann.FannType)([]fann.FannType){
	for i,_ := range input{
		if i < len(input) -1{
			input[i] = input[i+1]
		} else {
			input[i] = output[0]
		}
	}
	return input
}
