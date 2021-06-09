server_port = 9090

redis_url = 'redis://redis:6379'
redis_job_queue = 'imageSources'

loader_url = 'http://loader:8080'

rabbit_login = "guest"
rabbit_password = "guest"
rabbit_host = "rabbitmq"
rabbit_port = 5672
rabbit_queue = "indexing"

registered_images_file = './registered_images.txt'
rejected_images_by_IPFS_file = './rejected_images_by_IPFS.txt'
rejected_images_by_NN_file = './rejected_images_by_NN.txt'
