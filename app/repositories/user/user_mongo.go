package user

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepo struct {
	mongoDB *mongo.Client
}

func NewUserMongoRepo(mongoDB *mongo.Client) *UserMongoRepo {
	return &UserMongoRepo{
		mongoDB: mongoDB,
	}
}

func (r *UserMongoRepo) List(ctx context.Context) {

}
