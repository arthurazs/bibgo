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

type Element struct {
	key   string
	value string
}

func getNextElement(entry *strings.Reader) (Element, error) { // TODO should return a pointer?
	var key string
	var value string

	key, err = getElementKey(entry)
	if err != nil {
		return Element{}, err
	}

	value, err = getElementValue(entry)
	if err != nil {
		return Element{}, err
	}

	return Element{key, value}, err
}

type Entry struct {
	Category              string
	Key                   string
	Author                []string
	Abstract              string
	Title                 string
	Journal               string
	Year                  uint16
	Keywords              []string
	Volume                string
	Number                string
	Pages                 string
	Doi                   string
	Issn                  string
	Month                 string
	IssueDate             string
	Publisher             string
	Address               string
	Url                   string
	Numpages              uint16
	Articleno             uint16
	Note                  string
	Affiliations          []string
	AuthorKeywords        []string
	CorrespondenceAddress []string
	Language              string
	AbbrevSourceTitle     string
	PublicationStage      string
	Source                string
	Coden                 string
	Pmid                  uint32
}

func newEntry(category string, key string) Entry { // TODO should return a pointer?
	return Entry{
		Category: category,
		Key:      key,
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

func parseEntry(entry *strings.Reader) (Entry, error) { // TODO should return a pointer?
	var category, key string

	category, err = getCategory(entry)
	if err != nil {
		return Entry{}, errors.New("getCategory: " + err.Error())
	}

	key, err = getKey(entry)
	if err != nil {
		return Entry{}, errors.New("getKey: " + err.Error())
	}

	var parsed_entry Entry = newEntry(category, key)
	var element Element

Loop:
	for {
		element, err = getNextElement(entry)
		if err != nil {
			return Entry{}, errors.New("getNextElement: " + err.Error())
		}
		switch strings.ToLower(element.key) {
		case "author":
			parsed_entry.Author = splitAndTrim(element.value, " and ")
		case "abstract":
			parsed_entry.Abstract = element.value
		case "title":
			parsed_entry.Title = element.value
		case "journal":
			parsed_entry.Journal = element.value
		case "year":
			parsed_entry.Year = string2uint16(element.value, "year")
		case "keywords":
			parsed_entry.Keywords = splitAndTrim(element.value, ",")
		case "volume":
			parsed_entry.Volume = element.value
		case "number":
			parsed_entry.Number = element.value
		case "pages":
			parsed_entry.Pages = element.value
		case "doi":
			parsed_entry.Doi = element.value
		case "issn":
			parsed_entry.Issn = element.value
		case "month":
			parsed_entry.Month = element.value
		case "issue_date":
			parsed_entry.IssueDate = element.value
		case "publisher":
			parsed_entry.Publisher = element.value
		case "address":
			parsed_entry.Address = element.value
		case "url":
			parsed_entry.Url = element.value
		case "numpages":
			parsed_entry.Numpages = string2uint16(element.value, "numpages")
		case "articleno":
			parsed_entry.Articleno = string2uint16(element.value, "articleno")
		case "note":
			parsed_entry.Note = element.value
		case "affiliations":
			parsed_entry.Affiliations = splitAndTrim(element.value, ";")
		case "author_keywords":
			parsed_entry.AuthorKeywords = splitAndTrim(element.value, ";")
		case "correspondence_address":
			parsed_entry.CorrespondenceAddress = splitAndTrim(element.value, ";")
		case "language":
			parsed_entry.Language = element.value
		case "abbrev_source_title":
			parsed_entry.AbbrevSourceTitle = element.value
		case "publication_stage":
			parsed_entry.PublicationStage = element.value
		case "source":
			parsed_entry.Source = element.value
		case "coden":
			parsed_entry.Coden = element.value
		case "pmid":
			parsed_entry.Pmid = string2uint32(element.value, "pmid")
		case "}", "":
			break Loop
		case "type":
			if parsed_entry.Category != strings.ToLower(element.value) {
				fmt.Printf("Entry category \"%s\" differs from element type \"%s\"\n", parsed_entry.Category, element.value)
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
