This is an utility of image plagiarism detection using ResNet50 deep neural network.


For plagiarism detection we extract features with Convolutional Neural Network and then calculate cosine distance between this features. While the distance of two features will lie in [-1, 1], we move them into [0, 1] where 0 means images are very simular, and 1 means images are different. As a result, for new image, we return minimum over all scores between this image and all images in data base.  