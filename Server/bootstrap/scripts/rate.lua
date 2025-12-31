local key = KEYS[1]
local limit = tonumber(ARGV[1])
local expire_time = tonumber(ARGV[2])

local current = tonumber(redis.call('get', key) or "0")

if current + 1 > limit then
    return 0 -- 超限
end

-- 计数 +1
redis.call("incr", key)

-- 如果是第一次（刚才读出来是0），则设置过期时间
if current == 0 then
    redis.call("expire", key, expire_time)
end

return 1 -- 通过