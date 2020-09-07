<?php
define('URL', 'https://api.scryfall.com/cards/search?order=name&q=');

/**
 * Deck.
 */
class Deck {
    private $cards = array();

    /**
     * Add card to deck.
     * @param $name card name.
     * @param $amount card amount.
     */
    public function addCard($name, $amount) {
        $n = strtolower($name);
        if (isset($cards[$n])) {
            die("Card $n already mentioned in deck. Aborting.\n");
        }
        $this->cards[$n] = new Card($n, $amount);
    }

    /**
     * Update card prices using Black Lotus Project API.
     */
    public function updatePrices() {
        $cs = '';
        foreach ($this->cards as $n => $c) {
            $cs .= urlencode($n) . '+OR+';
        }

        $url = URL . $cs;

        $json = file_get_contents($url);
        $data = json_decode($json, true);


        foreach ($data['data'] as $card) {
            $name = strtolower($card['name']);
            if (is_null($this->cards[$name])) {
                // print("...IGNORED '$name'");
                continue;
            }
            $this->cards[$name]->updatePrice($card['prices']['usd']);
        }
    }

    /**
     * @return printable format of the deck card prices.
     */
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
				$fmt = numfmt_create('en_US', NumberFormatter::CURRENCY);
        foreach ($ord as $n => $v) {
            $c = $this->cards[$n]; 

            $a = $c->getAmount();
            $p = $c->getPrice();
            $t = $a * $p;

            $sn = str_pad(ucwords($n), $npad);
            $sa = str_pad($a, $apad);
            //$sp = money_format("%#" . $ppad . ".2n", $p);
            //$st = money_format("%#" . $tpad . ".2n", $t);
						$sp = numfmt_format_currency($fmt, $p, "USD");
						$st = numfmt_format_currency($fmt, $t, "USD");

            $s .= "$sa x $sn = $sa x$sp =$st" . PHP_EOL;
            $total += $t;
            $tc += $a;
        }
        $s .= "TOTAL: $tc cards : $total" . PHP_EOL;
        return $s;
    }
}

/**
 * Card.
 */
class Card {
    private $name; 
    private $amount;
    private $price;

    /**
     * Constructor.
     * @param $name card name.
     * @param $amount number of copies of that card.
     */
    function __construct($name, $amount) {
        $this->name = $name;
        $this->amount = $amount;
    }

    /**
     * Update price variable.
     * Always keeps the lowest price.
     * @param $price price to be tested against the current price.
     */
    public function updatePrice($price) {
        if (is_null($this->price) || $this->price > $price) {
            $this->price = $price;
        }
    }

    /**
     * @return the amount.
     */
    public function getAmount() {
        return $this->amount;
    }

    /**
     * @return the price.
     */
    public function getPrice() {
        return $this->price;
    }

    /**
     * @return the name.
     */
    public function getName() {
        return $this->name;
    }
}

/**
 * Deck reader interface.
 */
interface DeckReader {
    /**
     * Read the given file and creates a Deck object.
     * @param file name of the file with the deck description.
     * @return deck object created from the given file.
     */
    public function readDeck($file);
}

/**
 * Deck reader that understands the simple txt format.
 * 
 * The format expected is:
 * <amount> <card name>
 *
 * Example:
 * 2 Mountain Goat
 * 4 Giant Growth
 * 10 Mountain
 * 8 Forest
 */
class TxtDeckReader implements DeckReader {
    /**
     * Read deck from txt file.
     * @param $file name of the file with the deck description.
     * @return deck object.
     */
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

/**
 * Deck reader that understands the Cockatrice deck format.
 * @param $file file with deck description in cockatrice format.
 * @return deck object.
 */
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
