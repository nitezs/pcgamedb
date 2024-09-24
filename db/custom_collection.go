package db

import (
	"context"
	"pcgamedb/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomCollection struct {
	collName string
	coll     *mongo.Collection
}

func (c *CustomCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.Find(ctx, filter, opts...)
}

func (c *CustomCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.FindOne(ctx, filter, opts...)
}

func (c *CustomCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.UpdateOne(ctx, filter, update, opts...)
}

func (c *CustomCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.Aggregate(ctx, pipeline, opts...)
}

func (c *CustomCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.DeleteOne(ctx, filter, opts...)
}

func (c *CustomCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.DeleteMany(ctx, filter, opts...)
}

func (c *CustomCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	CheckConnect()
	if c.coll == nil {
		c.coll = mongoDB.Database(config.Config.Database.Database).Collection(c.collName)
	}
	return c.coll.CountDocuments(ctx, filter, opts...)
}
