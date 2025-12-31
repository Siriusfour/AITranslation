-- scripts/seckill.lua
local stock_key = KEYS[1]
local user_key  = KEYS[2]
local user_id   = ARGV[1]

-- 1. 校验用户是否已抢购
if redis.call('sismember', user_key, user_id) == 1 then
    return -1
end

-- 2. 校验库存
local stock = tonumber(redis.call('get', stock_key))
if stock == nil then
    return -2
end
if stock <= 0 then
    return 0
end

-- 3. 扣库存 + 记录用户
redis.call('decr', stock_key)
redis.call('sadd', user_key, user_id)

return 1
