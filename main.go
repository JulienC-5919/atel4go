package main

import (
	"errors"
	"os"
	"sync"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	println(tailleRepertoire("."))
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
