import os

def get_version():
    current_version = '0.0.0'
    if 'version' in os.listdir():
        with open('version', 'r', encoding='utf-8') as current_version_file:
            current_version = current_version_file.readline().strip()
    return current_version
