package bibgo

import (
	"io"
	"strings"
)

func nextEntry(bib *strings.Reader) (strings.Reader, error) {
	var buffer []byte = make([]byte, 1)
	var entry strings.Builder = strings.Builder{}
	var found bool = false
	var counter int = 0
	var err error

	for {
		_, err = bib.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return strings.Reader{}, err
		}
		entry.Write(buffer)

		is_open := buffer[0] == []byte("{")[0]
		if !found {
			found = is_open
		}

		if is_open {
			counter++
		}
		if buffer[0] == []byte("}")[0] {
			counter--
		}

		if found && counter == 0 {
			break
		}
	}

	return *strings.NewReader(entry.String()), err
}

func getCategory(entry *strings.Reader) (string, error) {
	var buffer []byte = make([]byte, 1)
	var category strings.Builder = strings.Builder{}
	var found_at bool = false
	var err error
	for {
		_, err = entry.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if !found_at {
			found_at = buffer[0] == []byte("@")[0]
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
	var err error
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
	var err error
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
    var counter int = 0
	var err error
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
