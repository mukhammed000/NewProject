package handler

import (
	pb "api/genproto/auth"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	Auth pb.AuthServiceClient
	Rdb  *redis.Client
}

func NewHandler(auth pb.AuthServiceClient, rdb redis.Client) *Handler {
	return &Handler{
		Auth: auth,
		Rdb:  &rdb,
	}
}
