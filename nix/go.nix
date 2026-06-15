{...}: {
  perSystem = {
    pkgs,
    config,
    ...
  }: {
    config = {
      shellPackages = with pkgs; [
        go
        gopls
        golangci-lint
        gofumpt
        gotools
      ];

      devShells.default = pkgs.mkShell {
        name = "ponte-dev-env";
        packages = config.shellPackages;
        shellHook = config.shellHooks;
      };
    };
  };
}
