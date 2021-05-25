import redis
import config

redisClient = redis.from_url(config.redis_url)


def consume_events():
    while redisClient.llen(config.redis_job_queue) != 0:
        yield redisClient.lpop(config.redis_job_queue)

    return None
