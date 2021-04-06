from web3 import Web3
import json
import asyncio
import os

from . import config
from . import download_images


async def handle_event(contract, event):
    print(f'New event: {event["args"]["_id"]}')
    await download_images.download_image_data(contract, event['args']['_id'])


async def iterate_events(contract, event_filter, poll_interval):
    while True:
        for event in event_filter.get_new_entries():
            await handle_event(contract, event)
        await asyncio.sleep(poll_interval)


def listen_images():
    if not os.path.exists(config.config.DATA_PATH):
        os.mkdir(config.config.DATA_PATH)
    if not os.path.exists(config.config.SOURCE_PATH):
        os.mkdir(config.config.SOURCE_PATH)

    w3 = Web3(Web3.HTTPProvider(config.config.get_web3_provider()))
    with open(config.config.CONTRACT_ABI_FILENAME) as json_file:
        abi = json.load(json_file)
    contract = w3.eth.contract(address=config.config.CONTRACT_ADDRESS, abi=abi)

    contract_filter = contract.events.URI().createFilter(fromBlock='latest')
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    try:
        task = asyncio.gather(iterate_events(contract, contract_filter, 2))
        loop.run_until_complete(task)
    finally:
        loop.close()
