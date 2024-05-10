package role

import (
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/rcrowley/go-bson"
)

// Role is an interface that represents a role in the system.
type Role interface {
	// Name returns the name of the role.
	Name() string
	// Chat returns the chat format of the role.
	Chat(name string, msg string) string
	// Color returns the given name in the color of the role.
	Color(name string) string
}

type HeirRole interface {
	Inherits() Role
}

var (
	// roles contains all registered Role implementations.
	roles []Role
	// rolesByName contains all registered Role implementations indexed by their name.
	rolesByName = map[string]Role{}
)

// All returns all registered roles.
func All() []Role {
	return roles
}

// Register registers a role to the roles list. The hierarchy of roles is determined by the order of registration.
func Register(role Role) {
	roles = append(roles, role)
	rolesByName[role.Name()] = role
}

// ByName returns the role with the given name. If no role with the given name is registered, the second return value
// is false.
func ByName(name string) (Role, bool) {
	role, ok := rolesByName[name]
	return role, ok
}

// Tier returns the tier of a role based on its registration hierarchy.
func Tier(role Role) int {
	return slices.IndexFunc(roles, func(other Role) bool {
		return role == other
	})
}

// Staff returns if the role is a staff role.
func Staff(role Role) bool {
	return role.Name() == "admin" || role.Name() == "manager" || role.Name() == "mod" || role.Name() == "operator" || role.Name() == "trial"
}

type Roles struct {
	roleMu          sync.Mutex
	roles           []Role
	roleExpirations map[Role]time.Time
}

// NewRoles creates a new Roles instance.
func NewRoles(roles []Role, expirations map[Role]time.Time) *Roles {
	return &Roles{
		roles:           roles,
		roleExpirations: expirations,
	}
}

// Add adds a role to the manager's role list.
func (r *Roles) Add(ro Role) {
	r.checkExpiry()
	r.roleMu.Lock()
	r.roles = append(r.roles, ro)
	r.roleMu.Unlock()
	r.sortRoles()
}

// Remove removes a role from the manager's role list. Users are responsible for updating the highest role usages if
// changed.
func (r *Roles) Remove(ro Role) bool {
	r.roleMu.Lock()
	i := slices.IndexFunc(r.roles, func(other Role) bool {
		return ro == other
	})
	r.roles = slices.Delete(r.roles, i, i+1)
	delete(r.roleExpirations, ro)
	r.roleMu.Unlock()
	r.checkExpiry()
	r.sortRoles()
	return true
}

// Contains returns true if the manager has any of the given roles. Users are responsible for updating the highest role
// usages if changed.
func (r *Roles) Contains(roles ...Role) bool {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	var actualRoles []Role
	for _, ro := range r.roles {
		r.propagateRoles(&actualRoles, ro)
	}

	for _, r := range roles {
		if i := slices.IndexFunc(actualRoles, func(other Role) bool {
			return r == other
		}); i >= 0 {
			return true
		}
	}
	return false
}

// Expiration returns the expiration time for a role. If the role does not expire, the second return value will be false.
func (r *Roles) Expiration(ro Role) (time.Time, bool) {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	e, ok := r.roleExpirations[ro]
	return e, ok
}

// Expire sets the expiration time for a role. If the role does not expire, the second return value will be false.
func (r *Roles) Expire(ro Role, t time.Time) {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	r.roleExpirations[ro] = t
}

// Highest returns the highest role the manager has, in terms of hierarchy.
func (r *Roles) Highest() Role {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	return r.roles[len(r.roles)-1]
}

// All returns the user's roles.
func (r *Roles) All() []Role {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	return append(make([]Role, 0, len(r.roles)), r.roles...)
}

type rolesData struct {
	Roles       []string
	Expirations map[string]time.Time
}

// MarshalBSON ...
func (r *Roles) MarshalBSON() ([]byte, error) {
	var d rolesData
	d.Expirations = make(map[string]time.Time)

	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	for _, rl := range r.roles {
		e, _ := r.roleExpirations[rl]
		if !e.IsZero() {
			d.Expirations[rl.Name()] = e
		}
		d.Roles = append(d.Roles, rl.Name())
	}
	return bson.Marshal(d)
}

// UnmarshalBSON ...
func (r *Roles) UnmarshalBSON(b []byte) error {
	var d rolesData
	if err := bson.Unmarshal(b, &d); err != nil {
		return err
	}

	rls := d.Roles
	for _, rl := range rls {
		ro, ok := ByName(rl)
		if ok {
			r.Add(ro)
			e, ok := d.Expirations[rl]
			if ok {
				r.Expire(ro, e)
			}
		}
	}

	return nil
}

// propagateRoles propagates roles to the user's role list.
func (r *Roles) propagateRoles(actualRoles *[]Role, role Role) {
	*actualRoles = append(*actualRoles, role)
	if h, ok := role.(HeirRole); ok {
		r.propagateRoles(actualRoles, h.Inherits())
	}
}

// sortRoles sorts the roles in the user's role list.
func (r *Roles) sortRoles() {
	sort.SliceStable(r.roles, func(i, j int) bool {
		return Tier(r.roles[i]) < Tier(r.roles[j])
	})
}

// checkExpirations checks each role the user has and removes the expired ones.
func (r *Roles) checkExpiry() {
	r.roleMu.Lock()
	rl, expirations := r.roles, r.roleExpirations
	r.roleMu.Unlock()

	for _, ro := range rl {
		if t, ok := expirations[ro]; ok && time.Now().After(t) {
			r.Remove(ro)
		}
	}
}

func init() {
	Register(Operator{})
	Register(Default{})
	Register(Trial{})
	Register(Mod{})
	Register(Admin{})
	Register(Manager{})
}
