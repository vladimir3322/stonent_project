import pika
import config
import json
import time
import base64
import globals


def connect_with_retry():
    try:
        credentials = pika.PlainCredentials(config.rabbit_login, config.rabbit_password)
        parameters = pika.ConnectionParameters(config.rabbit_host, config.rabbit_port, '/', credentials)
        connection = pika.BlockingConnection(parameters=parameters)

        if connection.is_open:
            return connection
        else:
            raise Exception('Connection closed')
    except Exception as e:
        print(e)
        print('Wait until Rabbit becomes alive', flush=True)
        time.sleep(1)
        return connect_with_retry()


def get_safe_consumer(connection, queue):
    try:
        channel = connection.channel()

        return channel, channel.consume(queue)
    except Exception as e:
        print(e)
        print('Wait until Rabbit becomes alive', flush=True)
        time.sleep(1)
        return get_safe_consumer(connection, queue)


def safe_consume(queue):
    connection = connect_with_retry()
    channel, consumer = get_safe_consumer(connection, queue)

    while True:
        try:
            for method_frame, properties, body in consumer:
                channel.basic_ack(method_frame.delivery_tag)

                yield method_frame, properties, body
        except Exception as e:
            print(e)
            print('Wait until Rabbit becomes alive', flush=True)
            connection = connect_with_retry()
            channel, consumer = get_safe_consumer(connection, queue)

        time.sleep(1)


def consume_events():
    print('Start Rabbit', flush=True)

    for method_frame, properties, body in safe_consume(config.rabbit_queue):
        data = json.loads(body)

        contract_address = data['contractAddress']
        nft_id = data['nftID']
        image_bytes_source = base64.b64decode(data['data'])
        is_finite = data['isFinite']

        if is_finite:
            globals.all_images_has_been_downloaded = True
            print('Actual images have been downloaded', flush=True)
            continue

        print('Received from Rabbit:', nft_id, contract_address, flush=True)

        yield contract_address, nft_id, image_bytes_source
