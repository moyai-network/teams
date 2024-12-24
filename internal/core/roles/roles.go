package roles

import "github.com/bedrock-gophers/role/role"

func Operator() role.Role {
	return role.ByNameMust("operator")
}

func Admin() role.Role {
	return role.ByNameMust("admin")
}

func Default() role.Role {
	return role.ByNameMust("default")
}

func Famous() role.Role {
	return role.ByNameMust("famous")
}

func Khufu() role.Role {
	return role.ByNameMust("khufu")
}

func Manager() role.Role {
	return role.ByNameMust("manager")
}

func Media() role.Role {
	return role.ByNameMust("media")
}

func Menes() role.Role {
	return role.ByNameMust("menes")
}

func Mod() role.Role {
	return role.ByNameMust("mod")
}

func Nitro() role.Role {
	return role.ByNameMust("nitro")
}

func Owner() role.Role {
	return role.ByNameMust("owner")
}

func Partner() role.Role {
	return role.ByNameMust("partner")
}

func Pharaoh() role.Role {
	return role.ByNameMust("pharaoh")
}

func Ramses() role.Role {
	return role.ByNameMust("ramses")
}

func Trial() role.Role {
	return role.ByNameMust("trial")
}

func Voter() role.Role {
	return role.ByNameMust("voter")
}

func Premium(rl role.Role) bool {
	return rl.Tier() >= Khufu().Tier()
}

func Staff(rl role.Role) bool {
	return rl.Tier() >= Trial().Tier()
}
