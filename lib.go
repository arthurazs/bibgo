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
