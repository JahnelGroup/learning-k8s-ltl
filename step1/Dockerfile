FROM python:3.9-slim-buster

WORKDIR /app

COPY requirements.txt .
COPY . .

RUN pip install -r requirements.txt


CMD ["python", "main.py"]
