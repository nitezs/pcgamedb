package db

func IsChovkaCrawled(flag string) bool {
	return IsGameCrawled(flag, "chovka")
}
