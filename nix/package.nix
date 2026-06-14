{...}: {
  perSystem = {pkgs, ...}: {
    packages.default = pkgs.buildGoModule {
      pname = "agentsync";
      version = "0.1.0";
      src = ../.;
      subPackages = ["apps/agentsync"];
      vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
    };
  };
}
