package main

import (
	"io"
	"strings"
)

func NextEntry(bib strings.Reader) (strings.Reader, error) {
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

// const Helper = `
// @article{1,
// author = {Ahmad, Waqar and Hasan, Osman and Tahar, Sofiene},
// title = {Formal reliability and failure analysis of ethernet based communication networks in a smart grid substation},
// year = {2020},
// issue_date = {Feb 2020},
// publisher = {Springer-Verlag},
// address = {Berlin, Heidelberg},
// volume = {32},
// number = {1},
// issn = {0934-5043},
// url = {https://doi.org/10.1007/s00165-019-00503-1},
// doi = {10.1007/s00165-019-00503-1},
// journal = {Form. Asp. Comput.},
// month = {feb},
// pages = {71–111},
// numpages = {41},
// keywords = {Theorem proving, Higher-order logic, Fault tree, Reliability block diagrams, Smart grid}
// }
// `

func main() {
	// r, _ := NextEntry(*strings.NewReader(Helper))
	// println(r.Len())
}
