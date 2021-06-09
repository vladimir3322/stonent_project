# Loader

The main task of this service is downloading all nft-images from IPFS and sending them to RabbitMQ.

After launching, service waits RabbitMQ connection.
Then it makes request to Stonent Eth contract to get list of supported nft contract addresses.
After that, service starts to perform the next things for each nft-contract:

1. Download all existed URI events;
2. Listen new URI events during other stuff computing.

From URI event we can retrieve IPFS-address for nft-image metadata.
This metadata includes another IPFS-address for nft-image source.
That's all we need.
If downloading is succeeded, image source will be sent to RabbitMQ in base64 format.

However, there is a problem.
There is great amount of nft-pictures and if we start to make requests to IPFS nodes, they are able to became dead.
To prevent this occasion, special buffer has been developed.
Before request to the IPFS, all request data tries to be pushed to the buffer.
The buffer has fixed size which indicates the max count of pending requests to IPFS.
To make request, there should be free place in the buffer.

When all contracts are processed, the finite message will be sent to RabbitMQ.
It indicates that no more images will be pushed.

In additional to the functionality described above, there is an HTTP-server.
It is used only by ML-service and not available outside.
It has three endpoints:

1. Downloading image source by contract address and nft-id;
2. Providing detailed information about failed images downloading (can be received from [ML-service](../ml/swagger.yml));
3. Providing some statistics stuff (can be received from [ML-service](../ml/swagger.yml)).
