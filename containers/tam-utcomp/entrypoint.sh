#!/bin/bash

# Link against the XAdmin.ini provided
if [[ -e "Shared/XAdmin.ini" ]]; then
  rm -f System/XAdmin.ini
  ln -sf ../Shared/XAdmin.ini System/XAdmin.ini
fi

# Copy configs
cp -av Config/System/*.ini System/

# Run launcher
exec /usr/bin/launcher -ini Config/UT2004.ini -launch Config/launch.yml "$@"
