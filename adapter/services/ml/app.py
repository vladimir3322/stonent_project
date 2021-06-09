from threading import Thread
from globals import image_manager, mutex
from PIL import Image
import server
import config


if __name__ == '__main__':
    print('Start ML', flush=True)

    Image.MAX_IMAGE_PIXELS = None

    open(config.registered_images_file, 'w').close()
    open(config.rejected_images_by_IPFS_file, 'w').close()
    open(config.rejected_images_by_NN_file, 'w').close()

    registerer_thread = Thread(target=image_manager.register_new_images, args=[mutex])
    server_thread = Thread(target=server.run_server)

    registerer_thread.start()
    server_thread.start()
