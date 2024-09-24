package model

type GOGAppDetail struct {
	ID                         int    `json:"id"`
	Title                      string `json:"title"`
	PurchaseLink               string `json:"purchase_link"`
	Slug                       string `json:"slug"`
	ContentSystemCompatibility struct {
		Windows bool `json:"windows"`
		Osx     bool `json:"osx"`
		Linux   bool `json:"linux"`
	} `json:"content_system_compatibility"`
	Languages map[string]string `json:"languages"`
	Links     struct {
		PurchaseLink string `json:"purchase_link"`
		ProductCard  string `json:"product_card"`
		Support      string `json:"support"`
		Forum        string `json:"forum"`
	} `json:"links"`
	InDevelopment struct {
		Active bool        `json:"active"`
		Until  interface{} `json:"until"`
	} `json:"in_development"`
	IsSecret      bool   `json:"is_secret"`
	IsInstallable bool   `json:"is_installable"`
	GameType      string `json:"game_type"`
	IsPreOrder    bool   `json:"is_pre_order"`
	ReleaseDate   string `json:"release_date"`
	Images        struct {
		Background          string `json:"background"`
		Logo                string `json:"logo"`
		Logo2X              string `json:"logo2x"`
		Icon                string `json:"icon"`
		SidebarIcon         string `json:"sidebarIcon"`
		SidebarIcon2X       string `json:"sidebarIcon2x"`
		MenuNotificationAv  string `json:"menuNotificationAv"`
		MenuNotificationAv2 string `json:"menuNotificationAv2"`
	} `json:"images"`
	Dlcs      any `json:"dlcs"`
	Downloads struct {
		Installers []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Os           string `json:"os"`
			Language     string `json:"language"`
			LanguageFull string `json:"language_full"`
			Version      string `json:"version"`
			TotalSize    int    `json:"total_size"`
			Files        []struct {
				ID       string `json:"id"`
				Size     int    `json:"size"`
				Downlink string `json:"downlink"`
			} `json:"files"`
		} `json:"installers"`
		Patches       []interface{} `json:"patches"`
		LanguagePacks []interface{} `json:"language_packs"`
		BonusContent  []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Count     int    `json:"count"`
			TotalSize int    `json:"total_size"`
			Files     []struct {
				ID       int    `json:"id"`
				Size     int    `json:"size"`
				Downlink string `json:"downlink"`
			} `json:"files"`
		} `json:"bonus_content"`
	} `json:"downloads"`
	ExpandedDlcs []interface{} `json:"expanded_dlcs"`
	Description  struct {
		Lead             string `json:"lead"`
		Full             string `json:"full"`
		WhatsCoolAboutIt string `json:"whats_cool_about_it"`
	} `json:"description"`
	Screenshots []struct {
		ImageID              string `json:"image_id"`
		FormatterTemplateURL string `json:"formatter_template_url"`
		FormattedImages      []struct {
			FormatterName string `json:"formatter_name"`
			ImageURL      string `json:"image_url"`
		} `json:"formatted_images"`
	} `json:"screenshots"`
	Videos          []interface{} `json:"videos"`
	RelatedProducts []interface{} `json:"related_products"`
	Changelog       string        `json:"changelog"`
}

type GOGSearch struct {
	Products []struct {
		CustomAttributes []interface{} `json:"customAttributes"`
		Developer        string        `json:"developer"`
		Publisher        string        `json:"publisher"`
		Gallery          []string      `json:"gallery"`
		Video            struct {
			ID       string `json:"id"`
			Provider string `json:"provider"`
		} `json:"video"`
		SupportedOperatingSystems []string    `json:"supportedOperatingSystems"`
		Genres                    []string    `json:"genres"`
		GlobalReleaseDate         interface{} `json:"globalReleaseDate"`
		IsTBA                     bool        `json:"isTBA"`
		Price                     struct {
			Currency                   string      `json:"currency"`
			Amount                     string      `json:"amount"`
			BaseAmount                 string      `json:"baseAmount"`
			FinalAmount                string      `json:"finalAmount"`
			IsDiscounted               bool        `json:"isDiscounted"`
			DiscountPercentage         int         `json:"discountPercentage"`
			DiscountDifference         string      `json:"discountDifference"`
			Symbol                     string      `json:"symbol"`
			IsFree                     bool        `json:"isFree"`
			Discount                   int         `json:"discount"`
			IsBonusStoreCreditIncluded bool        `json:"isBonusStoreCreditIncluded"`
			BonusStoreCreditAmount     string      `json:"bonusStoreCreditAmount"`
			PromoID                    interface{} `json:"promoId"`
		} `json:"price"`
		IsDiscounted    bool        `json:"isDiscounted"`
		IsInDevelopment bool        `json:"isInDevelopment"`
		ID              int         `json:"id"`
		ReleaseDate     interface{} `json:"releaseDate"`
		Availability    struct {
			IsAvailable          bool `json:"isAvailable"`
			IsAvailableInAccount bool `json:"isAvailableInAccount"`
		} `json:"availability"`
		SalesVisibility struct {
			IsActive   bool `json:"isActive"`
			FromObject struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"fromObject"`
			From     int `json:"from"`
			ToObject struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"toObject"`
			To int `json:"to"`
		} `json:"salesVisibility"`
		Buyable    bool   `json:"buyable"`
		Title      string `json:"title"`
		Image      string `json:"image"`
		URL        string `json:"url"`
		SupportURL string `json:"supportUrl"`
		ForumURL   string `json:"forumUrl"`
		WorksOn    struct {
			Windows bool `json:"Windows"`
			Mac     bool `json:"Mac"`
			Linux   bool `json:"Linux"`
		} `json:"worksOn"`
		Category         string        `json:"category"`
		OriginalCategory string        `json:"originalCategory"`
		Rating           int           `json:"rating"`
		Type             int           `json:"type"`
		IsComingSoon     bool          `json:"isComingSoon"`
		IsPriceVisible   bool          `json:"isPriceVisible"`
		IsMovie          bool          `json:"isMovie"`
		IsGame           bool          `json:"isGame"`
		Slug             string        `json:"slug"`
		IsWishlistable   bool          `json:"isWishlistable"`
		ExtraInfo        []interface{} `json:"extraInfo"`
		AgeLimit         int           `json:"ageLimit"`
	} `json:"products"`
	Ts               interface{} `json:"ts"`
	Page             int         `json:"page"`
	TotalPages       int         `json:"totalPages"`
	TotalResults     string      `json:"totalResults"`
	TotalGamesFound  int         `json:"totalGamesFound"`
	TotalMoviesFound int         `json:"totalMoviesFound"`
}
