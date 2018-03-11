package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	. "github.com/ematvey/gostat"
)

type Dice struct {
	HitAt int `json:hitAt`
	No    int `json:no`
}

func main() {
	//dice := make(map[int]int)
	dice := make([]Dice, 0)
	diceBody := []byte(os.Args[1])
	err := json.Unmarshal(diceBody, &dice)
	if err != nil {
		fmt.Printf("Error: %#v \n", err)
	}

	turnHits := turn(dice)
	var keys []int
	for k := range turnHits {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, i := range keys {
		if i > 0 {
			hitWord := "hit"
			if i > 1 {
				hitWord = "hits"
			}
			fmt.Printf("Probablity of rolling at least %v %v is %.1f%%  \n", i, hitWord, turnHits[i]*100)
		} else {
			fmt.Printf("Probablity of rolling %v hits is %.1f%%  \n", i, turnHits[i]*100)
		}
	}
}

func turn(dice []Dice) map[int]float64 {

	//Create Initial Data to append
	firstDie := dice[0]
	otherDie := dice[1:]
	turnHitData := getHits(firstDie)

	for _, die := range otherDie {
		//Each new iteration replace the current turn data
		newHitData := make(map[int]float64)

		//Loop through hit,prop outcome for this roll
		rollHitData := getHits(die)
		for rollHits, rollProb := range rollHitData {
			if rollHits > 0 {
				//Capture all the hit prob for this number
				newHitData[rollHits] = rollProb + newHitData[rollHits]
				//Calculate exact hit prob for multiplication
				rollProb = rollProb - rollHitData[rollHits+1]
			}
			for turnHits, turnProb := range turnHitData {
				if turnHits > 0 {
					culmulativeHits := turnHits + rollHits
					culmulativeProb := turnProb * rollProb
					newHitData[culmulativeHits] = newHitData[culmulativeHits] + culmulativeProb
				}
			}
		}
		newHitData[0] = turnHitData[0] * rollHitData[0]

		turnHitData = newHitData
	}
	return turnHitData
}

func getHits(die Dice) map[int]float64 {

	p := float64(die.HitAt) / float64(6)
	n := int64(die.No)

	hitRates := make(map[int]float64)

	calc := Binomial_CDF(p, n)
	//have to do zero differently
	hitRates[0] = calc(int64(0))

	for i := 0; i < int(die.No); i++ {
		prob := 1 - calc(int64(i))
		atLeastNumber := i + 1
		hitRates[atLeastNumber] = prob
	}

	return hitRates
}
