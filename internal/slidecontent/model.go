package slidecontent

type Document struct {
	SchemaVersion int    `json:"schema_version"`
	SourcePDF     string `json:"source_pdf"`
	Pages         []Page `json:"pages"`
}

type Page struct {
	Number int    `json:"number"`
	Kind   string `json:"kind"`
	Role   string `json:"role,omitempty"`

	TextGeez     string `json:"text_geez,omitempty"`
	TextAmharic  string `json:"text_amharic,omitempty"`
	TextEnglish  string `json:"text_english,omitempty"`
	Title        string `json:"title,omitempty"`
	TitleGeez    string `json:"title_geez,omitempty"`
	TitleAmharic string `json:"title_amharic,omitempty"`
	TitleEnglish string `json:"title_english,omitempty"`

	Instruction        string `json:"instruction,omitempty"`
	InstructionAmharic string `json:"instruction_amharic,omitempty"`
	Reference          string `json:"reference,omitempty"`
	Note               string `json:"note,omitempty"`
	NeedsReview        *bool  `json:"needs_review,omitempty"`

	DeaconInstructionGeez    string `json:"deacon_instruction_geez,omitempty"`
	DeaconInstructionAmharic string `json:"deacon_instruction_amharic,omitempty"`
	DeaconInstructionEnglish string `json:"deacon_instruction_english,omitempty"`

	ResponsePeople        string `json:"response_people,omitempty"`
	ResponsePeopleGeez    string `json:"response_people_geez,omitempty"`
	ResponsePeopleAmharic string `json:"response_people_amharic,omitempty"`
	ResponsePeopleEnglish string `json:"response_people_english,omitempty"`
	ResponseAmharic       string `json:"response_amharic,omitempty"`
	ResponseEnglish       string `json:"response_english,omitempty"`
	TextGeezPeople        string `json:"text_geez_people,omitempty"`
	TextAmharicPeople     string `json:"text_amharic_people,omitempty"`
	TextEnglishPeople     string `json:"text_english_people,omitempty"`

	Parts         []Part `json:"parts,omitempty"`
	ResponseMixed []Part `json:"response_mixed,omitempty"`
}

type Part struct {
	Role        string `json:"role,omitempty"`
	TextGeez    string `json:"text_geez,omitempty"`
	TextAmharic string `json:"text_amharic,omitempty"`
	TextEnglish string `json:"text_english,omitempty"`
}
