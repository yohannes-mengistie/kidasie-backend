package liturgydoc

import (
	"os"
	"strings"
	"testing"
)

func TestDecodeManifestReadsRowsAndSkipsComments(t *testing.T) {
	source := strings.Join([]string{
		"# comment",
		"",
		strings.Join([]string{
			"apostles",
			"Anaphora-of-the-Apostles.json",
			"Anaphora of the Apostles",
			"የሐዋርያት ቅዳሴ",
			"Complete Liturgy and Anaphora of the Apostles",
			"ሙሉ ሥርዓተ ቅዳሴና የሐዋርያት ቅዳሴ",
		}, "\t"),
		// A source file whose name contains a space is why the manifest is
		// tab-separated rather than space-separated.
		strings.Join([]string{
			"st-gregory-second",
			"Qidase gorgoryos_kalie.json",
			"Anaphora of St. Gregory II",
			"የቅዱስ ጎርጎርዮስ ካልዕ ቅዳሴ",
			"Complete Liturgy and Anaphora of St. Gregory II",
			"ሙሉ ሥርዓተ ቅዳሴና የቅዱስ ጎርጎርዮስ ካልዕ ቅዳሴ",
		}, "\t"),
	}, "\n")

	entries, err := DecodeManifest(strings.NewReader(source))
	if err != nil {
		t.Fatalf("DecodeManifest returned error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}

	if entries[1].SourceFile != "Qidase gorgoryos_kalie.json" {
		t.Errorf("source file = %q", entries[1].SourceFile)
	}

	if entries[0].NameAm != "የሐዋርያት ቅዳሴ" {
		t.Errorf("name_am = %q", entries[0].NameAm)
	}
}

func TestDecodeManifestRejectsBadRows(t *testing.T) {
	cases := map[string]string{
		"short row":     "apostles\tfile.json\tName",
		"missing field": "apostles\tfile.json\tName\t\tTitle\tTitleAm",
		"duplicate slug": "apostles\ta.json\tA\tሀ\tT\tት\n" +
			"apostles\tb.json\tB\tለ\tT\tት",
		"empty": "# only a comment\n",
	}

	for label, source := range cases {
		if _, err := DecodeManifest(strings.NewReader(source)); err == nil {
			t.Errorf("DecodeManifest accepted %s", label)
		}
	}
}

// TestRepositoryManifestIsValid guards the checked-in manifest itself, since a
// bad row there breaks every conversion.
func TestRepositoryManifestIsValid(t *testing.T) {
	file, err := os.Open("../../content/updated-liturgies.tsv")
	if err != nil {
		t.Fatalf("open manifest: %v", err)
	}
	defer file.Close()

	entries, err := DecodeManifest(file)
	if err != nil {
		t.Fatalf("DecodeManifest returned error: %v", err)
	}

	if len(entries) != 14 {
		t.Errorf("manifest lists %d liturgies, want 14", len(entries))
	}
}
