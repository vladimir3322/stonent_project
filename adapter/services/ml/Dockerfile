FROM python:3.9

RUN apt-get update
RUN apt-get install -y \
        libswscale-dev \
        libtbb2 \
        libtbb-dev \
        libjpeg-dev \
        libpng-dev \
        libtiff-dev \
        libavformat-dev \
        libpq-dev
RUN rm -rf /var/lib/apt/lists/*

RUN pip3 install --upgrade pip
RUN pip3 install pipenv

COPY . /app
WORKDIR /app

RUN pipenv install --system --deploy --ignore-pipfile

CMD ["python", "app.py" ]
