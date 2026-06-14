{...}: {
  perSystem = {
    pkgs,
    lib,
    ...
  }: {
    options.shellPackages = lib.mkOption {
      type = lib.types.listOf lib.types.package;
      default = [];
    };

    options.shellHooks = lib.mkOption {
      type = lib.types.lines;
      default = "";
    };

    config.shellPackages = with pkgs; [
      just
      alejandra
      lefthook
    ];
  };
}
