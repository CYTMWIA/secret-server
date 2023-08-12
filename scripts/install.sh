#! /bin/bash
set -e

python3 -m venv venv
source venv/bin/activate
pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt
