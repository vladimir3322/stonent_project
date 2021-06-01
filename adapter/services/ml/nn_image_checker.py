import torch
import nmslib
import numpy as np

from torch import nn
from PIL import Image
from torchvision import transforms
from scipy.stats import logistic


class NNModelChecker:
    def __init__(self):
        """
        We will use renset50 trained on ImageNet as feature extractor.
        To get features we remove last classification layer of the nn.
        To find nearest features we use nmslib index.
        """
        model = torch.hub.load('pytorch/vision:v0.9.0', 'resnet50', pretrained=True)
        self.feature_extractor = nn.Sequential(*list(model._modules.values())[:-1])

        self.preprocess = transforms.Compose([
                          transforms.Resize(256),
                          transforms.CenterCrop(224),
                          transforms.ToTensor(),
                          transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
                        ])

        self._index_need_to_be_build = True

        self._index = nmslib.init(method='hnsw', space='cosinesimil')
        self._feature_dict = {}

    def _get_features(self, pil_image):
        """
        :param pil_image: image loaded py PIL library.
        :return: array of features for the image
        """
        with torch.no_grad():
            image = np.array(pil_image)
            if image.ndim == 2:
                image = image[..., None]
                image = np.concatenate([image, image, image], -1)
            if image.shape[-1] == 4:
                image = image[..., :3]
            input_image = self.preprocess(Image.fromarray(image))
            return self.feature_extractor(input_image[None, :])[0].reshape(-1)

    def _transform_scores(self, scores):
        """
        :param scores: raw cosine distance scores
        :return: scores scaled to [0, 1]

        mean was calculated on classical art dataset
        temp was choose to make sigmiod output close to 0 or 1
        """
        mean = 0.0037
        temp = 10000
        return logistic.cdf((scores - mean) * temp)

    @staticmethod
    def cosine_distance(input1, input2):
        """
        :param input1: first feature vector
        :param input2: second feature vector
        :return: cosine distance between vectors.
        """
        return np.dot(input1, input2.T) / np.sqrt(np.dot(input1, input1.T) * np.dot(input2, input2.T))

    def add_image_to_storage(self, pil_image, description):
        """
        :param pil_image: image loaded py PIL library.
        :param description: description of the image. Will be returned if image will be chosen as neighbour
        :return: None
        """
        features = self._get_features(pil_image)
        index = len(self._feature_dict)
        self._feature_dict[index] = description
        self._index.addDataPoint(data=features, id=index)
        self._index_need_to_be_build = True

    def find_most_similar_images(self, pil_image, num=5):
        """
        :param pil_image:  image loaded py PIL library.
        :param num: number of neighbours to return
        :return: scores, nearest_descriptions.
                 scores - Scores of simularity between pil_image and neighbours
                 nearest_descriptions — descriptions of neighbours
        """
        if self._index_need_to_be_build:
            self._index.createIndex({'post': 2})
            self._index_need_to_be_build = False
        features = self._get_features(pil_image)
        indexes, scores = self._index.knnQuery(features, k=num)
        nearest_descriptions = [self._feature_dict[index] for index in indexes]
        scores = self._transform_scores(scores)
        return scores, nearest_descriptions
