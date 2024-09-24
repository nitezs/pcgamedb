package constant

type language struct {
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

var IGDBLanguages map[int]language = map[int]language{
	1: {
		Name:       "Arabic",
		NativeName: "العربية",
	},
	2: {
		Name:       "Chinese (Simplified)",
		NativeName: "简体中文",
	},
	3: {
		Name:       "Chinese (Traditional)",
		NativeName: "繁體中文",
	},
	4: {
		Name:       "Czech",
		NativeName: "čeština",
	},
	5: {
		Name:       "Danish",
		NativeName: "Dansk",
	},
	6: {
		Name:       "Dutch",
		NativeName: "Nederlands",
	},
	7: {
		Name:       "English",
		NativeName: "English (US)",
	},
	8: {
		Name:       "English (UK)",
		NativeName: "English (UK)",
	},
	9: {
		Name:       "Spanish (Spain)",
		NativeName: "Español (España)",
	},
	10: {
		Name:       "Spanish (Mexico)",
		NativeName: "Español (Mexico)",
	},
	12: {
		Name:       "French",
		NativeName: "Français",
	},
	14: {
		Name:       "Hungarian",
		NativeName: "Magyar",
	},
	11: {
		Name:       "Finnish",
		NativeName: "Suomi",
	},
	15: {
		Name:       "Italian",
		NativeName: "Italiano",
	},
	13: {
		Name:       "Hebrew",
		NativeName: "עברית",
	},
	16: {
		Name:       "Japanese",
		NativeName: "日本語",
	},
	17: {
		Name:       "Korean",
		NativeName: "한국어",
	},
	18: {
		Name:       "Norwegian",
		NativeName: "Norsk",
	},
	20: {
		Name:       "Portuguese (Portugal)",
		NativeName: "Português (Portugal)",
	},
	21: {
		Name:       "Portuguese (Brazil)",
		NativeName: "Português (Brasil)",
	},
	19: {
		Name:       "Polish",
		NativeName: "Polski",
	},
	22: {
		Name:       "Russian",
		NativeName: "Русский",
	},
	24: {
		Name:       "Turkish",
		NativeName: "Türkçe",
	},
	25: {
		Name:       "Thai",
		NativeName: "ไทย",
	},
	26: {
		Name:       "Vietnamese",
		NativeName: "Tiếng Việt",
	},
	23: {
		Name:       "Swedish",
		NativeName: "Svenska",
	},
	27: {
		Name:       "German",
		NativeName: "Deutsch",
	},
	28: {
		Name:       "Ukrainian",
		NativeName: "українська",
	},
}
