# ML

The core of this service is neural network which make indexing nft-pictures and is able to check plagiarism.
This service is two-threading program completing the next tasks:

1. NN indexing

    New images are received from RabbitMQ and send to the NN.

2. Simple HTTP-server

    Provides interface to interact with NN and receive statistics data, see [Swagger](./swagger.yml) for more details.

After launch, service wait RabbitMQ connection.
Then service starts to consume images, indexes NN and processes HTTP-server.
