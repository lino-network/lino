package types

type VoterDuty int

const (
	DutyVoter     VoterDuty = 0
	DutyApp       VoterDuty = 1
	DutyValidator VoterDuty = 2
	DutyPending   VoterDuty = 3 // pending is when voter is in unassign period
)
