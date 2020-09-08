package main

import (
	"bufio"
	"crypto/sha1"
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

func readCards(fileName string) map[string]*card {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return cardsMap
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func updatePrices(cardsMap map[string]*card) {
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
	sort.Strings(cardNames)

	q := fmt.Sprintf("%s%s", API_BASE_URL, strings.Join(cardNames, "+OR+"))

	// Checking cache
	var data []byte
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(q)))
	cacheFile := fmt.Sprintf("%s.json", hash)
	if fileExists(cacheFile) {
		log.Println("Reading from cache...")
		data, _ = ioutil.ReadFile(cacheFile)
	} else {
		log.Println("Not found in cache, requesting prices...")
		resp, err := http.Get(q)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			log.Fatal("Status Code: ", resp.StatusCode)
		}
		data, _ = ioutil.ReadAll(resp.Body)
		err = ioutil.WriteFile(cacheFile, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	var cards Cards
	json.Unmarshal(data, &cards)

	for i := 0; i < len(cards.Data); i++ {
		name := cards.Data[i].Name
		price, _ := strconv.ParseFloat(cards.Data[i].Prices.Usd, 32)
		if c, found := cardsMap[name]; found {
			c.price = price
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], " <deck-file>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	cardsMap := readCards(fileName)
	updatePrices(cardsMap)

	cardNames := make([]string, len(cardsMap))
	i := 0
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
