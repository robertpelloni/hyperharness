package subagents

// SubagentType defines categories of specialized subagents.
type SubagentType string

const (
	TypeCode     SubagentType = "code"
	TypeResearch SubagentType = "research"
	TypeReview   SubagentType = "review"
	TypePlan     SubagentType = "plan"
	TypeDoc      SubagentType = "doc"
	TypeBuild    SubagentType = "build"
	TypeTest     SubagentType = "test"
	TypeDebug    SubagentType = "debug"
	TypeSecurity SubagentType = "security"
	TypeDevOps   SubagentType = "devops"
)