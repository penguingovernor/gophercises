package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/penguingovernor/excersize/quiz/internal/quiz"
)

type cliFlags struct {
	csv   string
	timer time.Duration
}

func getFlags() (*cliFlags, error) {
	var f cliFlags
	flag.StringVar(&f.csv, "csv", "", "the csv file to parse")
	flag.DurationVar(&f.timer, "timer", 30*time.Second, "how long you have to answer all the questsions")
	flag.Parse()

	if f.csv == "" {
		return nil, fmt.Errorf("missing required flag %q", "csv")
	}

	return &f, nil
}

func main() {
	flags, err := getFlags()
	if err != nil {
		log.Fatalln("could not get flags:", err)
	}

	fd, err := os.Open(flags.csv)
	if err != nil {
		log.Fatalf("could not open file %s: %v\n", flags.csv, err)
	}
	defer fd.Close()

	game, err := quiz.NewGameFromReader(fd)
	if err != nil {
		log.Fatalln("failed to create game:", err)
	}

	if err := game.PlayWithTime(os.Stdin, os.Stdout, flags.timer); err != nil {
		log.Fatalln("error encountered during gameplay:", err)
	}

}
