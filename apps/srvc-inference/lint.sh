mypy .
ruff check . --fix
bandit -c pyproject.toml -r .
black .
