FROM python:3.8.0

RUN apt-get update \
    && apt-get install -y \
        build-essential \
        cmake \
        git \
        wget \
        unzip \
        yasm \
        pkg-config \
        libswscale-dev \
        libtbb2 \
        libtbb-dev \
        libjpeg-dev \
        libpng-dev \
        libtiff-dev \
        libavformat-dev \
        libpq-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install --upgrade pip
RUN pip3 install pipenv
RUN pip3 install opencv-python imagehash torch torchvision
RUN pip3 install nmslib
RUN pip3 install opencv-python imagehash

COPY . /app
WORKDIR /app
RUN pipenv install --system --deploy --ignore-pipfile
RUN pip3 install py_dotenv
RUN pip3 install -r requirements.txt

CMD ["python", "app.py" ]
