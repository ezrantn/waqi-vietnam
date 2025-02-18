package main

import (
	"context"
	"time"
)

type AirQualityService interface {
	GetByCity(ctx context.Context, city string) (*AirQuality, error)
}

type WAQIService interface {
	GetByCity(ctx context.Context, city string) (*AirQuality, error)
}

type CacheService interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string)
}
