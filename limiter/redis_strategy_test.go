package limiter

import (
	"context"
	"testing"
	"time"
)

func TestRedisLimiter_AllowRequest(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    5,
		MaxRequestsToken: 50,
		BlockTime:        1 * time.Minute,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)
	ctx := context.Background()

	// Limpar o Redis para garantir um estado inicial
	rdb.Del(ctx, "limiter:test_user")

	// Testar que as primeiras 5 requisições são permitidas
	for i := 0; i < 5; i++ {
		if !limiter.AllowRequest("test_user", false) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// A 6ª requisição deve ser bloqueada
	if limiter.AllowRequest("test_user", false) {
		t.Error("expected 6th request to be blocked")
	}
}

func TestRedisLimiter_Token(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    10,
		MaxRequestsToken: 50,
		BlockTime:        1 * time.Minute,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)
	ctx := context.Background()

	// Limpar o Redis para garantir um estado inicial
	rdb.Del(ctx, "limiter:token_user")

	// Testar que as primeiras 50 requisições são permitidas
	for i := 0; i < 50; i++ {
		if !limiter.AllowRequest("token_user", true) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// A 51ª requisição deve ser bloqueada
	if limiter.AllowRequest("token_user", true) {
		t.Error("expected 51st request to be blocked")
	}
}

func TestRedisLimiter_BlockTimeLeft(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    5,
		MaxRequestsToken: 50,
		BlockTime:        1 * time.Minute,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)
	ctx := context.Background()

	// Limpar o Redis
	rdb.Del(ctx, "limiter:test_user_block")

	// Testar que a requisição é permitida e verifica o tempo de bloqueio
	limiter.AllowRequest("test_user_block", false)
	limiter.AllowRequest("test_user_block", false)

	ttl := limiter.BlockTimeLeft("test_user_block")
	if ttl <= 0 {
		t.Error("expected BlockTimeLeft to be greater than 0")
	}
}
