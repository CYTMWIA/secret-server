#! /bin/bash
set -e

source venv/bin/activate

gunicorn -w 4 -b '0.0.0.0' 'secret_server:create_app()'
