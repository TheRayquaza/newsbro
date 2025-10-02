{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    utils,
  }:
    utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShell = let
          python = pkgs.python313.withPackages (ppkgs:
            with ppkgs; [
              python-lsp-server
              uv
              pydantic
              fastapi-cli
              fastapi
              joblib
              streamlit
              pandas
              scikit-learn
              ruff
              validators
            ]);
        in
          pkgs.mkShell {
            buildInputs = [
              python
            ];
            UV_PYTHON = python;
            UV_PYTHON_PREFERENCE="only-system";
          };
      }
    );
}