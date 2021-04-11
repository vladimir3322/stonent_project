import os

from PIL import Image

from contract import config


class ImageManager:
    descriptions_path = config.Config.DATA_PATH
    images_path = config.Config.SOURCE_PATH

    def __init__(self, image_cheker, delete_registered=True):
        """
        :param image_cheker: image checker class to register new images
        :param delete_registered: flag to delete images after registration
        This class scan download directory and register all new images in image_cheker.
        """
        self.image_cheker = image_cheker
        self._registred_ids = set()
        self.delete_registered = delete_registered

    def _register_new_image(self, id):
        """
        :param id: image id. id is name of the files
        :return: None
        """
        image_path = os.path.join(self.images_path, id)
        description_path = os.path.join(self.descriptions_path, id)
        try:
            pil_image = Image.open(image_path)
            decription = str(id)

            self.image_cheker.add_image_to_storage(pil_image, decription)
        except Exception as e:
            print("error in registring new image", e)

        if self.delete_registered:
            os.remove(image_path)
            os.remove(description_path)
        else:
            self._registred_ids.add(id)

    def register_new_images(self):
        """
        :return: None
        This function will scan downloads folder and register all new images in it.
        """
        if not os.path.exists(self.descriptions_path) or not os.path.exists(self.images_path):
            return

        descriptions_ids = os.listdir(self.descriptions_path)
        images_ids = os.listdir(self.images_path)
        for image_id in images_ids:
            if image_id in self._registred_ids:
                continue
            if image_id not in descriptions_ids:
                continue
            self._register_new_image(image_id)