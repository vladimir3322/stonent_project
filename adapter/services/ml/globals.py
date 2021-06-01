from threading import Lock
from nn_image_checker import NNModelChecker
from image_manager import ImageManager

mutex = Lock()

image_checker = NNModelChecker()
image_manager = ImageManager(image_checker)
