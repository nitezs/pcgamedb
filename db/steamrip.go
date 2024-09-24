package db

func IsSteamRIPCrawled(flag string) bool {
	return IsGameCrawled(flag, "SteamRIP")
}
