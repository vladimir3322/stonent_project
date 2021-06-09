package config

const ServerPort = 8080
const ProviderUrl = "wss://mainnet.infura.io/ws/v3/844de29fabee4fcebf315309262d0836"

var IpfsLink = []string{"https://ipfs.daonomic.com", "https://ipfs.io"}

const RabbitLogin = "guest"
const RabbitPass = "guest"
const RabbitHost = "rabbitmq"
const RabbitPort = "5672"
const QueueIndexing = "indexing"

const MlUrl = "http://ml:9090"

const DownloadImageBufferSize = 4
const DownloadImageMaxCount = -1 // -1 for ignoring
