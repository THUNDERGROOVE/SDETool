print("Hello SDETool\nLet's open an SDE version")
sde.loadVersion("wl-uf-latest")
t = sde.getTypeByID(364022)
print(t.TypeName)