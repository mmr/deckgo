package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Cards struct {
	Data []struct {
		Name   string `json:"name"`
		Prices struct {
			Usd interface{} `json:"usd"`
		} `json:"prices"`
	} `json:"data"`
}

func main() {
	response, err := http.Get("https://api.scryfall.com/cards/search?order=name&q=Palladium%20Myr+OR+Myr%20Battlesphere+OR+Myr%20Enforcer+OR+Myr%20Galvanizer+OR+Myr%20Reservoir+OR+Mirrorworks+OR+Tree%20of%20Tales+OR+Vault%20of%20Whispers+OR+Great%20Furnace+OR+Seat%20of%20the%20Synod+OR+Lodestone%20Myr+OR+Thoughtcast+OR+Darksteel%20Citadel+OR+Silver%20Myr+OR+Myr%20Turbine+OR+Banefire+OR+Voltaic%20Key+OR+Ancient%20Den")

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
		var name = cards.Data[i].Name
		var price = cards.Data[i].Prices.Usd
		fmt.Println(name, " = ", price)
	}
}

// curl "https://scryfall.com/search?as=grid&order=name&q="%"28type"%"3Acreature+type"%"3Atreefolk"%"29+color"%"3DG"

/*
    public function updatePrices() {
        $cs = '';
        foreach ($this->cards as $n => $c) {
            $cs .= urlencode($n) . '|';
        }

        $url = URL . $cs;

        $json = file_get_contents($url);
        $data = json_decode($json, true);

        foreach ($data['cards'] as $card) {
            $name = strtolower($card['name']);
            $this->cards[$name]->updatePrice($card['average']);
        }
    }

    public function __toString() {
        $ord = array();

        // calculate paddings and sort by price (low -> high)
        $tpad = 0;
        $ppad = 0;
        $apad = 0;
        $npad = 0;
        foreach ($this->cards as $n => $c) {
            $p = $c->getPrice();
            $a = $c->getAmount();
            $t = $p * $a;

            $npad = max($npad, strlen($n));
            $ppad = max($ppad, strlen((int)$p));
            $tpad = max($tpad, strlen((int)$t));
            $apad = max($apad, strlen($a));

            $ord[$n] = $t;
        }
        asort($ord);

        $s = "";
        $tc = 0;
        $total = 0;
        foreach ($ord as $n => $v) {
            $c = $this->cards[$n];

            $a = $c->getAmount();
            $p = $c->getPrice();
            $t = $a * $p;

            $sn = str_pad(ucwords($n), $npad);
            $sa = str_pad($a, $apad);
            $sp = money_format("%#" . $ppad . ".2n", $p);
            $st = money_format("%#" . $tpad . ".2n", $t);

            $s .= "$sa x $sn = $sa x$sp =$st" . PHP_EOL;
            $total += $t;
            $tc += $a;
        }
        $s .= "TOTAL: $tc cards : $total" . PHP_EOL;
        return $s;
    }
}

class Card {
    private $name;
    private $amount;
    private $price;

    function __construct($name, $amount) {
        $this->name = $name;
        $this->amount = $amount;
    }

    public function updatePrice($price) {
        if (is_null($this->price) || $this->price > $price) {
            $this->price = $price;
        }
    }

    public function getAmount() {
        return $this->amount;
    }

    public function getPrice() {
        return $this->price;
    }

    public function getName() {
        return $this->name;
    }
}

interface DeckReader {
    public function readDeck($file);
}

class TxtDeckReader implements DeckReader {
    public function readDeck($file) {
        $deck = new Deck();
        $lines = file($file);
        foreach ($lines as $line) {
            if (preg_match("/^(\d+)\s*(.*)$/", $line, $matches)) {
                $deck->addCard(trim($matches[2]), $matches[1]);
            } else {
                die("Incorret format at line: $line\n");
            }
        }
        return $deck;
    }
}

class CockatriceDeckReader implements DeckReader {
    public function readDeck($file) {
        $deck = new Deck();
        $xml = simplexml_load_file($file);
        foreach ($xml->zone->card as $card) {
            $deck->addCard($card['name'], $card['number']);
        }
        return $deck;
    }
}

if (!isset($argv[1])) {
    die("Usage: " . $argv[0] . " file\n");
}

$file = $argv[1];
if (!is_readable($file)) {
    die("File $file not found\n");
}

$ext = pathinfo($file, PATHINFO_EXTENSION);
$reader = $ext == 'cod' ? new CockatriceDeckReader() : new TxtDeckReader();
$deck = $reader->readDeck($file);
$deck->updatePrices();
print $deck;
?>
*/
