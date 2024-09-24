package model

type SteamAppDetail struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string        `json:"type"`
		Name                string        `json:"name"`
		SteamAppid          int           `json:"steam_appid"`
		RequiredAge         any           `json:"required_age"`
		IsFree              bool          `json:"is_free"`
		ControllerSupport   string        `json:"controller_support"`
		DetailedDescription string        `json:"detailed_description"`
		AboutTheGame        string        `json:"about_the_game"`
		ShortDescription    string        `json:"short_description"`
		SupportedLanguages  string        `json:"supported_languages"`
		HeaderImage         string        `json:"header_image"`
		CapsuleImage        string        `json:"capsule_image"`
		CapsuleImagev5      string        `json:"capsule_imagev5"`
		Website             string        `json:"website"`
		PcRequirements      any           `json:"pc_requirements"`
		MacRequirements     any           `json:"mac_requirements"`
		LinuxRequirements   any           `json:"linux_requirements"`
		LegalNotice         string        `json:"legal_notice"`
		Developers          []string      `json:"developers"`
		Publishers          []string      `json:"publishers"`
		PackageGroups       []interface{} `json:"package_groups"`
		Platforms           struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int    `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Mp4 struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"mp4"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		BackgroundRaw      string `json:"background_raw"`
		ContentDescriptors struct {
			Ids   []interface{} `json:"ids"`
			Notes interface{}   `json:"notes"`
		} `json:"content_descriptors"`
		Ratings struct {
			Esrb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
				UseAgeGate  string `json:"use_age_gate"`
				RequiredAge string `json:"required_age"`
			} `json:"esrb"`
			Pegi struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"pegi"`
			Oflc struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"oflc"`
		} `json:"ratings"`
	} `json:"data"`
}

type SteamPackageDetail struct {
	Success bool `json:"success"`
	Data    struct {
		Name        string `json:"name"`
		PageContent string `json:"page_content"`
		PageImage   string `json:"page_image"`
		HeaderImage string `json:"header_image"`
		SmallLogo   string `json:"small_logo"`
		Apps        []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"apps"`
		Price struct {
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
			Individual      int    `json:"individual"`
		} `json:"price"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Controller struct {
			FullGamepad bool `json:"full_gamepad"`
		} `json:"controller"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
	} `json:"data"`
}
