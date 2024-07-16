package main

import (
	"fmt"
	"sync"
	"time"
)

type Stage struct {
	name string
	workerCount int
	process func(chan int, chan int, *sync.WaitGroup)

	inChan chan int
	outChan chan int
	wg sync.WaitGroup
}

type Pipeline struct {
	stages []Stage
}

func (p *Pipeline) GetInChan() chan int {
	return p.stages[0].inChan
}

func (p *Pipeline) GetOutChan() chan int {
	return p.stages[len(p.stages) - 1].outChan
}

func main() {
	var pipelineWaitGroup sync.WaitGroup
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	pipeline := ConstructPipelineFromStages([]Stage{
		Stage{
			name: "Map",
			workerCount: 5,
			process: processMap,
		},
		Stage{
			name: "Filter",
			workerCount: 3,
			process: processFilter,
		},
		Stage{
			name: "Reduce",
			workerCount: 1,
			process: processReduce,
		},
	})

	for _, stage := range pipeline.stages {
		time.Sleep(1 * time.Second)
		fmt.Printf("Starting process for stage %s\n", stage.name)
		if stage.name == "Map" {
			fmt.Println(stage.inChan == pipeline.GetInChan())
		}
		pipelineWaitGroup.Add(1)
		go func(stage Stage) {
			stage.Process(&pipelineWaitGroup)
		}(stage)
	}

	time.Sleep(1 * time.Second)

	go func() {
		for val := range pipeline.GetOutChan() {
			fmt.Println(val)
		}
	}()

	go producer(pipeline.GetInChan(), input)

	pipelineWaitGroup.Wait()
	
	fmt.Println("Done!")
}

func ConstructPipelineFromStages(stages []Stage) Pipeline {
	inChan := make(chan int)

	for i := range stages {
		stages[i].inChan = inChan
		stages[i].outChan = make(chan int)

		inChan = stages[i].outChan
	}

	return Pipeline{stages: stages}
}

func processMap(inChan chan int, outChan chan int, wg *sync.WaitGroup) {
	fmt.Println("processMap call")
	defer wg.Done()
	for input := range inChan {
		fmt.Printf("Processing %v; map\n", input)
		if input % 2 == 0 {
			outChan <- input / 2
		}
		outChan <- 3 * input + 1
	}
}

func processFilter(inChan chan int, outChan chan int, wg *sync.WaitGroup) {
	fmt.Println("processFilter call")
	defer wg.Done()
	for input := range inChan {
		fmt.Printf("Processing %v; filter\n", input)
		if input % 2 != 0 {
			outChan <- -1
		} else {
			outChan <- input
		}
	}
}

func processReduce(inChan chan int, outChan chan int, wg *sync.WaitGroup) {
	fmt.Println("processReduce call")
	defer wg.Done()
	sum := 0
	count := 0
	for input := range inChan {
		fmt.Printf("Processing %v; reduce\n", input)
		if input == -1 {
			continue
		}
		sum += input
		count += 1
	}
	fmt.Printf("processReduce results: %v, %v; output: %v\n", sum, count, sum / count)
	outChan <- sum / count
}

func (s *Stage) Process(pipelineWaitGroup *sync.WaitGroup) {
	defer pipelineWaitGroup.Done()
	defer close(s.outChan)
	for i := 0; i < s.workerCount; i++ {
		time.Sleep(500 * time.Millisecond)
		s.wg.Add(1)
		go s.process(s.inChan, s.outChan, &s.wg)
	}

	fmt.Println("Finished spawning workers for stage " + s.name)
	s.wg.Wait()
	fmt.Println("Finished processing for stage " + s.name)
}

func producer(inChan chan int, input []int) {
	defer close(inChan)

	for _, num := range input {
		fmt.Printf("Producing %v\n", num)
		inChan <-num
	}
}
