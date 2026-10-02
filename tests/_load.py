"""把 skills/<x>/scripts/*.py 當模組載入。腳本不是套件，用路徑載。"""
import importlib.util
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load(rel: str):
    path = ROOT / rel
    spec = importlib.util.spec_from_file_location(path.stem, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module
