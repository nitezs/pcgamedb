package db

func IsGnarlyCrawled(flag string) bool {
	return IsGameCrawled(flag, "gnarly")
}
