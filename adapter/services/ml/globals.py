from threading import Lock
from nn_image_checker import NNImageChecker
from image_manager import ImageManager

mutex = Lock()

image_checker = NNImageChecker()
image_manager = ImageManager(image_checker)
all_images_has_been_downloaded = False
