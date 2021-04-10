import os

from PIL import Image

from contract import config


class ImageManager:
    descriptions_path = config.Config.DATA_PATH
    images_path = config.Config.SOURCE_PATH

    def __init__(self, image_cheker, delete_registered=True):
        self.image_cheker = image_cheker
        self._registred_ids = set()
        self.delete_registered = delete_registered

    def _register_new_image(self, id):
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