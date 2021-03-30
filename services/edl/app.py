import os
import redis
from flask import Flask, make_response
app = Flask(__name__)

redis_host = os.getenv('REDIS_HOST')

@app.route('/')
def edl():
    redisClient = redis.StrictRedis(host=redis_host,
                                    port=6379,
                                    db=0)
    urls = redisClient.smembers("edl")
    res = ""
    for url in urls:
        res += str(url, 'UTF-8') + '\n'
    response = make_response(res, 200)
    response.mimetype = "text/plain"
    return response

@app.route('/health')
def health():
    return 'OK'