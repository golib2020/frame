package queue

const LuaScriptPop = `
local job = redis.call('lpop', KEYS[1])
local reserved = false
if(job) then
 	reserved = job
    redis.call('zadd', KEYS[2], ARGV[1], reserved)
end
return reserved
`

const LuaScriptDelete = `
local reserved = ''
reserved = cjson.decode(ARGV[1])
reserved = cjson.encode(reserved)
return redis.call('zrem', KEYS[1], reserved)
`

const LuaScriptRelease = `
reserved = cjson.decode(ARGV[1])
reserved = cjson.encode(reserved)
redis.call('zrem', KEYS[2], reserved)
redis.call('zadd', KEYS[1], ARGV[2], reserved)
return true
`

const LuaScriptMigrateExpiredJobs = `
local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])
if(next(val) ~= nil) then
    redis.call('zremrangebyrank', KEYS[1], 0, #val - 1)
    for i = 1, #val, 100 do
        redis.call('rpush', KEYS[2], unpack(val, i, math.min(i+99, #val)))
    end
end
return true
`

const LuaScriptSize = `
return redis.call('llen', KEYS[1]) + redis.call('zcard', KEYS[2]) + redis.call('zcard', KEYS[3])
`
