# Stonent

### Fast start

Just run using
```
docker compose up
```

To check process status use:
```
GET http://localhost:9090/statistics
```

To get image score use:
```
GET http://localshost:9090/check?contract_address=CONTRACT_ADDRESS&nft_id=NFT_ID
```

### What about stuffing?

There are three services which interact with each other:
1. [Loader](./services/loader/README.md) - downloading nft-images from IPFS;
2. [ML](./services/ml/README.md) - neural network checking plagiarism;
3. RabbitMQ - database-mediator between Loader and ML for comfortable nft-images transfer.

The whole system has such scheme:

![Scheme](./media/scheme.png)
