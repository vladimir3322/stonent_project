from web3 import Web3
import json
import asyncio
import aiohttp
import aiofile
from urllib.parse import urlparse
import os
import time
from functools import wraps, partial

from . import config


def async_wrap(func):
    @wraps(func)
    async def run(*args, loop=None, executor=None, **kwargs):
        if loop is None:
            loop = asyncio.get_event_loop()
        pfunc = partial(func, *args, **kwargs)
        return await loop.run_in_executor(executor, pfunc)
    return run


def show_progress():
    print('Active tasks count: ', len(
        [task for task in asyncio.Task.all_tasks() if not task.done()])
    )


async def download_image_data(contract, image_id):
    data_ipfs_url = contract.functions.uri(image_id).call()
    parsed_data_ipfs_url = urlparse(data_ipfs_url)

    if not parsed_data_ipfs_url.path or parsed_data_ipfs_url.path == '/':
        return False

    try:
        async with aiohttp.ClientSession() as request:
            async with request.get(data_ipfs_url) as data_response:
                if data_response.status != 200:
                    return False

                data = await data_response.json()
    except Exception as e:
        print(e)
        return False

    image_source_ipfs_link = urlparse(data['image'])
    image_source_path = image_source_ipfs_link.path

    if not image_source_path or image_source_path == '/':
        return False

    ipfs_prefix = '/ipfs'

    if image_source_path[0:len(ipfs_prefix)] != ipfs_prefix:
        image_source_path = ipfs_prefix + image_source_path

    image_source_ipfs_url = f'{config.config.IPFS_IMAGE_SOURCE_URL}{image_source_path}'

    try:
        async with aiohttp.ClientSession() as request:
            async with request.get(image_source_ipfs_url) as data_response:
                if data_response.status != 200:
                    return False
                image_source = await data_response.read()
    except Exception as e:
        print(e)
        return False

    try:
        async with aiofile.async_open(f'{config.config.SOURCE_PATH}/{image_id}', 'wb', encoding='utf-8') as file:
            await file.write(image_source)
    except Exception as e:
        print(e)
        return False

    data['image_source_path'] = image_source_path

    try:
        async with aiofile.async_open(f'{config.config.DATA_PATH}/{image_id}', 'w', encoding='utf-8') as file:
            await file.write(json.dumps(data))
    except Exception as e:
        print(e)
        return False

    return True


async def run_blocks_interval(contract, block_from, block_to):
    show_progress()

    contract_filter = contract.events.URI().createFilter(fromBlock=block_from, toBlock=block_to)
    events = await async_wrap(contract_filter.get_all_entries)()

    for event in events:
        await download_image_data(contract, event["args"]["_id"])

    show_progress()


def download_images():
    if not os.path.exists(config.config.DATA_PATH):
        os.mkdir(config.config.DATA_PATH)
    if not os.path.exists(config.config.SOURCE_PATH):
        os.mkdir(config.config.SOURCE_PATH)

    w3 = Web3(Web3.HTTPProvider(config.config.get_web3_provider()))
    with open(config.config.CONTRACT_ABI_FILENAME) as json_file:
        abi = json.load(json_file)
    contract = w3.eth.contract(address=config.config.CONTRACT_ADDRESS, abi=abi)


    last_block = w3.eth.getBlock('latest')
    last_block_number = last_block['number']

    start = time.time()
    print(config.config.FIRST_BLOCK_NUMBER, last_block_number)
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    futures = [run_blocks_interval(contract, i, i + config.config.SAFE_BLOCK_STEP) for i in range(config.config.FIRST_BLOCK_NUMBER, last_block_number, config.config.SAFE_BLOCK_STEP)]
    loop.run_until_complete(asyncio.wait(futures))

    end = time.time()

    print('Downloaded images successfully', end - start)
