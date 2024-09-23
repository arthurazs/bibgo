package main

import (
	"strings"
	"testing"
)

const ACMText = `
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

@article{2,
author = {Formby, David and Walid, Anwar and Beyah, Raheem},
title = {A Case Study in Power Substation Network Dynamics},
year = {2017},
issue_date = {June 2017},
publisher = {Association for Computing Machinery},
address = {New York, NY, USA},
volume = {1},
number = {1},
url = {https://doi.org/10.1145/3084456},
doi = {10.1145/3084456},
journal = {Proc. ACM Meas. Anal. Comput. Syst.},
month = {jun},
articleno = {19},
numpages = {24},
keywords = {scada, power grid, network characterization}
}
`

func TestNextEntry(t *testing.T) {
	cases := []struct {
		bib, entry strings.Reader
	}{
		{
			*strings.NewReader(ACMText), *strings.NewReader(ACMText[:619]),
		},
	}

	for _, c := range cases {
		got, _ := NextEntry(c.bib)
		if got != c.entry {
			t.Errorf("NextEntry(%q)\nexpected %q\n     got %q", c.bib, c.entry, got)
		}
	}
}
