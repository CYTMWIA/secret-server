import hashlib
import os
import tempfile
import urllib.parse

from cryptography.hazmat.primitives.ciphers.aead import AESGCM
from flask import Flask, abort, request, send_file
from sqlitedict import SqliteDict

from . import config


def parse_query_string(qs: bytes):
    raw = urllib.parse.parse_qs(urllib.parse.unquote(request.query_string))
    for k in raw:
        if isinstance(raw[k], list) and len(raw[k]) == 1:
            raw[k] = raw[k][0]
    return raw


def sha3_256(s: str):
    encode = s.strip().encode("utf-8")
    return hashlib.sha3_256(encode).digest().hex()


def generate_key(password: str, bits=256):
    if bits % 8:
        raise ValueError("`bits` mod 8 MUST equal to 0")
    key_len = bits
    password = password.encode("utf-8")
    key = (password * (int(key_len / len(password)) + 1))[: key_len // 8]
    return key


def encrypt_data(data: bytes, associated_data: str, password: str):
    key = generate_key(password)
    aesgcm = AESGCM(key)

    iv = os.urandom(96)
    data = aesgcm.encrypt(iv, data, associated_data.encode("utf-8"))

    return data, iv


def decrypt_data(data: bytes, associated_data: str, password: str, iv: bytes):
    key = generate_key(password)
    aesgcm = AESGCM(key)

    data = aesgcm.decrypt(iv, data, associated_data.encode("utf-8"))

    return data


def remove_leading_slash(s: str):
    return s.lstrip("/\\")


class Core:
    def __init__(self, cfg: config.Config) -> None:
        print("Working Dir: ", os.getcwd())

        self.cfg = cfg

        self.root_dir = os.path.abspath(os.path.join(os.getcwd(), "root"))
        os.makedirs(self.root_dir, exist_ok=True)
        self.db = SqliteDict(cfg.sqlite_path, autocommit=True)

    def has_permissions(self, api_key: str):
        if len(self.cfg.api_key) == 0:
            return True
        return sha3_256(api_key) in self.cfg.api_key

    def get_real_path(self, client_path):
        return os.path.join(self.root_dir, client_path)

    def upload_file(self, path):
        path = remove_leading_slash(path)

        query = parse_query_string(request.query_string)
        if not self.has_permissions(query.get("api_key", None)):
            abort(403)
        print(query)
        real_path = self.get_real_path(path)
        data = request.data
        if "file_key" in query:
            data, iv = encrypt_data(data, path, query["file_key"])
            self.db[path] = iv
        with open(real_path, "wb") as f:
            f.write(data)

        return ""

    def download_file(self, path):
        path = remove_leading_slash(path)

        query = parse_query_string(request.query_string)
        real_path = self.get_real_path(path)

        not_exists = [
            not os.path.exists(real_path),
            (("file_key" in query) and (path not in self.db)),
        ]
        if any(not_exists):
            abort(404)

        if "file_key" in query:
            with open(real_path, "rb") as f:
                data = f.read()
            plain_data = decrypt_data(data, path, query["file_key"], self.db[path])
            with tempfile.NamedTemporaryFile() as tmp:
                tmp.write(plain_data)
                tmp.seek(0)
                return send_file(tmp.name, download_name=os.path.basename(path))
        else:
            return send_file(real_path, download_name=os.path.basename(path))


def create_app():
    cfg = config.read(
        ["config.json", "config/config.json", "config/config-example.json"]
    )
    print(cfg)

    core = Core(cfg)

    app = Flask(__name__)
    app.route("/<path:path>", methods=["PUT"])(core.upload_file)
    app.route("/<path:path>", methods=["GET"])(core.download_file)

    return app


if __name__ == "__main__":
    print(("-" * 20) + " Debug Mode " + ("-" * 20))
    create_app().run("0.0.0.0", 16969, debug=True)
