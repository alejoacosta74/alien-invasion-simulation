package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/alejoacosta74/allien_invasion/app"
	"github.com/alejoacosta74/allien_invasion/log"
	"github.com/alejoacosta74/allien_invasion/world"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	aliens_num     = kingpin.Flag("aliens", "Number of aliens. Defaults to system's number of CPUs.").Default(strconv.Itoa(runtime.NumCPU())).Short('a').Int()
	iterations_num = kingpin.Flag("iterations", "Number of iterations. Defaults to 10.000.").Default(strconv.Itoa(10000)).Short('i').Int()
	debug          = kingpin.Flag("debug", "debug mode").Short('d').Default("false").Bool()
	sourceFile     = kingpin.Flag("file", "input source file").Required().Short('f').String()
)
var logger *logrus.Logger

func init() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	mainLogger, err := log.GetLogger(
		log.WithDebugLevel(*debug),
		log.WithWriter(os.Stdout),
	)
	if err != nil {
		logrus.Panic(err)
	}
	logger = mainLogger
}

func main() {
	if sourceFile == nil {
		logger.Fatal("No input file provided")
	}
	//read world from file
	worldMap, err := world.NewWorldFromFile(*sourceFile)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("World loaded from file")
	worldMap.PrintWorld()

	//channel to handle Ctrl+C signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//create context to cancel goroutines (aliens)
	ctx, cancelFunc := context.WithCancel(context.Background())

	// starts aliens and coordinates their movement until the aliens are destroyed or the iterations are completed
	done := app.StartInvation(ctx, cancelFunc, *aliens_num, *iterations_num, sigs, worldMap)

	select {
	case <-sigs:
		// handle Ctrl+C signal from user input
		logger.Warn("Received ^C ... exiting")
		cancelFunc()
		os.Exit(0)
	case <-done:
		logger.Info("Invasion terminated ... ")
		worldMap.PrintWorld()
	}
	logger.Print("Program finished")
}
