{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  # Packages that you want to include in your shell environment
  buildInputs = [
    pkgs.go           # The Go compiler and tools
    pkgs.python312
    pkgs.git          # Git, in case you need version control
    pkgs.air

    pkgs.python312Packages.requests
  ];

  # Set the GO111MODULE environment variable (optional)
  GO111MODULE = "on";

  # Set additional environment variables if necessary
  shellHook = ''
    export GOPATH=$HOME/go
    echo "HTTP Server and SH Script"
    '';
}

