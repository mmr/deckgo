package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const API_BASE_URL = "https://api.scryfall.com/cards/search?q="

type Cards struct {
	Data []struct {
		Name   string `json:"name"`
		Prices struct {
			Usd string `json:"usd"`
		} `json:"prices"`
	} `json:"data"`
}

type card struct {
	name   string
	amount int
	price  float64
}

/*
func cleanName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
*/

func readCards(fileName string) (map[string]*card, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cardsMap := make(map[string]*card)
	re := regexp.MustCompile(`^(\d+)\s*(.*)$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		vals := re.FindStringSubmatch(line)
		amount, _ := strconv.Atoi(vals[1])
		name := vals[2]
		cardsMap[name] = &card{name: name, amount: amount}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cardsMap, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], " <deck-file>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	cardsMap, err := readCards(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	i := 0
	cardNames := make([]string, len(cardsMap))
	for k := range cardsMap {
		if cardsMap[k].amount > 4 {
			// Probably land. Ignore
			continue
		}
		cardNames[i] = url.QueryEscape(k)
		i++
	}

	q := fmt.Sprintf("%s%s", API_BASE_URL, strings.Join(cardNames, "+OR+"))
	response, err := http.Get(q)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cards Cards
	json.Unmarshal(responseData, &cards)

	for i := 0; i < len(cards.Data); i++ {
		name := cards.Data[i].Name
		price, _ := strconv.ParseFloat(cards.Data[i].Prices.Usd, 32)
		if c, found := cardsMap[name]; found {
			c.price = price
		}
	}

	i = 0
	for k := range cardsMap {
		cardNames[i] = k
		i++
	}
	sort.Strings(cardNames)

	cardsTotal := 0
	priceTotal := 0.0
	for _, name := range cardNames {
		c := cardsMap[name]
		amount := c.amount
		price := c.price
		total := float64(amount) * price
		fmt.Printf("%2d x %-20v = $%.2f\n", amount, name, total)
		cardsTotal += amount
		priceTotal += total
	}
	fmt.Printf("%2d %-22v = $%.2f\n", cardsTotal, "cards", priceTotal)
}
