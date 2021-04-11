This is an utility of image plagiarism detection using deep neural network.

For plagiarism detection we extract features with Convolutional Neural Network — the most popular in computer vision for it's shift invariant structure type of neural network.
In particular, we created a feature extractor by removing the last classification layer from resnet50 trained on image net.
With this feature extractor we get features for rarrible art and build a nmslib index for cosine distance on it. With this index we simply calculate image plagiarism int 3 steps:
1. Get image features.
2. Find the closest images by cosine disatance between features.
3. Check the cosine distance close to 1.

Cosine distance is just a cos of angles between featrues. The closer cos is to 1, the more similar the  images are.
If there is an image in our index with distance close to 1, it means that the image we are checking is a plagiarism.
For better interpretation we scale the distance to [0, 1] where 0 means that the image is new, and 1 means it is plagiarism.
 This value we call score. And this is the score we return.
 
Here are examples of our system output.

![Alt text](images/output_examples.png?raw=true "Title")!