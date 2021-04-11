# Stonent project

This project has been created to help NFT-painters fight with intellectual property theft.
We use neural networks to detect plagiarism of paintings.
This project is hosted by chainlink external adapter.
Every painting is certificated by oracles.
Originality score of each painting is kept by the [contract](https://rinkeby.etherscan.io/tx/0x1d3e80be8475b53fcef3a71e65d3a24e14fae47f4dd32e1079cb666585ec358b).
So, everybody can see a picture and check its originality score!
There is [the project demo](http://stonent.us-east-2.elasticbeanstalk.com/) accessed by Rinkeby test network.

Our solution for plagiarism detection is inspired by Content-Based Image Retrieval or CBIR.
 CBIR is a task of image searching. There are different ways to solve this problem.
  But all of them are about building descriptions of images. So the problem solving consists of 3 steps.
1. Training feature extractor. (We took pre-trained)
2. Indexing dataset. (We made it via nmslib)
3. Searching the best matches (Also via nmslib)

For our solution we took resnet50 model pre-trained for ImageNet classification problem.
 While ImageNet is not an CBIR problem, features in deep model layers are enough meaningful to achieve appreciable result.
 So we took resnet50, removed the classification layer from the model and by this made the model feature extractor.
 With this feature extractor we created features for images from rarrible and build a index on it.
 To compare features with each other, we use cosine distance. Cosine distance is a classical metric for neural network features.
 If features are vectors in some space, then cosine distance is cos of angle between this vectors.
  So, if distance is close to 1, images are similar, and if they are close to -1, they are different.
 
After we found most simular image we need to understand if it is similar enough to call image plagiarism. 
 For this task we took dataset of classical art, and build a simple classificator.
  Classificator was trained on distances between different images and between image and it's augmented (slightly changed) copy.

The result output of our system is the output of this classificator. Here you can see example of the results. 


![Alt text](images/output_examples.png?raw=true "Title")


## Run adapter:
1. Go to the adapter dir: `cd stonent/adapter`
2. Set your Web3 http provider as environment variable: `WEB_3_PROVIDER`
3. Run docker: `docker-compose up`
4. ...
5. Profit!
