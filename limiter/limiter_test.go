package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// Função auxiliar para inicializar o Redis
func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// Teste para limitação por IP
func TestRateLimiter_IP(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    5,
		MaxRequestsToken: 50,
		BlockTime:        1 * time.Minute,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)

	identifier := "192.168.1.1"

	// Realiza 5 requisições, todas devem ser permitidas
	for i := 0; i < 5; i++ {
		if !limiter.AllowRequest(identifier, false) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// A 6ª requisição deve ser bloqueada
	if limiter.AllowRequest(identifier, false) {
		t.Error("expected 6th request to be blocked")
	}
}

// Teste para limitação por Token
func TestRateLimiter_Token(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    10,
		MaxRequestsToken: 50,
		BlockTime:        1 * time.Minute,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)

	token := "abc123"

	// Realiza 50 requisições, todas devem ser permitidas
	for i := 0; i < 50; i++ {
		if !limiter.AllowRequest(token, true) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// A 51ª requisição deve ser bloqueada
	if limiter.AllowRequest(token, true) {
		t.Error("expected 51st request to be blocked")
	}
}

// Teste para expiração do tempo de bloqueio
func TestRateLimiter_BlockTimeExpiration(t *testing.T) {
	rdb := setupRedis()
	defer rdb.FlushAll(context.Background()) // Limpar Redis após o teste

	config := RateLimiterConfig{
		MaxRequestsIP:    3,
		MaxRequestsToken: 50,
		BlockTime:        2 * time.Second,
		RedisAddress:     "localhost:6379",
	}
	limiter := NewRedisLimiter(config)

	identifier := "192.168.1.2"

	// Realiza 3 requisições, todas devem ser permitidas
	for i := 0; i < 3; i++ {
		if !limiter.AllowRequest(identifier, false) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	// A 4ª requisição deve ser bloqueada
	if limiter.AllowRequest(identifier, false) {
		t.Error("expected 4th request to be blocked")
	}

	// Espera o tempo de bloqueio expirar
	time.Sleep(3 * time.Second)

	// Agora a requisição deve ser permitida novamente
	if !limiter.AllowRequest(identifier, false) {
		t.Error("expected request to be allowed after block time expiration")
	}
}
