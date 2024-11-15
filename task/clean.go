package task

import (
	"github.com/nitezs/pcgamedb/db"

	"go.uber.org/zap"
)

func Clean(logger *zap.Logger) {
	ids, err := db.DeduplicateGames()
	if err != nil {
		logger.Error("Failed to deduplicate games", zap.Error(err))
	}
	for _, id := range ids {
		logger.Info("Deduplicated game", zap.Any("game_id", id))
	}
	idmap, err := db.CleanOrphanGamesInGameInfos()
	if err != nil {
		logger.Error("Failed to clean orphan games", zap.Error(err))
	}
	for _, id := range idmap {
		logger.Info("Cleaned orphan game in game info", zap.Any("in", id), zap.Any("removed", idmap[id]))
	}
	ids, err = db.CleanGameInfoWithEmptyGameIDs()
	if err != nil {
		logger.Error("Failed to clean game info with empty game ids", zap.Error(err))
	}
	for _, id := range ids {
		logger.Info("Cleaned game info with empty game ids", zap.Any("game_id", id))
	}
	err = db.MergeSameNameGameInfos()
	if err != nil {
		logger.Error("Failed to merge same name game infos", zap.Error(err))
	}
}
