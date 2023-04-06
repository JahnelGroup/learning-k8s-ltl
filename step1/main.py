from flask import Flask
from os import environ

app = Flask(__name__)


@app.route('/')
def hello_world():
    return 'Hello, World!'


if __name__ == '__main__':
    app.run(debug=bool(environ.get("DEBUG", "false")),
            port=8888, host='0.0.0.0')
