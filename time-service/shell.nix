let
  pkgs = import <nixpkgs>{};
  proto_src = "./proto";
  proto_dst = "./proto";
in
  pkgs.mkShell {
    packages = (with pkgs; [
      protobuf
      protoc-gen-go
    ]);

    shellHook = ''
    regen-protobuf() {
      protoc \
        -I=${proto_src} \
        --go_out=${proto_dst} \
        --go_opt=paths=source_relative \
        ${proto_src}/messages.proto
      }
    '';
  }
