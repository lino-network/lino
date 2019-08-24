package types

type VoterDuty int

const (
	DutyUnimplemented VoterDuty = -1 // voter duty module not implemented yet.
	DutyNop           VoterDuty = 0  // not a voter.
	DutyVoter         VoterDuty = 1
	DutyApp           VoterDuty = 2
	DutyValidator     VoterDuty = 3
)
