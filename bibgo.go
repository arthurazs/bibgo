package bibgo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var err error

func nextEntry(bib *strings.Reader) (*strings.Reader, error) {
	var buffer []byte = make([]byte, 1)
	var entry strings.Builder = strings.Builder{}
	var found bool = false
	var counter int16 = 0

	for {
		_, err = bib.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		entry.Write(buffer)

		isOpen := buffer[0] == []byte("{")[0]
		if !found {
			found = isOpen
		}

		if isOpen {
			counter++
		}
		if buffer[0] == []byte("}")[0] {
			counter--
		}

		if found && counter == 0 {
			break
		}
	}

	return strings.NewReader(entry.String()), err
}

func getCategory(entry *strings.Reader) (string, error) {
	var buffer []byte = make([]byte, 1)
	var category strings.Builder = strings.Builder{}
	var foundAt bool = false
	for {
		_, err = entry.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if !foundAt {
			foundAt = buffer[0] == []byte("@")[0]
			continue
		}
		if buffer[0] == []byte("{")[0] {
			break
		}
		category.Write(buffer)
	}
	return strings.ToLower(category.String()), nil
}

func getKey(entry *strings.Reader) (string, error) {
	var buffer []byte = make([]byte, 1)
	var key strings.Builder = strings.Builder{}
	for {
		_, err = entry.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if buffer[0] == []byte(",")[0] {
			break
		}
		key.Write(buffer)
	}
	return key.String(), nil
}

func getElementKey(entry *strings.Reader) (string, error) {
	var buffer []byte = make([]byte, 1)
	var elementKey strings.Builder = strings.Builder{}
	for {
		_, err = entry.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if buffer[0] == []byte("=")[0] {
			break
		}
		elementKey.Write(buffer)
	}
	return strings.TrimSpace(strings.TrimPrefix(elementKey.String(), ",")), nil
}

func getElementValue(entry *strings.Reader) (string, error) {
	var buffer []byte = make([]byte, 1)
	var elementValue strings.Builder = strings.Builder{}
	var started bool = false
	var counter int16 = 0
	for {
		_, err = entry.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		isOpen := buffer[0] == []byte("{")[0]
		if !started {
			if isOpen {
				started = true
				counter++
			}
			continue
		}

		if isOpen {
			counter++
		}
		if buffer[0] == []byte("}")[0] {
			counter--
		}

		if started && counter == 0 {
			break
		}
		elementValue.Write(buffer)
	}
	return strings.TrimSpace(elementValue.String()), nil
}

type bibElement struct {
	key   string
	value string
}

func getNextElement(entry *strings.Reader) (bibElement, error) { // TODO should return a pointer?
	var key string
	var value string

	key, err = getElementKey(entry)
	if err != nil {
		return bibElement{}, err
	}

	value, err = getElementValue(entry)
	if err != nil {
		return bibElement{}, err
	}

	return bibElement{key, value}, err
}

type bibEntry struct {
	category              string
	key                   string
	author                []string
	abstract              string
	title                 string
	journal               string
	year                  uint16
	keywords              []string
	volume                string
	number                string
	pages                 string
	doi                   string
	issn                  string
	month                 string
	issueDate             string
	publisher             string
	address               string
	url                   string
	numpages              uint16
	articleno             uint16
	note                  string
	affiliations          []string
	authorKeywords        []string
	correspondenceAddress []string
	language              string
	abbrevSourceTitle     string
	publicationStage      string
	source                string
	coden                 string
	pmid                  uint32
}

func newEntry(category string, key string) bibEntry { // TODO should return a pointer?
	return bibEntry{
		category: category,
		key:      key,
	}
}

func splitAndTrim(text string, separator string) []string {
	var values []string = strings.Split(text, separator)
	var index int
	var value string

	for index, value = range values {
		values[index] = strings.TrimSpace(value)
	}
	return values
}

func string2uint16(text string, name string) uint16 {
	var value uint64
	value, err = strconv.ParseUint(text, 10, 16)
	if err != nil {
		fmt.Printf("Could not parse %s: %v", name, err)
		return 0
	}
	return uint16(value)
}

func string2uint32(text string, name string) uint32 {
	var value uint64
	value, err = strconv.ParseUint(text, 10, 32)
	if err != nil {
		fmt.Printf("Could not parse %s: %v", name, err)
		return 0
	}
	return uint32(value)
}

func parseEntry(entry *strings.Reader) (bibEntry, error) { // TODO should return a pointer?
	var category, key string

	category, err = getCategory(entry)
	if err != nil {
		return bibEntry{}, errors.New("getCategory: " + err.Error())
	}

	key, err = getKey(entry)
	if err != nil {
		return bibEntry{}, errors.New("getKey: " + err.Error())
	}

	var parsed_entry bibEntry = newEntry(category, key)
	var element bibElement

Loop:
	for {
		element, err = getNextElement(entry)
		if err != nil {
			return bibEntry{}, errors.New("getNextElement: " + err.Error())
		}
		switch strings.ToLower(element.key) {
		case "author":
			parsed_entry.author = splitAndTrim(element.value, " and ")
		case "abstract":
			parsed_entry.abstract = element.value
		case "title":
			parsed_entry.title = element.value
		case "journal":
			parsed_entry.journal = element.value
		case "year":
			parsed_entry.year = string2uint16(element.value, "year")
		case "keywords":
			parsed_entry.keywords = splitAndTrim(element.value, ",")
		case "volume":
			parsed_entry.volume = element.value
		case "number":
			parsed_entry.number = element.value
		case "pages":
			parsed_entry.pages = element.value
		case "doi":
			parsed_entry.doi = element.value
		case "issn":
			parsed_entry.issn = element.value
		case "month":
			parsed_entry.month = element.value
		case "issue_date":
			parsed_entry.issueDate = element.value
		case "publisher":
			parsed_entry.publisher = element.value
		case "address":
			parsed_entry.address = element.value
		case "url":
			parsed_entry.url = element.value
		case "numpages":
			parsed_entry.numpages = string2uint16(element.value, "numpages")
		case "articleno":
			parsed_entry.articleno = string2uint16(element.value, "articleno")
		case "note":
			parsed_entry.note = element.value
		case "affiliations":
			parsed_entry.affiliations = splitAndTrim(element.value, ";")
		case "author_keywords":
			parsed_entry.authorKeywords = splitAndTrim(element.value, ";")
		case "correspondence_address":
			parsed_entry.correspondenceAddress = splitAndTrim(element.value, ";")
		case "language":
			parsed_entry.language = element.value
		case "abbrev_source_title":
			parsed_entry.abbrevSourceTitle = element.value
		case "publication_stage":
			parsed_entry.publicationStage = element.value
		case "source":
			parsed_entry.source = element.value
		case "coden":
			parsed_entry.coden = element.value
		case "pmid":
			parsed_entry.pmid = string2uint32(element.value, "pmid")
		case "}", "":
			break Loop
		case "type":
			if parsed_entry.category != strings.ToLower(element.value) {
				fmt.Printf("Entry category \"%s\" differs from element type \"%s\"\n", parsed_entry.category, element.value)
			}
		default:
			fmt.Println("Skipping unknown element: ", element.key)
		}
	}
	return parsed_entry, nil
}

func ParseFile(filePath string) uint64 {
	fmt.Printf("Parsing %s...\n", filePath)
	var raw_data []byte
	raw_data, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var raw_entry *strings.Reader
	var counter uint64 = 0
	var data *strings.Reader = strings.NewReader(string(raw_data))
	for {
		raw_entry, err = nextEntry(data)
		if err == io.EOF {
			fmt.Printf("Reached end of file, read %d entries...\n", counter)
			return counter
		}
		if err != nil {
			panic(err)
		}
		_, err = parseEntry(raw_entry)
		if err != nil {
			fmt.Printf("parseEntry: %v\n", err)
			return counter
		}
		counter++
	}
}
