{inputs, ...}: {
  perSystem = {
    pkgs,
    config,
    ...
  }: let
    python = pkgs.python314;

    workspace = inputs.uv2nix.lib.workspace.loadWorkspace {workspaceRoot = ../.;};

    overlay = workspace.mkPyprojectOverlay {
      sourcePreference = "wheel";
    };

    pythonSet =
      (pkgs.callPackage inputs.pyproject-nix.build.packages {
        inherit python;
      })
      .overrideScope (
        pkgs.lib.composeManyExtensions [
          inputs.pyproject-build-systems.overlays.default
          overlay
        ]
      );
  in {
    config = {
      shellPackages = [python pkgs.uv];

      shellHooks = ''
        if [ -f "pyproject.toml" ]; then
          uv sync --frozen 2>/dev/null || uv sync
        fi
      '';

      devShells.default = pkgs.mkShell {
        name = "agentsync-dev-env";
        packages = config.shellPackages;
        shellHook = config.shellHooks;
      };

      packages.default = pythonSet.mkVirtualEnv "agentsync-env" workspace.deps.default;
    };
  };
}
