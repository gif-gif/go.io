package goredisc

import (
	"context"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gogf/gf/container/garray"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoRedisC struct {
	Redis  *redis.ClusterClient
	Config Config
	Ctx    context.Context
}

func New(conf Config) (cli *GoRedisC, err error) {
	if conf.Type == "" {
		conf.Type = "node"
	}

	cli = &GoRedisC{
		Ctx:    context.Background(),
		Config: conf,
	}

	cli.Redis = redis.NewClusterClient(&redis.ClusterOptions{
		// 指定集群中的节点地址，至少需要一个有效节点
		Addrs: conf.Addrs,

		// 连接池配置
		PoolSize:     10, // 连接池大小
		MinIdleConns: 5,  // 最小空闲连接数

		// 认证信息（如果需要）
		Password: conf.Password, // 如果有密码，请设置

		// 读写超时设置
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,

		// 连接超时
		DialTimeout: 5 * time.Second,

		// 最大重试次数
		MaxRetries: 3,

		// 集群刷新间隔
		ClusterSlots: func(ctx context.Context) ([]redis.ClusterSlot, error) {
			// 自定义集群槽位获取逻辑，通常使用默认值即可
			return nil, nil
		},
	})

	ctx := context.Background()
	if err = cli.Redis.Ping(ctx).Err(); err != nil {
		golog.WithTag("goredis").Error(err)
		return
	}

	if conf.AutoPing {
		gj, _ := gojob.New()
		gj.Start()
		gj.SecondX(nil, 5, func() {
			if err := cli.Redis.Ping(ctx).Err(); err != nil {
				golog.WithTag("goredis").Fatal("redis ping error:", err)
			}
		})
	}

	return cli, nil
}

func (s *GoRedisC) WrapKey(key string) string {
	if s.Config.Prefix == "" {
		return key
	}
	return s.Config.Prefix + ":" + key
}

func (s *GoRedisC) WrapKeys(keys ...string) []string {
	arr := garray.NewStrArrayFromCopy(keys)
	return arr.Walk(func(val string) string { return s.WrapKey(val) }).Slice()
}

// BitCount is redis bitcount command implementation.
func (s *GoRedisC) BitCount(key string, start, end int64) *redis.IntCmd {
	return s.Redis.BitCount(s.Ctx, s.WrapKey(key), &redis.BitCount{Start: start, End: end})
}

// BitOpAnd is redis bit operation (and) command implementation.
func (s *GoRedisC) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	return s.Redis.BitOpAnd(s.Ctx, s.WrapKey(destKey), s.WrapKeys(keys...)...)
}

// BitOpNot is redis bit operation (not) command implementation.
func (s *GoRedisC) BitOpNot(destKey, key string) *redis.IntCmd {
	return s.Redis.BitOpNot(s.Ctx, s.WrapKey(destKey), s.WrapKey(key))
}

// BitOpOr is redis bit operation (or) command implementation.
func (s *GoRedisC) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	return s.Redis.BitOpOr(s.Ctx, s.WrapKey(destKey), s.WrapKeys(keys...)...)
}

// BitOpXor is redis bit operation (xor) command implementation.
func (s *GoRedisC) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	return s.Redis.BitOpXor(s.Ctx, s.WrapKey(destKey), s.WrapKeys(keys...)...)
}

// BitPos is redis bitpos command implementation.
func (s *GoRedisC) BitPos(key string, bit, start, end int64) *redis.IntCmd {
	return s.Redis.BitPos(s.Ctx, s.WrapKey(key), bit, start, end)
}

// Blpop uses passed in redis connection to execute blocking queries.
// Doesn't benefit from pooling redis connections of blocking queries
func (s *GoRedisC) BLPop(timeout time.Duration, key string) *redis.StringSliceCmd {
	return s.Redis.BLPop(s.Ctx, timeout, s.WrapKey(key))
}

// Del deletes keys.
func (s *GoRedisC) Del(keys ...string) *redis.IntCmd {
	return s.Redis.Del(s.Ctx, s.WrapKeys(keys...)...)
}

// Eval is the implementation of redis eval command.
func (s *GoRedisC) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return s.Redis.Eval(s.Ctx, script, s.WrapKeys(keys...), args...)
}

// EvalSha is the implementation of redis evalsha command.
func (s *GoRedisC) EvalSha(sha string, keys []string, args ...interface{}) *redis.Cmd {
	return s.Redis.EvalSha(s.Ctx, sha, s.WrapKeys(keys...), args...)
}

// Exists is the implementation of redis exists command.
func (s *GoRedisC) Exists(key string) *redis.IntCmd {
	return s.Redis.Exists(s.Ctx, s.WrapKey(key))
}

// Expire is the implementation of redis expire command.
func (s *GoRedisC) Expire(key string, seconds int) *redis.BoolCmd {
	return s.Redis.Expire(s.Ctx, s.WrapKey(key), time.Duration(seconds)*time.Second)
}

// Expireat is the implementation of redis expireat command.
func (s *GoRedisC) ExpireAt(key string, time time.Time) *redis.BoolCmd {
	return s.Redis.ExpireAt(s.Ctx, s.WrapKey(key), time)
}

// GeoAdd is the implementation of redis geoadd command.
func (s *GoRedisC) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return s.Redis.GeoAdd(s.Ctx, s.WrapKey(key), geoLocation...)
}

// GeoDist is the implementation of redis geodist command.
func (s *GoRedisC) GeoDist(key, member1, member2, unit string) *redis.FloatCmd {
	return s.Redis.GeoDist(s.Ctx, s.WrapKey(key), member1, member2, unit)
}

// GeoHash is the implementation of redis geohash command.
func (s *GoRedisC) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	return s.Redis.GeoHash(s.Ctx, s.WrapKey(key), members...)
}

// GeoRadius is the implementation of redis georadius command.
func (s *GoRedisC) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return s.Redis.GeoRadius(s.Ctx, s.WrapKey(key), longitude, latitude, query)
}

// GeoRadiusByMember is the implementation of redis georadiusbymember command.
func (s *GoRedisC) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return s.Redis.GeoRadiusByMember(s.Ctx, s.WrapKey(key), member, query)
}

// GeoPos is the implementation of redis geopos command.
func (s *GoRedisC) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	return s.Redis.GeoPos(s.Ctx, s.WrapKey(key), members...)
}

// Get is the implementation of redis get command.
func (s *GoRedisC) Get(key string) *redis.StringCmd {
	return s.Redis.Get(s.Ctx, s.WrapKey(key))
}

// GetBit is the implementation of redis getbit command.
func (s *GoRedisC) GetBit(key string, offset int64) *redis.IntCmd {
	return s.Redis.GetBit(s.Ctx, s.WrapKey(key), offset)
}

// Hdel is the implementation of redis hdel command.
func (s *GoRedisC) HDel(key string, fields ...string) *redis.IntCmd {
	return s.Redis.HDel(s.Ctx, s.WrapKey(key), fields...)
}

// Hexists is the implementation of redis hexists command.
func (s *GoRedisC) HExists(key, field string) *redis.BoolCmd {
	return s.Redis.HExists(s.Ctx, s.WrapKey(key), field)
}

// Hget is the implementation of redis hget command.
func (s *GoRedisC) HGet(key, field string) *redis.StringCmd {
	return s.Redis.HGet(s.Ctx, s.WrapKey(key), field)
}

// Hgetall is the implementation of redis hgetall command.
func (s *GoRedisC) HGetAll(key string) *redis.MapStringStringCmd {
	return s.Redis.HGetAll(s.Ctx, s.WrapKey(key))
}

// Hincrby is the implementation of redis hincrby command.
func (s *GoRedisC) HIncrBy(key, field string, increment int64) *redis.IntCmd {
	return s.Redis.HIncrBy(s.Ctx, s.WrapKey(key), field, increment)
}

// Hkeys is the implementation of redis hkeys command.
func (s *GoRedisC) HKeys(key string) *redis.StringSliceCmd {
	return s.Redis.HKeys(s.Ctx, s.WrapKey(key))
}

// Hlen is the implementation of redis hlen command.
func (s *GoRedisC) HLen(key string) *redis.IntCmd {
	return s.Redis.HLen(s.Ctx, s.WrapKey(key))
}

// Hmget is the implementation of redis hmget command.
func (s *GoRedisC) HMGet(key string, fields ...string) *redis.SliceCmd {
	return s.Redis.HMGet(s.Ctx, s.WrapKey(key), fields...)
}

// Hset is the implementation of redis hset command.
func (s *GoRedisC) HSet(key, field, value string) *redis.IntCmd {
	return s.Redis.HSet(s.Ctx, s.WrapKey(key), field, value)
}

// Hsetnx is the implementation of redis hsetnx command.
func (s *GoRedisC) HSetNX(key, field, value string) *redis.BoolCmd {
	return s.Redis.HSetNX(s.Ctx, s.WrapKey(key), field, value)
}

// Hmset is the implementation of redis hmset command.
func (s *GoRedisC) HMSet(key string, fieldsAndValues map[string]string) *redis.BoolCmd {
	return s.Redis.HMSet(s.Ctx, s.WrapKey(key), fieldsAndValues)
}

// Hscan is the implementation of redis hscan command.
func (s *GoRedisC) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return s.Redis.HScan(s.Ctx, s.WrapKey(key), cursor, match, count)
}

// Hvals is the implementation of redis hvals command.
func (s *GoRedisC) HVals(key string) *redis.StringSliceCmd {
	return s.Redis.HVals(s.Ctx, s.WrapKey(key))
}

// Incr is the implementation of redis incr command.
func (s *GoRedisC) Incr(key string) *redis.IntCmd {
	return s.Redis.Incr(s.Ctx, s.WrapKey(key))
}

// Incrby is the implementation of redis incrby command.
func (s *GoRedisC) IncrBy(key string, increment int64) *redis.IntCmd {
	return s.Redis.IncrBy(s.Ctx, s.WrapKey(key), increment)
}

// Keys is the implementation of redis keys command.
func (s *GoRedisC) Keys(pattern string) *redis.StringSliceCmd {
	return s.Redis.Keys(s.Ctx, pattern)
}

// Llen is the implementation of redis llen command.
func (s *GoRedisC) LLen(key string) *redis.IntCmd {
	return s.Redis.LLen(s.Ctx, s.WrapKey(key))
}

// Lpop is the implementation of redis lpop command.
func (s *GoRedisC) LPop(key string) *redis.StringCmd {
	return s.Redis.LPop(s.Ctx, s.WrapKey(key))
}

// Lpush is the implementation of redis lpush command.
func (s *GoRedisC) LPush(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.LPush(s.Ctx, s.WrapKey(key), values...)
}

// Lrange is the implementation of redis lrange command.
func (s *GoRedisC) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return s.Redis.LRange(s.Ctx, s.WrapKey(key), start, stop)
}

// Lrem is the implementation of redis lrem command.
func (s *GoRedisC) LRem(key string, count int64, value string) *redis.IntCmd {
	return s.Redis.LRem(s.Ctx, s.WrapKey(key), count, value)
}

// Mget is the implementation of redis mget command.
func (s *GoRedisC) MGet(keys ...string) *redis.SliceCmd {
	return s.Redis.MGet(s.Ctx, s.WrapKeys(keys...)...)
}

// Persist is the implementation of redis persist command.
func (s *GoRedisC) Persist(key string) *redis.BoolCmd {
	return s.Redis.Persist(s.Ctx, s.WrapKey(key))
}

// Pfadd is the implementation of redis pfadd command.
func (s *GoRedisC) PFAdd(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.PFAdd(s.Ctx, s.WrapKey(key), values...)
}

// Pfcount is the implementation of redis pfcount command.
func (s *GoRedisC) PFCount(key string) *redis.IntCmd {
	return s.Redis.PFCount(s.Ctx, s.WrapKey(key))
}

// Pfmerge is the implementation of redis pfmerge command.
func (s *GoRedisC) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return s.Redis.PFMerge(s.Ctx, dest, s.WrapKeys(keys...)...)
}

// Ping is the implementation of redis ping command.
func (s *GoRedisC) Ping() *redis.StatusCmd {
	return s.Redis.Ping(s.Ctx)
}

// Pipelined lets fn to execute pipelined commands.
// fn key must call GetKey or GetKeys to add prefix.
func (s *GoRedisC) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return s.Redis.Pipelined(s.Ctx, fn)
}

// Rpop is the implementation of redis rpop command.
func (s *GoRedisC) RPop(key string) *redis.StringCmd {
	return s.Redis.RPop(s.Ctx, s.WrapKey(key))
}

// Rpush is the implementation of redis rpush command.
func (s *GoRedisC) RPush(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.RPush(s.Ctx, s.WrapKey(key), values...)
}

// Sadd is the implementation of redis sadd command.
func (s *GoRedisC) SAdd(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.SAdd(s.Ctx, s.WrapKey(key), values...)
}

// Scan is the implementation of redis scan command.
func (s *GoRedisC) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return s.Redis.Scan(s.Ctx, cursor, match, count)
}

// SetBit is the implementation of redis setbit command.
func (s *GoRedisC) SetBit(key string, offset int64, value int) *redis.IntCmd {
	return s.Redis.SetBit(s.Ctx, s.WrapKey(key), offset, value)
}

// Sscan is the implementation of redis sscan command.
func (s *GoRedisC) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return s.Redis.SScan(s.Ctx, s.WrapKey(key), cursor, match, count)
}

// Scard is the implementation of redis scard command.
func (s *GoRedisC) SCard(key string) *redis.IntCmd {
	return s.Redis.SCard(s.Ctx, s.WrapKey(key))
}

// ScriptLoad is the implementation of redis script load command.
func (s *GoRedisC) ScriptLoad(script string) *redis.StringCmd {
	return s.Redis.ScriptLoad(s.Ctx, script)
}

// Set is the implementation of redis set command.
func (s *GoRedisC) Set(key, value string) *redis.StatusCmd {
	return s.Redis.Set(s.Ctx, s.WrapKey(key), value, 0)
}

func (s *GoRedisC) Set1(key, value string, expiration time.Duration) *redis.StatusCmd {
	return s.Redis.Set(s.Ctx, s.WrapKey(key), value, expiration)
}

// Setex is the implementation of redis setex command.
func (s *GoRedisC) SetEx(key, value string, expiration time.Duration) *redis.StatusCmd {
	return s.Redis.SetEx(s.Ctx, s.WrapKey(key), value, expiration)
}

func (s *GoRedisC) SetNX(key, value string, expiration time.Duration) *redis.BoolCmd {
	return s.Redis.SetNX(s.Ctx, s.WrapKey(key), value, expiration)
}

// Sismember is the implementation of redis sismember command.
func (s *GoRedisC) SIsMember(key string, value interface{}) *redis.BoolCmd {
	return s.Redis.SIsMember(s.Ctx, s.WrapKey(key), value)
}

// Smembers is the implementation of redis smembers command.
func (s *GoRedisC) SMembers(key string) *redis.StringSliceCmd {
	return s.Redis.SMembers(s.Ctx, s.WrapKey(key))
}

// Spop is the implementation of redis spop command.
func (s *GoRedisC) SPop(key string) *redis.StringCmd {
	return s.Redis.SPop(s.Ctx, s.WrapKey(key))
}

// Srandmember is the implementation of redis srandmember command.
func (s *GoRedisC) SRandMember(key string) *redis.StringCmd {
	return s.Redis.SRandMember(s.Ctx, s.WrapKey(key))
}

// Srem is the implementation of redis srem command.
func (s *GoRedisC) SRem(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.SRem(s.Ctx, s.WrapKey(key), values...)
}

// Sunion is the implementation of redis sunion command.
func (s *GoRedisC) SUnion(keys ...string) *redis.StringSliceCmd {
	return s.Redis.SUnion(s.Ctx, s.WrapKeys(keys...)...)
}

// Sunionstore is the implementation of redis sunionstore command.
func (s *GoRedisC) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return s.Redis.SUnionStore(s.Ctx, destination, s.WrapKeys(keys...)...)
}

// Sdiff is the implementation of redis sdiff command.
func (s *GoRedisC) SDiff(keys ...string) *redis.StringSliceCmd {
	return s.Redis.SDiff(s.Ctx, s.WrapKeys(keys...)...)
}

// Sdiffstore is the implementation of redis sdiffstore command.
func (s *GoRedisC) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return s.Redis.SDiffStore(s.Ctx, destination, s.WrapKeys(keys...)...)
}

// Sinter is the implementation of redis sinter command.
func (s *GoRedisC) SInter(keys ...string) *redis.StringSliceCmd {
	return s.Redis.SInter(s.Ctx, s.WrapKeys(keys...)...)
}

// Sinterstore is the implementation of redis sinterstore command.
func (s *GoRedisC) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return s.Redis.SInterStore(s.Ctx, destination, s.WrapKeys(keys...)...)
}

// Ttl is the implementation of redis ttl command.
func (s *GoRedisC) TTL(key string) *redis.DurationCmd {
	return s.Redis.TTL(s.Ctx, s.WrapKey(key))
}

// Zadd is the implementation of redis zadd command.
func (s *GoRedisC) ZAdd(key string, value ...redis.Z) *redis.IntCmd {
	return s.Redis.ZAdd(s.Ctx, s.WrapKey(key), value...)
}

// Zcard is the implementation of redis zcard command.
func (s *GoRedisC) ZCard(key string) *redis.IntCmd {
	return s.Redis.ZCard(s.Ctx, s.WrapKey(key))
}

// Zcount is the implementation of redis zcount command.
func (s *GoRedisC) ZCount(key string, max, min string) *redis.IntCmd {
	return s.Redis.ZCount(s.Ctx, s.WrapKey(key), min, max)
}

// Zincrby is the implementation of redis zincrby command.
func (s *GoRedisC) ZIncrBy(key string, increment float64, field string) *redis.FloatCmd {
	return s.Redis.ZIncrBy(s.Ctx, s.WrapKey(key), increment, field)
}

// Zscore is the implementation of redis zscore command.
func (s *GoRedisC) ZScore(key, value string) *redis.FloatCmd {
	return s.Redis.ZScore(s.Ctx, s.WrapKey(key), value)
}

// Zrank is the implementation of redis zrank command.
func (s *GoRedisC) ZRank(key, field string) *redis.IntCmd {
	return s.Redis.ZRank(s.Ctx, s.WrapKey(key), field)
}

// Zrem is the implementation of redis zrem command.
func (s *GoRedisC) Zrem(key string, values ...interface{}) *redis.IntCmd {
	return s.Redis.ZRem(s.Ctx, s.WrapKey(key), values...)
}

// Zremrangebyscore is the implementation of redis zremrangebyscore command.
func (s *GoRedisC) ZRemRangeByScore(key string, max, min string) *redis.IntCmd {
	return s.Redis.ZRemRangeByScore(s.Ctx, s.WrapKey(key), min, max)
}

// Zremrangebyrank is the implementation of redis zremrangebyrank command.
func (s *GoRedisC) Zremrangebyrank(key string, start, stop int64) *redis.IntCmd {
	return s.Redis.ZRemRangeByRank(s.Ctx, s.WrapKey(key), start, stop)
}

// Zrange is the implementation of redis zrange command.
func (s *GoRedisC) Zrange(key string, start, stop int64) *redis.StringSliceCmd {
	return s.Redis.ZRange(s.Ctx, s.WrapKey(key), start, stop)
}

// ZrangeWithScores is the implementation of redis zrange command with scores.
func (s *GoRedisC) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return s.Redis.ZRangeWithScores(s.Ctx, s.WrapKey(key), start, stop)
}

// ZRevRangeWithScores is the implementation of redis zrevrange command with scores.
func (s *GoRedisC) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return s.Redis.ZRevRangeWithScores(s.Ctx, s.WrapKey(key), start, stop)
}

// ZrangebyscoreWithScores is the implementation of redis zrangebyscore command with scores.
func (s *GoRedisC) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return s.Redis.ZRangeByScoreWithScores(s.Ctx, s.WrapKey(key), opt)
}

// Zrevrange is the implementation of redis zrevrange command.
func (s *GoRedisC) Zrevrange(key string, start, stop int64) *redis.StringSliceCmd {
	return s.Redis.ZRevRange(s.Ctx, s.WrapKey(key), start, stop)
}

// ZrevrangebyscoreWithScores is the implementation of redis zrevrangebyscore command with scores.
func (s *GoRedisC) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return s.Redis.ZRevRangeByScoreWithScores(s.Ctx, s.WrapKey(key), opt)
}

// Zrevrank is the implementation of redis zrevrank command.
func (s *GoRedisC) ZRevRank(key, field string) *redis.IntCmd {
	return s.Redis.ZRevRank(s.Ctx, s.WrapKey(key), field)
}

// Zunionstore is the implementation of redis zunionstore command.
func (s *GoRedisC) ZUnionStore(dest string, store *redis.ZStore) *redis.IntCmd {
	return s.Redis.ZUnionStore(s.Ctx, dest, store)
}
