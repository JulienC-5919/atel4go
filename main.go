package main

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	_ "github.com/disintegration/imaging"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	println(tailleRepertoire("."))

	srcImage, _ := imaging.Open("goblin.png")

	dstImage80 := imaging.Fit(srcImage, 80, 80, imaging.NearestNeighbor)

	imaging.Save(dstImage80, "goblin_80.png")

	a, _ := Decoder([]byte{1, 2, 3, 4})
	println(reflect.DeepEqual(a, Message{Type: 1, Charge: []byte{2, 3, 4}}))

	ch := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for range 1000 {
		go func(c chan int) {
			select {
			case c <- rand.Int():
			case <-ctx.Done():
			}
		}(ch)
	}

	println(runtime.NumGoroutine())

	time.Sleep(2 * time.Second)
	println(runtime.NumGoroutine())
}

type Message struct {
	Type   byte
	Charge []byte
}

func (m Message) Equal(other Message) bool {
	if m.Type != other.Type {
		return false
	}
	if len(m.Charge) != len(other.Charge) {
		return false
	}
	for i := range m.Charge {
		if m.Charge[i] != other.Charge[i] {
			return false
		}
	}
	return true
}

func Decoder(donnees []byte) (Message, error) {
	if len(donnees) < 1 {
		return Message{}, errors.New("données trop courtes")
	}

	msg := Message{
		Type:   donnees[0],
		Charge: donnees[1:],
	}

	return msg, nil
}

func minmaxmoy(ints ...int) (int, int, int, error) {
	if len(ints) == 0 {
		return 0, 0, 0, errors.New("aucun nombre fourni")
	}

	min, max, total := ints[0], ints[0], ints[0]

	for _, i := range ints[1:] {
		if i < min {
			min = i
		} else if i > max {
			max = i
		}
		total += i
	}

	return min, max, total / len(ints), nil
}

func tailleRepertoire(dir string) (int64, error) {

	chSomme := make(chan int64)
	var somme int64
	var wg sync.WaitGroup

	wg.Add(1)
	grTailleRepertoire(dir, chSomme, &wg)

	go func() {
		wg.Wait()
		close(chSomme)
	}()

	for taille := range chSomme {
		somme += taille
	}

	return somme, nil
}

func grTailleRepertoire(dir string, c chan int64, wg *sync.WaitGroup) {
	go func(dir string, c chan int64) {
		defer wg.Done()

		files, err := os.ReadDir(dir)

		if err != nil {
			return // todo rendre l'erreur détectable
		}

		for _, file := range files {
			if file.IsDir() {
				wg.Add(1)
				grTailleRepertoire(dir+"/"+file.Name(), c, wg)
			} else if file.Type().IsRegular() {
				info, err := file.Info()
				if err != nil {
					return // todo rendre l'erreur détectable
				}

				c <- info.Size()
			}
		}
	}(dir, c)
}
