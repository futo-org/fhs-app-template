data "external_schema" "bun" {
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-bun",
    "load",
    "--path", "./internal/database/entity",
    "--dialect", "sqlite",
  ]
}

env "bun" {
  src = data.external_schema.bun.url
  dev = "sqlite://dev?mode=memory"

  migration {
    dir = "file://internal/database/migrations?format=goose"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
