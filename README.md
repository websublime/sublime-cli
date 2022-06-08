# Workspace

go run main.go workspace --name ws-libs-ui --scope @ws --repo sublime/ws-libs-ui --username miguelramos --email miguel@websublime.dev

# Lib or Pkg

go run main.go --root ws-libs-ui/ create --name utils --type lib --template lit