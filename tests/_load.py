"""把 skills/<x>/scripts/*.py 當模組載入。腳本不是套件，用路徑載。"""
import importlib.util
import sys
from pathlib import Path

# 載入 skills/<x>/scripts/ 底下的腳本時不要寫 __pycache__：audit_skill.py 會把
# 任何沒被 SKILL.md 指路的檔案當成孤兒報錯，而 .pyc 永遠不會被指路。
sys.dont_write_bytecode = True

ROOT = Path(__file__).resolve().parents[1]


def load(rel: str):
    path = ROOT / rel
    spec = importlib.util.spec_from_file_location(path.stem, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module
