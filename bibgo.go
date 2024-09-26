package bibgo

import (
	"errors"
	"io"
	"strings"
)

var err error

func NextEntry(bib *strings.Reader) (*strings.Reader, error) {
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
	return strings.TrimSpace(elementKey.String()), nil
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
	category               string
	key                    string
	author                 []string
	abstract               string
	title                  string
	journal                string
	year                   uint16
	keywords               []string
	volume                 string
	number                 string
	pages                  string
	doi                    string
	issn                   string
	month                  string
	issue_date             string
	publisher              string
	address                string
	url                    string
	numpages               uint16
	articleno              uint16
	note                   string
	affiliations           []string
	author_keywords        []string
	correspondence_address []string
	language               string
	abbrev_source_title    string
	publication_stage      string
	source                 string
	coden                  string
	pmid                   uint32
}

func newEntry(category string, key string) Entry { // TODO should return a pointer?
	return Entry{
		category: category,
		key:      key,
	}
}

func ParseEntry(entry *strings.Reader) (Entry, error) { // TODO should return a pointer?
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
	var count uint8 = 0
	for {
		element, err = getNextElement(entry)
		if err != nil {
			return Entry{}, errors.New("getNextElement: " + err.Error())
		}
		switch element.key {
		case "author":
			println(">>", element.value)
		case "}", "":
			println("end")
		default:
			println("!!unknown", "|"+element.key+"|")
		}
		if element.key == "}" || element.key == "" {
			break
		}
		count++
		// println(count, element.key, element.value)
	}
	return parsed_entry, err
}

// TODO @arthurazs: add parseEntry, ParseFile
