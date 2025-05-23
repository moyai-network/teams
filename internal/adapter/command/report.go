package command

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"time"

	"github.com/moyai-network/teams/pkg/lang"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

// Report is a command used to report other players.
type Report struct {
	Targets []cmd.Target `cmd:"target"`
	Reason  reason       `cmd:"reason"`
}

// Run ...
func (r Report) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	l := locale(s)
	p := s.(*player.Player)
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	if len(r.Targets) < 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	t, ok := r.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if s == t {
		o.Error(lang.Translatef(l, "command.report.self"))
		return
	}
	if u.Teams.Report.Active() {
		o.Error(lang.Translatef(l, "command.report.cooldown", u.Teams.Report.Remaining().Round(time.Millisecond*10)))
		return
	}
	u.Teams.Report.Set(time.Minute)
	core.UserRepository.Save(u)

	internal.Messagef(p, "command.report.success")
	internal.Alertf(tx, s, "staff.alert.report", t.Name(), r.Reason)
	// webhook.Send(webhook.Report, hook.Webhook{
	// 	Embeds: []hook.Embed{{
	// 		Title: "Report (Practice)",
	// 		Color: 0xFFFFFF,
	// 		Description: strings.Join([]string{
	// 			fmt.Sprintf("**Player:** %v", t.Name()),
	//             fmt.Sprintf("**Reporter:** %v", p.Name()),
	// 			fmt.Sprintf("**Reason:** %v", cases.Title(language.English).String(string(r.Reason))),
	// 		}, "\n"),
	// 	}},
	// })
}

// Allow ...
func (Report) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

type reason string

// Type ...
func (reason) Type() string {
	return "reason"
}

// Options ...
func (reason) Options(cmd.Source) []string {
	return []string{
		"cheating",
		"allying",
		"spam",
		"threats",
		"glitching",
		"hostage",
		"exploiting",
		"toxic",
	}
}
