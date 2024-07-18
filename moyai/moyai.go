package moyai

import (
	"fmt"
	"time"

	"github.com/bedrock-gophers/provider/provider"
	"github.com/google/uuid"

	"github.com/diamondburned/arikawa/v3/state"

	"github.com/moyai-network/teams/moyai/data"

	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sirupsen/logrus"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
)

func init() {
    go func() {
        for Overworld() == nil {
            <-time.After(time.Millisecond)
            continue
        }
        //go tickAirDrop(Overworld())
        go tickAutomaticSave(Overworld())
    }()
}

var (
    // discordState is the discord state of the Moyai bot.
    discordState *state.State
    // srv is the server instance of the Moyai server.
    srv *server.Server
    // playerProvider is the player provider of the Moyai server.
    playerProvider *provider.Provider
    // end is the world of the End dimension.
    end *world.World
    // nether is the world of the Nether dimension.
    nether *world.World
    // deathban is the world of the Deathban arena.
    deathban *world.World

    // lastBlackMarket is the time at which the last black market was opened.
    lastBlackMarket time.Time
    // blackMarketOpened is the time at which the black market was last opened.
    blackMarketOpened time.Time

    // chatGameWord is the word to guess in the chat game.
    chatGameWord string
)

func SetDiscordState(s *state.State) {
    discordState = s
}

func DiscordState() *state.State {
    return discordState
}

func Worlds() []*world.World {
    return []*world.World{Overworld(), Nether(), End(), Deathban()}
}

func Overworld() *world.World {
    if srv == nil {
        return nil
    }
    return srv.World()
}

func End() *world.World {
    return end
}

func Nether() *world.World {
    return nether
}

func Deathban() *world.World {
    return deathban
}

func Players() []*player.Player {
    return srv.Players()
}

func PlayerProvider() *provider.Provider {
    return playerProvider
}

func LoadPlayerData(uuid uuid.UUID) (player.Data, error) {
    dat, err := playerProvider.Load(uuid, func(dimension world.Dimension) *world.World {
        switch dimension {
        case world.Overworld:
            return Overworld()
        case world.Nether:
            return Nether()
        case world.End:
            return End()
        }
        return nil
    })
    return dat, err
}

func wrld(dim world.Dimension) *world.World {
    switch dim {
    case world.Overworld:
        return Overworld()
    case world.Nether:
        return Nether()
    case world.End:
        return End()
    }
    return nil
}

func NewServer(config server.Config) *server.Server {
    providerSettings := provider.DefaultSettings()
    providerSettings.UseServerWorld = false
    providerSettings.World = wrld
    providerSettings.FlushRate = time.Minute * 10
    providerSettings.SaveEffects = false

    playerProvider = provider.NewProvider(&config, providerSettings)
    srv = config.New()
    return srv
}

func LastBlackMarket() time.Time {
    return lastBlackMarket
}

func SetLastBlackMarket(t time.Time) {
    lastBlackMarket = t
}

func BlackMarketOpened() time.Time {
    return blackMarketOpened
}

func SetBlackMarketOpened(t time.Time) {
    blackMarketOpened = t
}

func ChatGameWord() string {
    return chatGameWord
}

func SetChatGameWord(w string) {
    chatGameWord = w
}

func ConfigureDimensions(reg world.EntityRegistry, netherFolder, endFolder string) (*world.World, *world.World) {
    endProv, err := mcdb.Open(endFolder)
    if err != nil {
        panic(err)
    }
    end = world.Config{
        Provider: endProv,
        Dim:      world.End,
        Entities: reg,
    }.New()

    netherProv, err := mcdb.Open(netherFolder)
    if err != nil {
        panic(err)
    }

    nether = world.Config{
        Provider: netherProv,
        Dim:      world.Nether,
        Entities: reg,
    }.New()

    return nether, end
}

func ConfigureDeathban(reg world.EntityRegistry, folder string) *world.World {
    deathbanProv, err := mcdb.Open(folder)
    if err != nil {
        panic(err)
    }
    deathban = world.Config{
        Provider:        deathbanProv,
        Dim:             world.Overworld,
        Entities:        reg,
        RandomTickSpeed: -1,
    }.New()

    return deathban
}

func Close() {
    data.FlushCache()
    destroyAirDrop(srv.World(), lastDropPos)
    for _, p := range Players() {
        u, err := data.LoadUserFromName(p.Name())
        if err != nil {
            continue
        }
        h, ok := p.Handler().(userHandler)
        if !ok {
            continue
        }
        u.PlayTime += h.LogTime()
        data.SaveUser(u)

        p.Disconnect("Server is shutting down.")
    }

    time.Sleep(time.Millisecond * 500)

    sotw.Save()
    _ = Nether().Close()
    _ = End().Close()
    _ = srv.World().Close()
    if err := srv.Close(); err != nil {
        logrus.Fatalln("close server: %v", err)
    }
}

func tickAutomaticSave(w *world.World) {
    for {
        <-time.After(time.Minute * 1)
        w.Save()
        for _, p := range Players() {
            u, err := data.LoadUserFromName(p.Name())
            if err != nil {
                fmt.Println("load user: ", err)
                continue
            }
            if u.StaffMode {
                continue
            }

            err = playerProvider.SavePlayer(p)
            if err != nil {
                fmt.Printf("save player: %v\n", err)
            }
        }
    }
}

type userHandler interface {
    player.Handler
    LogTime() time.Duration
}
