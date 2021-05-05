import os
from os.path import join
from py_dotenv import read_dotenv


class Config:
    DATA_PATH = './data'
    SOURCE_PATH = './source'
    CONTRACT_ADDRESS = '0xd07dc4262BCDbf85190C01c996b4C06a461d2430'
    CONTRACT_ABI_FILENAME = './contract/abi.json'
    # FIRST_BLOCK_NUMBER = 10147631
    FIRST_BLOCK_NUMBER = 12214955
    SAFE_BLOCK_STEP = 400
    LAST_URLS_COUNT = 10

    def __init__(self):
        pass

    def get_web3_provider(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        provider_url = os.getenv('WEB_3_PROVIDER')

        return provider_url

    def get_ipfs_url(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        ipfs_url = os.getenv('IPFS_URL')

        return ipfs_url

    def get_max_count_NFTs_loader(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        max_count_NFTs_loader = os.getenv('MAX_COUNT_NFTS_LOADER')

        return None if not max_count_NFTs_loader else int(max_count_NFTs_loader)

    def get_max_count_NFTs_watcher(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        max_count_NFTs_watcher = os.getenv('MAX_COUNT_NFTS_WATCHER')

        return None if not max_count_NFTs_watcher else int(max_count_NFTs_watcher)

    def get_save_images(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        save_images = os.getenv('SAVE_IMAGES')

        return True if save_images else False


config = Config()
