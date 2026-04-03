package redis

import goredis "github.com/redis/go-redis/v9"

var allowScript = goredis.NewScript(`
local key = KEYS[1]
local now_ms = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local refill_per_ms = tonumber(ARGV[3])
local ttl_seconds = tonumber(ARGV[4])

local data = redis.call('HMGET', key, 'tokens', 'last_refill_ms')
local tokens = tonumber(data[1])
local last_refill_ms = tonumber(data[2])

if tokens == nil or last_refill_ms == nil then
  tokens = capacity
  last_refill_ms = now_ms
end

local elapsed = now_ms - last_refill_ms
if elapsed > 0 then
  tokens = math.min(capacity, tokens + (elapsed * refill_per_ms))
  last_refill_ms = now_ms
end

local allowed = 0
local remaining = math.floor(math.max(tokens - 1, 0))
local reset_seconds = math.max(0, math.ceil((capacity - tokens) / (refill_per_ms * 1000)))
local retry_after_seconds = 0

if tokens >= 1 then
  tokens = tokens - 1
  allowed = 1
  remaining = math.floor(tokens)
else
  retry_after_seconds = math.max(1, math.ceil((1 - tokens) / (refill_per_ms * 1000)))
  if reset_seconds < retry_after_seconds then
    reset_seconds = retry_after_seconds
  end
end

redis.call('HMSET', key, 'tokens', tokens, 'last_refill_ms', last_refill_ms)
redis.call('EXPIRE', key, ttl_seconds)

return {allowed, remaining, reset_seconds, retry_after_seconds}
`)
