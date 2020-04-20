package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

type question struct {
	Problem string
	Answer  string
}

type Game struct {
	tQuestions []question
	score      int
}

func NewGameFromReader(r io.Reader) (*Game, error) {
	csvRD := csv.NewReader(r)
	records, err := csvRD.ReadAll()
	if err != nil {
		return nil, err
	}
	var g Game
	for _, row := range records {
		g.tQuestions = append(g.tQuestions, question{
			Problem: row[0],
			Answer:  row[1],
		})
	}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(g.tQuestions), func(i, j int) {
		g.tQuestions[i], g.tQuestions[j] = g.tQuestions[j], g.tQuestions[i]
	})
	return &g, nil
}

func (g *Game) Play(r io.Reader, w io.Writer) error {

	bufRD := bufio.NewReader(r)

	for i := 0; i < len(g.tQuestions); i++ {

		fmt.Fprintf(w, "Question #%d: %s?\n", i+1, g.tQuestions[i].Problem)
		fmt.Fprintf(w, "Your answer: ")

		line, err := bufRD.ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Println()

		if strings.ToLower(strings.TrimSpace(line)) == strings.ToLower(strings.TrimSpace(g.tQuestions[i].Answer)) {
			g.score++
		}

	}

	fmt.Fprintln(w, "All done!")
	fmt.Fprintf(w, "Final Score: %02d / %02d\n", g.score, len(g.tQuestions))
	return nil
}

func (g *Game) PlayWithTime(r io.Reader, w io.Writer, timeout time.Duration) error {
	gameChan := make(chan struct{}, 1)
	timerChan := make(chan struct{}, 2)
	bufRD := bufio.NewReader(r)

	fmt.Fprintf(w, "Press enter to start the %s timer...", timeout)
	if _, err := bufRD.ReadString('\n'); err != nil {
		return err
	}

	fmt.Fprintln(w, "Time started goodluck!")

	go func(c chan<- struct{}) {
		<-time.After(timeout)
		timerChan <- struct{}{}
	}(timerChan)

	go func(r io.Reader, w io.Writer, c chan<- struct{}) error {
		if err := g.Play(r, w); err != nil {
			return err
		}
		c <- struct{}{}
		return nil
	}(r, w, gameChan)

	select {
	case <-gameChan:
		fmt.Fprintln(w, "You beat the timer!")
	case <-timerChan:
		fmt.Println()
		fmt.Fprintln(w, "Sorry, time ran out!")
		fmt.Fprintf(w, "Final score: %d / %d\n", g.score, len(g.tQuestions))
	}

	return nil
}
