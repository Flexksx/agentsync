import typer


app = typer.Typer(name="agentsync", help="Sync skills and instructions across AI agent vendors.")


def main() -> None:
    app()
