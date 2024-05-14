#!/bin/zsh

# go
go mod tidy

# python

cd interactive
python3 -m venv venv
source venv/bin/activate
python3 -m pip install -r requirements.txt
cd ..

# other tools
pre-commit


if [[ "${CODESPACES}" == true ]]; then
  echo "Fixing directory ownership for GitHub Codespaces..." >&2
  sudo chown -R nonroot:nonroot /home/nonroot
  sudo chown -R nonroot:nonroot /workspace
fi
