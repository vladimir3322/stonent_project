import asyncio
import os
from functools import wraps, partial

from . import config
from . import download_image_data as did
from . import get_contract as gc


metadata = {
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


async def handle_event(event):
    image_id = event["args"]["_id"]
    image_url = f'{config.config.get_ipfs_url()}{event["args"]["_value"]}'

    print(f'New event: {image_id}')

    metadata['found_images'] += 1

    download_error = await did.download_image_data(image_url, image_id)

    if not download_error:
        metadata['downloaded_images'] += 1


async def iterate_events(event_filter, poll_interval):
    while True:
        for event in event_filter.get_new_entries():
            await handle_event(event)

            max_count_NFTs_watcher = config.config.get_max_count_NFTs_watcher()

            if not max_count_NFTs_watcher:
                continue
            if metadata['downloaded_images'] >= max_count_NFTs_watcher:
                return

        await asyncio.sleep(poll_interval)


def listen_images():
    if not os.path.exists(config.config.DATA_PATH):
        os.mkdir(config.config.DATA_PATH)
    if not os.path.exists(config.config.SOURCE_PATH):
        os.mkdir(config.config.SOURCE_PATH)

    contract = gc.get_contract()
    contract_filter = contract.events.URI().createFilter(fromBlock='latest')
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    try:
        print('Images listener started!')
        task = asyncio.gather(iterate_events(contract_filter, 2))
        loop.run_until_complete(task)
    finally:
        metadata['found_images'] = 0
        loop.close()
