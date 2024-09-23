package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisLimiter define a estrutura para trabalhar com Redis
type RedisLimiter struct {
	client           *redis.Client
	maxRequestsIP    int
	maxRequestsToken int
	blockTime        time.Duration
}

// RateLimiterConfig define a estrutura de configuração do limiter
type RateLimiterConfig struct {
	MaxRequestsIP    int
	MaxRequestsToken int
	BlockTime        time.Duration
	RedisAddress     string
}

// NewRedisLimiter cria uma nova instância do RedisLimiter
func NewRedisLimiter(config RateLimiterConfig) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
	})
	return &RedisLimiter{
		client:           rdb,
		maxRequestsIP:    config.MaxRequestsIP,
		maxRequestsToken: config.MaxRequestsToken,
		blockTime:        config.BlockTime,
	}
}

// AllowRequest verifica se uma requisição é permitida com base no IP ou Token
func (rl *RedisLimiter) AllowRequest(identifier string, isToken bool) bool {
	ctx := context.Background()
	key := "limiter:" + identifier

	var maxRequests int
	if isToken {
		maxRequests = rl.maxRequestsToken
	} else {
		maxRequests = rl.maxRequestsIP
	}

	// Verifica o número de requisições no último segundo
	reqCount, err := rl.client.Get(ctx, key).Int()
	if err == redis.Nil {
		// Não existe contador para este IP/Token, então começamos do zero
		reqCount = 0
	} else if err != nil {
		// Se ocorrer qualquer outro erro ao buscar a chave, permitimos a requisição
		return true
	}

	if reqCount >= maxRequests {
		// Se o número de requisições já alcançou o limite, bloqueamos
		return false
	}

	// Incrementa a contagem de requisições
	_, err = rl.client.Incr(ctx, key).Result()
	if err != nil {
		// Em caso de erro ao incrementar, permitimos a requisição
		return true
	}

	// Define uma expiração de 1 segundo para garantir que o limite é por segundo
	rl.client.Expire(ctx, key, time.Second)
	return true
}

// BlockTimeLeft retorna o tempo de bloqueio restante para o identificador
func (rl *RedisLimiter) BlockTimeLeft(identifier string) time.Duration {
	ctx := context.Background()
	key := "limiter:" + identifier
	ttl, _ := rl.client.TTL(ctx, key).Result()
	return ttl
}
