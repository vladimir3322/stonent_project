from web3 import Web3
import asyncio
import os
import time
from functools import wraps, partial

from . import config
from . import download_image_data as did
from . import get_contract as gc


metadata = {
    'active_tasks': 0,
    'found_images': 0,
    'downloaded_images': 0,
}


def async_wrap(func):
    @wraps(func)
    async def run(*args, loop=None, executor=None, **kwargs):
        if loop is None:
            loop = asyncio.get_event_loop()
        pfunc = partial(func, *args, **kwargs)
        return await loop.run_in_executor(executor, pfunc)
    return run


async def run_blocks_interval(loop, contract, block_from, block_to):
    contract_filter = contract.events.URI().createFilter(fromBlock=block_from, toBlock=block_to)
    events = await async_wrap(contract_filter.get_all_entries)()

    for event in events:
        metadata['found_images'] += 1

        image_url = f'{config.config.get_ipfs_url()}{event["args"]["_value"]}'
        download_error = await did.download_image_data(image_url, event["args"]["_id"])
        max_count_NFTs_loader = config.config.get_max_count_NFTs_loader()
        if download_error:
            print(download_error)
            continue

        metadata['downloaded_images'] += 1
        print(f'Downloaded images: {metadata["downloaded_images"]}')

        if not max_count_NFTs_loader:
            continue

        if metadata['downloaded_images'] >= max_count_NFTs_loader:
            loop.stop()
            loop.close()
            return

    metadata['active_tasks'] -= 1
    print(f'Active tasks: {metadata["active_tasks"]}')


def download_images():
    if not os.path.exists(config.config.DATA_PATH):
        os.mkdir(config.config.DATA_PATH)
    if not os.path.exists(config.config.SOURCE_PATH):
        os.mkdir(config.config.SOURCE_PATH)

    w3 = Web3(Web3.HTTPProvider(config.config.get_web3_provider()))

    contract = gc.get_contract()

    last_block = w3.eth.getBlock('latest')
    last_block_number = last_block['number']

    start = time.time()
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    futures = [run_blocks_interval(loop, contract, i, i + config.config.SAFE_BLOCK_STEP) for i in range(config.config.FIRST_BLOCK_NUMBER, last_block_number, config.config.SAFE_BLOCK_STEP)]
    metadata['active_tasks'] = len(futures)

    try:
        print('Images downloader started!')
        loop.run_until_complete(asyncio.wait(futures))
    except:
        pass

    metadata['found_images'] = 0
    end = time.time()

    print('Downloaded images successfully', end - start)
