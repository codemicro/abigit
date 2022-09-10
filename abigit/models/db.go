package models

import "github.com/bwmarrin/snowflake"

type SigningKey struct {
	ID  snowflake.ID `bun:"id,pk"`
	Key []byte       `bun:"key,notnull"`
}

type User struct {
	ID           snowflake.ID `bun:"id,pk"`
	EmailAddress string       `bun:"email,unique,nullzero"`
	ExternalID   string       `bun:"extern_id,unique,notnull"`
}

// type Game struct {
// 	bun.BaseModel `bun:"table:games"`

// 	ID             uuid.UUID `bun:"id,pk,type:varchar"`
// 	IntroMessageID string    `bun:"intro_message_id,notnull"`
// 	GuildID        string    `bun:"guild_id,notnull"`
// 	ChannelID      string    `bun:"channel_id,notnull"`
// 	LeaderID       string    `bun:"leader_id,notnull"`

// 	Started       bool          `bun:"started"`
// 	PhaseDuration time.Duration `bun:"phase_duration"`
// 	Phase         GamePhase     `bun:"phase,nullzero"`
// 	PhaseEndsAt   time.Time     `bun:"phase_ends_at,nullzero"`
// }

// type GamePhase uint

// func (g GamePhase) IsDay() bool {
// 	// Even numbers are daytimes
// 	return g%2 == 0
// }

// type Player struct {
// 	bun.BaseModel `bun:"table:players"`

// 	ID        uuid.UUID `bun:"id,pk,type:varchar"`
// 	DiscordID string    `bun:"discord_id,notnull"`
// 	GameID    uuid.UUID `bun:"game_id,notnull,type:varchar"`
// 	Game      *Game     `bun:"game,rel:belongs-to,join:game_id=id"`

// 	GameRole GameRole `bun:"game_role,nullzero"`
// 	IsDead   bool     `bun:"is_dead"`
// }
