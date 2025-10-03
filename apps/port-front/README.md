# Port-Front

### Dev

For Nix users:

```bash
nix develop .
uv venv
source ./.venv/bin/activate
uv pip install -r requirements.txt --python .venv/bin/python
```

For normal users:

```bash
python3 -m venv .venv
pip install -r
streamlit run app.py
```

## Tools

```bash
ruff check --fix
ty check
```

## Variables

| Name                 | Value                                 |
| -------------------- | ------------------------------------- |
| ACCOUNT_BASE_URL     | `https://account.newsbro.cc`          |
| ARTICLE_BASE_URL     | `https://article.newsbro.cc`          |
