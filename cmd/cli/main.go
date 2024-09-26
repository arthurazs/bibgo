package main

import (
	"bibgo"
	"fmt"
	"strings"
)

const debugText string = `
@article{1,
author = {Ahmad, Waqar and Hasan, Osman and Tahar, Sofiene},
title = {Formal reliability and failure analysis of ethernet based communication networks in a smart grid substation},
year = {2020},
issue_date = {Feb 2020},
publisher = {Springer-Verlag},
address = {Berlin, Heidelberg},
volume = {32},
number = {1},
issn = {0934-5043},
url = {https://doi.org/10.1007/s00165-019-00503-1},
doi = {10.1007/s00165-019-00503-1},
journal = {Form. Asp. Comput.},
month = {feb},
pages = {71–111},
numpages = {41},
keywords = {Theorem proving, Higher-order logic, Fault tree, Reliability block diagrams, Smart grid}
}
`

func main() {
	parsed_entry, _ := bibgo.ParseEntry(strings.NewReader(debugText))
	fmt.Println(parsed_entry)
	fmt.Printf("\n%+v\n", parsed_entry)
}
