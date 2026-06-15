{...}: {
  perSystem = {pkgs, ...}: {
    packages.default = pkgs.buildGoModule {
      pname = "ponte";
      version = "0.1.0";
      src = ../.;
      subPackages = ["apps/ponte"];
      vendorHash = "sha256-n58Qmiv3gik1qkuXQFbQ+soeOQtUz1dUocEAJepqp/E=";
    };
  };
}
