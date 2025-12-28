import os
import json
import logging

from datetime import datetime

logger = logging.getLogger(__name__)

def load_json(filename):
    with open(filename, 'r') as f:
        return json.load(f)

def dump_json(data, filename):
    with open(filename, 'w') as f:
        json.dump(data, f, indent=4)

def get_current_datetime():
    return datetime.now().strftime('%Y-%m-%d %H:%M:%S')

def create_dir_if_not_exists(dir_path):
    if not os.path.exists(dir_path):
        os.makedirs(dir_path)