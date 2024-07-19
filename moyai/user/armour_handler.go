package user

import (
    "github.com/df-mc/atomic"
    "github.com/df-mc/dragonfly/server/event"
    "github.com/df-mc/dragonfly/server/item"
    "github.com/df-mc/dragonfly/server/player"
    "time"
)

type ArmourHandler struct {
    stormBreakerHelmet item.Stack
    stormBreakerStatus atomic.Bool
    stormBreakerCancel chan struct{}
    p                  *player.Player
}

func NewArmourHandler(p *player.Player) *ArmourHandler {
    return &ArmourHandler{p: p}
}

func (a *ArmourHandler) HandleTake(ctx *event.Context, slot int, itm item.Stack) {
    _, ok := itm.Value("storm_breaker")
    if ok {
        ctx.Cancel()
        return
    }

    h, ok := a.p.Handler().(*Handler)
    if ok {
        sortClassEffects(h)
        sortArmourEffects(h)
    }
}
func (a *ArmourHandler) HandlePlace(ctx *event.Context, slot int, itm item.Stack) {
    _, ok := itm.Value("storm_breaker")
    if ok {
        ctx.Cancel()
        return
    }

    h, ok := a.p.Handler().(*Handler)
    if ok {
        sortClassEffects(h)
        sortArmourEffects(h)
    }
}
func (a *ArmourHandler) HandleDrop(ctx *event.Context, slot int, itm item.Stack) {
    _, ok := itm.Value("storm_breaker")
    if ok {
        ctx.Cancel()
        return
    }

    h, ok := a.p.Handler().(*Handler)
    if ok {
        sortClassEffects(h)
        sortArmourEffects(h)
    }
}

func (a *ArmourHandler) stormBreak() {
    if a.stormBreakerStatus.Load() {
        a.stormBreakerCancel <- struct{}{}
        a.handleStormBreakProcess()
        return
    }

    a.handleStormBreakProcess()
    a.stormBreakerStatus.Store(true)
    a.stormBreakerHelmet = a.p.Armour().Helmet()
    a.p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: item.ColourBrown().RGBA()}}, 1).WithValue("storm_breaker", true))
}

func (a *ArmourHandler) stormBreakCancel() {
    if a.stormBreakerStatus.Load() {
        a.stormBreakerStatus.Store(false)
        a.stormBreakerCancel <- struct{}{}
        a.p.Armour().SetHelmet(a.stormBreakerHelmet)
    }
}

func (a *ArmourHandler) handleStormBreakProcess() {
    a.stormBreakerCancel = make(chan struct{})
    go func() {
        select {
        case <-time.After(time.Second * 5):
            a.stormBreakerStatus.Store(false)
            a.p.Armour().SetHelmet(a.stormBreakerHelmet)
        case <-a.stormBreakerCancel:
            return
        }
    }()
}
