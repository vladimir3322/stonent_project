import aiohttp
import aiofile
from urllib.parse import urlparse
import json

from . import config


metadata = {
    'last_urls': [],
}

errors = {
    'not_found_by_contract': 'not_found_by_contract',
    'failed_metadata_request': 'failed_metadata_request',
    'invalid_image_source_url': 'invalid_image_source_url',
    'failed_image_source_request': 'failed_image_source_request',
    'failed_metadata_saving': 'failed_metadata_saving',
    'failed_image_source_saving': 'failed_image_source_saving',
}


async def download_image_data(data_ipfs_url, image_id, save_to_disk=True):
    parsed_data_ipfs_url = urlparse(data_ipfs_url)

    if not parsed_data_ipfs_url.path or parsed_data_ipfs_url.path == '/':
        return errors['not_found_by_contract']

    metadata['last_urls'].insert(0, data_ipfs_url)
    metadata['last_urls'] = metadata['last_urls'][0:config.config.LAST_URLS_COUNT]

    try:
        async with aiohttp.ClientSession() as request:
            async with request.get(data_ipfs_url) as data_response:
                if data_response.status != 200:
                    print(f'Metadata request response code: {data_response.status}')
                    print(f'For url: {data_ipfs_url}')

                    try:
                        data = await data_response.text()
                        print(data)
                    except:
                        pass

                    return errors['failed_metadata_request']

                data = await data_response.json()
    except Exception as e:
        print(e)
        return errors['failed_metadata_request']

    image_source_ipfs_link = urlparse(data['image'])
    image_source_path = image_source_ipfs_link.path

    if not image_source_path or image_source_path == '/':
        return errors['invalid_image_source_url']

    ipfs_prefix = '/ipfs'

    if image_source_path[0:len(ipfs_prefix)] != ipfs_prefix:
        image_source_path = ipfs_prefix + image_source_path

    image_source_ipfs_url = f'{config.config.get_ipfs_url()}{image_source_path}'

    try:
        async with aiohttp.ClientSession() as request:
            async with request.get(image_source_ipfs_url) as data_response:
                if data_response.status != 200:
                    print(f'Source request response code: {data_response.status}')
                    print(f'For url: {image_source_ipfs_url}')

                    try:
                        data = await data_response.text()
                        print(data)
                    except:
                        pass

                    return errors['failed_image_source_request']
                image_source = await data_response.read()
    except Exception as e:
        print(e)
        return errors['failed_image_source_request']

    if not config.config.get_save_images():
        return True

    data['image_source_path'] = image_source_path

    if save_to_disk:
        try:
            async with aiofile.async_open(f'{config.config.SOURCE_PATH}/{image_id}', 'wb', encoding='utf-8') as file:
                await file.write(image_source)
        except Exception as e:
            print(e)
            return errors['failed_metadata_saving']

        try:
            async with aiofile.async_open(f'{config.config.DATA_PATH}/{image_id}', 'w', encoding='utf-8') as file:
                await file.write(json.dumps(data))
        except Exception as e:
            print(e)
            return errors['failed_image_source_saving']

        return None
    else:
        return image_source, json.dumps(data)

