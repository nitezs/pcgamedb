package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameInfo struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	Aliases     []string             `json:"aliases" bson:"aliases"`
	Developers  []string             `json:"developers" bson:"developers"`
	Publishers  []string             `json:"publishers" bson:"publishers"`
	IGDBID      int                  `json:"igdb_id" bson:"igdb_id"`
	SteamID     int                  `json:"steam_id" bson:"steam_id"`
	GOGID       int                  `json:"-" bson:"gog_id"`
	Cover       string               `json:"cover" bson:"cover"`
	Languages   []string             `json:"languages" bson:"languages"`
	Screenshots []string             `json:"screenshots" bson:"screenshots"`
	GameIDs     []primitive.ObjectID `json:"game_ids" bson:"games"`
	Games       []*GameItem          `json:"game_downloads" bson:"-"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

type GameItem struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"speculative_name" bson:"name"`
	RawName    string             `json:"raw_name,omitempty" bson:"raw_name"`
	Download   string             `json:"download_link,omitempty" bson:"download"`
	Size       string             `json:"size,omitempty" bson:"size"`
	Url        string             `json:"url" bson:"url"`
	Password   string             `json:"password,omitempty" bson:"password"`
	Author     string             `json:"author,omitempty" bson:"author"`
	UpdateFlag string             `json:"-" bson:"update_flag,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
