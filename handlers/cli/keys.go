package cli

//nolint:lll // cannot break struct tags
type KeysCliHandler struct {
	List   ListKeysCliHandler  `cmd:"" name:"list" help:"list identities"`
	Import ImportCliHandler    `cmd:"" name:"import" help:"import an identity - currently only x25519 identities are supported"`
	Remove RemoveKeyCliHandler `cmd:"" name:"remove" aliases:"rm,del" help:"remove an identity"`
}
