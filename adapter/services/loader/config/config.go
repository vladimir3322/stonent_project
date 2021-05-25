package config

const ServerPort = 5000
const ProviderUrl = "wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836"

var IpfsLink = []string{"https://ipfs.daonomic.com", "https://ipfs.io"}

const RabbitMQUrl = "amqp://rabbitmq:rabbitmq@rabbit1:5672/"
const RabbitMQQueueName = "imageSources"

const RedisUrl = "redis:6379"
const RedisJobQueue = "imageSources"

const DownloadImageBufferSize = 2
const DownloadImageMaxCount = -1 // -1 for ignoring

const RabbitLogin = "guest"
const RabbitPass = "guest"
const RabbitHost = "rabbitmq"
const RabbitPort = "5672"
const QueueIndexing = "indexing"
